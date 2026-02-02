package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/Tencent/WeKnora/internal/utils"
	"github.com/gin-gonic/gin"
)

// SocialMediaHandler 自媒体处理器
type SocialMediaHandler struct {
	config           *config.Config
	kbService        interfaces.KnowledgeBaseService
	knowledgeService interfaces.KnowledgeService
}

// NewSocialMediaHandler 创建自媒体处理器
func NewSocialMediaHandler(
	config *config.Config,
	kbService interfaces.KnowledgeBaseService,
	knowledgeService interfaces.KnowledgeService,
) *SocialMediaHandler {
	return &SocialMediaHandler{
		config:           config,
		kbService:        kbService,
		knowledgeService: knowledgeService,
	}
}

// ExtractContentRequest 提取内容请求
type ExtractContentRequest struct {
	Platform string `json:"platform" binding:"required"` // xiaohongshu, douyin
	VideoURL string `json:"videoUrl" binding:"required"`
	KbID     string `json:"kbId" binding:"required"`
}

// CozeWorkflowRequest Coze workflow 请求
type CozeWorkflowRequest struct {
	WorkflowID string                 `json:"workflow_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

// CozeStreamResponse Coze 流式响应
type CozeStreamResponse struct {
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	NodeType       string `json:"node_type"`
	NodeID         string `json:"node_id"`
	NodeSeqID      string `json:"node_seq_id"`
	NodeTitle      string `json:"node_title"`
	NodeIsFinish   bool   `json:"node_is_finish"`
	NodeExecuteUUID string `json:"node_execute_uuid"`
}

// ExtractContent godoc
// @Summary      提取自媒体文案
// @Description  从自媒体平台提取文案内容并添加到知识库
// @Tags         自媒体
// @Accept       json
// @Produce      json
// @Param        request  body      ExtractContentRequest  true  "提取请求"
// @Success      200      {object}  map[string]interface{}  "提取成功"
// @Failure      400      {object}  errors.AppError         "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /social-media/extract [post]
func (h *SocialMediaHandler) ExtractContent(c *gin.Context) {
	ctx := c.Request.Context()

	var req ExtractContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to parse extract content request", err)
		c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	// 验证知识库是否存在
	kb, err := h.kbService.GetKnowledgeBaseByID(ctx, req.KbID)
	if err != nil || kb == nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"kbId": utils.SanitizeForLog(req.KbID)})
		c.Error(errors.NewNotFoundError("知识库不存在"))
		return
	}

	// 获取知识库的 Aliyun API Key
	aliyunAPIKey := h.getAliyunAPIKey(kb)
	if aliyunAPIKey == "" {
		logger.Error(ctx, "Aliyun API Key not configured for knowledge base")
		c.Error(errors.NewBadRequestError("知识库未配置阿里云 API Key"))
		return
	}

	// 调用 Coze workflow 提取内容
	content, err := h.callCozeWorkflow(ctx, req.VideoURL, aliyunAPIKey)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"platform": req.Platform,
			"videoUrl": utils.SanitizeForLog(req.VideoURL),
		})
		c.Error(errors.NewInternalServerError("提取文案失败: " + err.Error()))
		return
	}

	// 创建知识条目 - 使用 CreateKnowledgeFromPassage 方法
	passages := []string{content}
	knowledge, err := h.knowledgeService.CreateKnowledgeFromPassage(ctx, req.KbID, passages)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"kbId": req.KbID})
		c.Error(errors.NewInternalServerError("创建知识条目失败: " + err.Error()))
		return
	}

	// 更新知识条目的标题和来源
	knowledge.Title = fmt.Sprintf("%s_%s", req.Platform, time.Now().Format("20060102_150405"))
	knowledge.Source = fmt.Sprintf("%s:%s", req.Platform, req.VideoURL)
	if err := h.knowledgeService.UpdateKnowledge(ctx, knowledge); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"knowledgeId": knowledge.ID})
		// 不返回错误，因为知识已经创建成功
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文案提取成功",
		"data": gin.H{
			"content":     content,
			"knowledgeId": knowledge.ID,
		},
	})
}

// getAliyunAPIKey 从知识库配置中获取阿里云 API Key
func (h *SocialMediaHandler) getAliyunAPIKey(kb *types.KnowledgeBase) string {
	// 从知识库的 StorageConfig 中获取
	if kb.StorageConfig.AliyunAPIKey != "" {
		return kb.StorageConfig.AliyunAPIKey
	}
	return ""
}

// callCozeWorkflow 调用 Coze workflow API
func (h *SocialMediaHandler) callCozeWorkflow(ctx context.Context, videoURL, aliyunAPIKey string) (string, error) {
	// Coze API 配置
	cozeAPIURL := "https://api.coze.cn/v1/workflow/stream_run"
	cozeToken := "sat_bvn5AjfiAhnVaEIffM4myBfa443BqzzAdFu3Y9VjGsFNhL24fctmrUFafTKnMjw8"
	workflowID := "7596601391050588212"

	// 构建请求
	reqBody := CozeWorkflowRequest{
		WorkflowID: workflowID,
		Parameters: map[string]interface{}{
			"video_url": videoURL,
			"aliy_api":  aliyunAPIKey,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", cozeAPIURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+cozeToken)
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求 - 增加超时时间以应对 Coze Workflow 执行时间较长的情况
	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 解析流式响应
	content, err := h.parseStreamResponse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parse response failed: %w", err)
	}

	return content, nil
}

// parseStreamResponse 解析流式响应，提取 content 字段
func (h *SocialMediaHandler) parseStreamResponse(body io.Reader) (string, error) {
	scanner := bufio.NewScanner(body)
	var finalContent string
	var lastError string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// 处理 SSE 格式: "data: {...}"
		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")
			
			// 尝试解析为错误响应
			var errorResp struct {
				ErrorMessage string `json:"error_message"`
				ErrorCode    int    `json:"error_code"`
			}
			if err := json.Unmarshal([]byte(jsonData), &errorResp); err == nil && errorResp.ErrorMessage != "" {
				lastError = errorResp.ErrorMessage
				continue
			}

			// 尝试解析为普通响应
			var streamResp CozeStreamResponse
			if err := json.Unmarshal([]byte(jsonData), &streamResp); err != nil {
				// 跳过无法解析的行
				continue
			}

			// 提取 content 字段 - 收集所有非空内容
			if streamResp.Content != "" {
				// 优先使用 End 节点的内容
				if streamResp.NodeType == "End" {
					finalContent = streamResp.Content
				} else if finalContent == "" {
					// 如果还没有找到 End 节点，暂存其他节点的内容
					finalContent = streamResp.Content
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read stream failed: %w", err)
	}

	// 如果有错误信息，返回错误
	if lastError != "" {
		return "", fmt.Errorf("workflow execution failed: %s", lastError)
	}

	if finalContent == "" {
		return "", fmt.Errorf("no content found in response")
	}

	return finalContent, nil
}
