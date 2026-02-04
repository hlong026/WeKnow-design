package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BackupHandler handles backup and restore operations
type BackupHandler struct {
	db *gorm.DB
}

// NewBackupHandler creates a new backup handler
func NewBackupHandler(db *gorm.DB) *BackupHandler {
	return &BackupHandler{db: db}
}

// ExportRequest defines the export request parameters
type ExportRequest struct {
	IncludeTenants       bool `json:"include_tenants"`
	IncludeUsers         bool `json:"include_users"`
	IncludeKnowledgeBases bool `json:"include_knowledge_bases"`
	IncludeKnowledge     bool `json:"include_knowledge"`
	IncludeChunks        bool `json:"include_chunks"`
	IncludeSessions      bool `json:"include_sessions"`
	IncludeMessages      bool `json:"include_messages"`
	IncludeModels        bool `json:"include_models"`
	IncludeCredentials   bool `json:"include_credentials"`
	IncludeTags          bool `json:"include_tags"`
	IncludeAgents        bool `json:"include_agents"`
	IncludeMCPServices   bool `json:"include_mcp_services"`
}

// ExportData represents the exported data structure
type ExportData struct {
	Version        string                  `json:"version"`
	ExportTime     time.Time               `json:"export_time"`
	Tenants        []types.Tenant          `json:"tenants,omitempty"`
	Users          []types.User            `json:"users,omitempty"`
	KnowledgeBases []types.KnowledgeBase   `json:"knowledge_bases,omitempty"`
	Knowledge      []types.Knowledge       `json:"knowledge,omitempty"`
	Chunks         []types.Chunk           `json:"chunks,omitempty"`
	Sessions       []types.Session         `json:"sessions,omitempty"`
	Messages       []types.Message         `json:"messages,omitempty"`
	Models         []types.Model           `json:"models,omitempty"`
	Credentials    []types.ProviderCredential      `json:"credentials,omitempty"`
	Tags           []types.KnowledgeTag    `json:"tags,omitempty"`
	Agents         []types.CustomAgent     `json:"agents,omitempty"`
	MCPServices    []types.MCPService      `json:"mcp_services,omitempty"`
}

// ImportResult represents the import result
type ImportResult struct {
	TenantsImported       int      `json:"tenants_imported"`
	UsersImported         int      `json:"users_imported"`
	KnowledgeBasesImported int     `json:"knowledge_bases_imported"`
	KnowledgeImported     int      `json:"knowledge_imported"`
	ChunksImported        int      `json:"chunks_imported"`
	SessionsImported      int      `json:"sessions_imported"`
	MessagesImported      int      `json:"messages_imported"`
	ModelsImported        int      `json:"models_imported"`
	CredentialsImported   int      `json:"credentials_imported"`
	TagsImported          int      `json:"tags_imported"`
	AgentsImported        int      `json:"agents_imported"`
	MCPServicesImported   int      `json:"mcp_services_imported"`
	Errors                []string `json:"errors,omitempty"`
}

// Export godoc
// @Summary      导出数据库数据
// @Description  导出选定的数据库表数据为 JSON 格式
// @Tags         备份
// @Accept       json
// @Produce      application/octet-stream
// @Param        request body ExportRequest true "导出选项"
// @Success      200  {file}  file  "导出的数据文件"
// @Failure      400  {object}  map[string]interface{}  "请求参数错误"
// @Failure      500  {object}  map[string]interface{}  "服务器错误"
// @Router       /system/backup/export [post]
func (h *BackupHandler) Export(c *gin.Context) {
	ctx := logger.CloneContext(c.Request.Context())

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Invalid request parameters: " + err.Error()))
		return
	}

	// Get tenant ID from context
	tenantID := c.GetUint64(types.TenantIDContextKey.String())

	exportData := ExportData{
		Version:    Version,
		ExportTime: time.Now(),
	}

	// Export tenants
	if req.IncludeTenants {
		var tenants []types.Tenant
		if err := h.db.Where("id = ?", tenantID).Find(&tenants).Error; err != nil {
			logger.Error(ctx, "Failed to export tenants", "error", err)
		} else {
			exportData.Tenants = tenants
		}
	}

	// Export users
	if req.IncludeUsers {
		var users []types.User
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&users).Error; err != nil {
			logger.Error(ctx, "Failed to export users", "error", err)
		} else {
			// Clear sensitive data
			for i := range users {
				users[i].PasswordHash = ""
			}
			exportData.Users = users
		}
	}

	// Export knowledge bases
	if req.IncludeKnowledgeBases {
		var kbs []types.KnowledgeBase
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&kbs).Error; err != nil {
			logger.Error(ctx, "Failed to export knowledge bases", "error", err)
		} else {
			exportData.KnowledgeBases = kbs
		}
	}

	// Export knowledge
	if req.IncludeKnowledge {
		var knowledge []types.Knowledge
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&knowledge).Error; err != nil {
			logger.Error(ctx, "Failed to export knowledge", "error", err)
		} else {
			exportData.Knowledge = knowledge
		}
	}

	// Export chunks
	if req.IncludeChunks {
		var chunks []types.Chunk
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&chunks).Error; err != nil {
			logger.Error(ctx, "Failed to export chunks", "error", err)
		} else {
			exportData.Chunks = chunks
		}
	}

	// Export sessions
	if req.IncludeSessions {
		var sessions []types.Session
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&sessions).Error; err != nil {
			logger.Error(ctx, "Failed to export sessions", "error", err)
		} else {
			exportData.Sessions = sessions
		}
	}

	// Export messages
	if req.IncludeMessages {
		var messages []types.Message
		// Get session IDs for this tenant first
		var sessionIDs []uint64
		h.db.Model(&types.Session{}).Where("tenant_id = ?", tenantID).Pluck("id", &sessionIDs)
		if len(sessionIDs) > 0 {
			if err := h.db.Where("session_id IN ?", sessionIDs).Find(&messages).Error; err != nil {
				logger.Error(ctx, "Failed to export messages", "error", err)
			} else {
				exportData.Messages = messages
			}
		}
	}

	// Export models
	if req.IncludeModels {
		var models []types.Model
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&models).Error; err != nil {
			logger.Error(ctx, "Failed to export models", "error", err)
		} else {
			exportData.Models = models
		}
	}

	// Export credentials
	if req.IncludeCredentials {
		var credentials []types.ProviderCredential
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&credentials).Error; err != nil {
			logger.Error(ctx, "Failed to export credentials", "error", err)
		} else {
			// Clear sensitive data
			for i := range credentials {
				credentials[i] = *credentials[i].HideSensitiveInfo()
			}
			exportData.Credentials = credentials
		}
	}

	// Export tags
	if req.IncludeTags {
		var tags []types.KnowledgeTag
		// Get knowledge base IDs for this tenant
		var kbIDs []uint64
		h.db.Model(&types.KnowledgeBase{}).Where("tenant_id = ?", tenantID).Pluck("id", &kbIDs)
		if len(kbIDs) > 0 {
			if err := h.db.Where("knowledge_base_id IN ?", kbIDs).Find(&tags).Error; err != nil {
				logger.Error(ctx, "Failed to export tags", "error", err)
			} else {
				exportData.Tags = tags
			}
		}
	}

	// Export agents
	if req.IncludeAgents {
		var agents []types.CustomAgent
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&agents).Error; err != nil {
			logger.Error(ctx, "Failed to export agents", "error", err)
		} else {
			exportData.Agents = agents
		}
	}

	// Export MCP services
	if req.IncludeMCPServices {
		var mcpServices []types.MCPService
		if err := h.db.Where("tenant_id = ?", tenantID).Find(&mcpServices).Error; err != nil {
			logger.Error(ctx, "Failed to export MCP services", "error", err)
		} else {
			exportData.MCPServices = mcpServices
		}
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		c.Error(errors.NewInternalServerError("Failed to serialize export data: " + err.Error()))
		return
	}

	// Create ZIP file
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add JSON file to ZIP
	jsonFile, err := zipWriter.Create("backup.json")
	if err != nil {
		c.Error(errors.NewInternalServerError("Failed to create zip file: " + err.Error()))
		return
	}
	if _, err := jsonFile.Write(jsonData); err != nil {
		c.Error(errors.NewInternalServerError("Failed to write to zip file: " + err.Error()))
		return
	}

	if err := zipWriter.Close(); err != nil {
		c.Error(errors.NewInternalServerError("Failed to close zip file: " + err.Error()))
		return
	}

	// Set response headers
	filename := fmt.Sprintf("weknora_backup_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))

	logger.Info(ctx, "Export completed successfully", "size", buf.Len())
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}


// Import godoc
// @Summary      导入数据库数据
// @Description  从 ZIP 文件导入数据库数据
// @Tags         备份
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formance file true "备份文件"
// @Param        skip_existing formData bool false "跳过已存在的记录"
// @Success      200  {object}  ImportResult  "导入结果"
// @Failure      400  {object}  map[string]interface{}  "请求参数错误"
// @Failure      500  {object}  map[string]interface{}  "服务器错误"
// @Router       /system/backup/import [post]
func (h *BackupHandler) Import(c *gin.Context) {
	ctx := logger.CloneContext(c.Request.Context())

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(errors.NewBadRequestError("No file uploaded: " + err.Error()))
		return
	}

	// Check file extension
	if file.Filename[len(file.Filename)-4:] != ".zip" {
		c.Error(errors.NewBadRequestError("Invalid file format. Please upload a .zip file"))
		return
	}

	// Get skip_existing parameter
	skipExisting := c.PostForm("skip_existing") == "true"

	// Get tenant ID from context
	tenantID := c.GetUint64(types.TenantIDContextKey.String())

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.Error(errors.NewInternalServerError("Failed to open uploaded file: " + err.Error()))
		return
	}
	defer src.Close()

	// Read file content
	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.Error(errors.NewInternalServerError("Failed to read uploaded file: " + err.Error()))
		return
	}

	// Open ZIP file
	zipReader, err := zip.NewReader(bytes.NewReader(fileContent), int64(len(fileContent)))
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid ZIP file: " + err.Error()))
		return
	}

	// Find and read backup.json
	var exportData ExportData
	for _, f := range zipReader.File {
		if f.Name == "backup.json" {
			rc, err := f.Open()
			if err != nil {
				c.Error(errors.NewInternalServerError("Failed to open backup.json: " + err.Error()))
				return
			}
			defer rc.Close()

			jsonData, err := io.ReadAll(rc)
			if err != nil {
				c.Error(errors.NewInternalServerError("Failed to read backup.json: " + err.Error()))
				return
			}

			if err := json.Unmarshal(jsonData, &exportData); err != nil {
				c.Error(errors.NewBadRequestError("Invalid backup.json format: " + err.Error()))
				return
			}
			break
		}
	}

	if exportData.Version == "" {
		c.Error(errors.NewBadRequestError("backup.json not found in ZIP file"))
		return
	}

	result := ImportResult{}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Import knowledge bases (must be before knowledge and tags)
	if len(exportData.KnowledgeBases) > 0 {
		for _, kb := range exportData.KnowledgeBases {
			kb.TenantID = tenantID
			if skipExisting {
				var existing types.KnowledgeBase
				if tx.Where("id = ? AND tenant_id = ?", kb.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&kb).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import knowledge base %s: %v", kb.Name, err))
			} else {
				result.KnowledgeBasesImported++
			}
		}
	}

	// Import tags
	if len(exportData.Tags) > 0 {
		for _, tag := range exportData.Tags {
			if skipExisting {
				var existing types.KnowledgeTag
				if tx.Where("id = ?", tag.ID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&tag).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import tag %s: %v", tag.Name, err))
			} else {
				result.TagsImported++
			}
		}
	}

	// Import knowledge
	if len(exportData.Knowledge) > 0 {
		for _, k := range exportData.Knowledge {
			k.TenantID = tenantID
			if skipExisting {
				var existing types.Knowledge
				if tx.Where("id = ? AND tenant_id = ?", k.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&k).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import knowledge %s: %v", k.Title, err))
			} else {
				result.KnowledgeImported++
			}
		}
	}

	// Import chunks
	if len(exportData.Chunks) > 0 {
		for _, chunk := range exportData.Chunks {
			chunk.TenantID = tenantID
			if skipExisting {
				var existing types.Chunk
				if tx.Where("id = ? AND tenant_id = ?", chunk.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&chunk).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import chunk %d: %v", chunk.ID, err))
			} else {
				result.ChunksImported++
			}
		}
	}

	// Import sessions
	if len(exportData.Sessions) > 0 {
		for _, session := range exportData.Sessions {
			session.TenantID = tenantID
			if skipExisting {
				var existing types.Session
				if tx.Where("id = ? AND tenant_id = ?", session.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&session).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import session %s: %v", session.Title, err))
			} else {
				result.SessionsImported++
			}
		}
	}

	// Import messages
	if len(exportData.Messages) > 0 {
		for _, msg := range exportData.Messages {
			if skipExisting {
				var existing types.Message
				if tx.Where("id = ?", msg.ID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&msg).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import message %d: %v", msg.ID, err))
			} else {
				result.MessagesImported++
			}
		}
	}

	// Import models
	if len(exportData.Models) > 0 {
		for _, model := range exportData.Models {
			model.TenantID = tenantID
			if skipExisting {
				var existing types.Model
				if tx.Where("id = ? AND tenant_id = ?", model.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&model).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import model %s: %v", model.Name, err))
			} else {
				result.ModelsImported++
			}
		}
	}

	// Import agents
	if len(exportData.Agents) > 0 {
		for _, agent := range exportData.Agents {
			agent.TenantID = tenantID
			if skipExisting {
				var existing types.CustomAgent
				if tx.Where("id = ? AND tenant_id = ?", agent.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&agent).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import agent %s: %v", agent.Name, err))
			} else {
				result.AgentsImported++
			}
		}
	}

	// Import MCP services
	if len(exportData.MCPServices) > 0 {
		for _, svc := range exportData.MCPServices {
			svc.TenantID = tenantID
			if skipExisting {
				var existing types.MCPService
				if tx.Where("id = ? AND tenant_id = ?", svc.ID, tenantID).First(&existing).Error == nil {
					continue
				}
			}
			if err := tx.Create(&svc).Error; err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to import MCP service %s: %v", svc.Name, err))
			} else {
				result.MCPServicesImported++
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Error(errors.NewInternalServerError("Failed to commit import transaction: " + err.Error()))
		return
	}

	logger.Info(ctx, "Import completed", "result", result)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"success": true,
		"data":    result,
	})
}

// GetExportOptions godoc
// @Summary      获取导出选项
// @Description  获取可用的导出选项和数据统计
// @Tags         备份
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "导出选项"
// @Router       /system/backup/options [get]
func (h *BackupHandler) GetExportOptions(c *gin.Context) {
	ctx := logger.CloneContext(c.Request.Context())
	tenantID := c.GetUint64(types.TenantIDContextKey.String())

	// Get counts for each table
	var tenantCount, userCount, kbCount, knowledgeCount, chunkCount int64
	var sessionCount, messageCount, modelCount, credentialCount, tagCount int64
	var agentCount, mcpServiceCount int64

	h.db.Model(&types.Tenant{}).Where("id = ?", tenantID).Count(&tenantCount)
	h.db.Model(&types.User{}).Where("tenant_id = ?", tenantID).Count(&userCount)
	h.db.Model(&types.KnowledgeBase{}).Where("tenant_id = ?", tenantID).Count(&kbCount)
	h.db.Model(&types.Knowledge{}).Where("tenant_id = ?", tenantID).Count(&knowledgeCount)
	h.db.Model(&types.Chunk{}).Where("tenant_id = ?", tenantID).Count(&chunkCount)
	h.db.Model(&types.Session{}).Where("tenant_id = ?", tenantID).Count(&sessionCount)
	
	// Count messages through sessions
	var sessionIDs []uint64
	h.db.Model(&types.Session{}).Where("tenant_id = ?", tenantID).Pluck("id", &sessionIDs)
	if len(sessionIDs) > 0 {
		h.db.Model(&types.Message{}).Where("session_id IN ?", sessionIDs).Count(&messageCount)
	}
	
	h.db.Model(&types.Model{}).Where("tenant_id = ?", tenantID).Count(&modelCount)
	h.db.Model(&types.ProviderCredential{}).Where("tenant_id = ?", tenantID).Count(&credentialCount)
	
	// Count tags through knowledge bases
	var kbIDs []uint64
	h.db.Model(&types.KnowledgeBase{}).Where("tenant_id = ?", tenantID).Pluck("id", &kbIDs)
	if len(kbIDs) > 0 {
		h.db.Model(&types.KnowledgeTag{}).Where("knowledge_base_id IN ?", kbIDs).Count(&tagCount)
	}
	
	h.db.Model(&types.CustomAgent{}).Where("tenant_id = ?", tenantID).Count(&agentCount)
	h.db.Model(&types.MCPService{}).Where("tenant_id = ?", tenantID).Count(&mcpServiceCount)

	options := []map[string]interface{}{
		{"key": "include_tenants", "label": "租户配置", "count": tenantCount},
		{"key": "include_users", "label": "用户数据", "count": userCount},
		{"key": "include_knowledge_bases", "label": "知识库", "count": kbCount},
		{"key": "include_knowledge", "label": "知识条目", "count": knowledgeCount},
		{"key": "include_chunks", "label": "知识分块", "count": chunkCount},
		{"key": "include_sessions", "label": "会话记录", "count": sessionCount},
		{"key": "include_messages", "label": "消息记录", "count": messageCount},
		{"key": "include_models", "label": "模型配置", "count": modelCount},
		{"key": "include_credentials", "label": "凭证配置", "count": credentialCount},
		{"key": "include_tags", "label": "标签数据", "count": tagCount},
		{"key": "include_agents", "label": "智能体配置", "count": agentCount},
		{"key": "include_mcp_services", "label": "MCP服务", "count": mcpServiceCount},
	}

	logger.Info(ctx, "Export options retrieved")
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"success": true,
		"data":    options,
	})
}
