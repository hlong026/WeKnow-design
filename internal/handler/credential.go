package handler

import (
	"net/http"

	"github.com/Tencent/WeKnora/internal/application/service"
	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	secutils "github.com/Tencent/WeKnora/internal/utils"
	"github.com/gin-gonic/gin"
)

// CredentialHandler 凭证处理器
type CredentialHandler struct {
	service interfaces.CredentialService
}

// NewCredentialHandler 创建凭证处理器
func NewCredentialHandler(service interfaces.CredentialService) *CredentialHandler {
	return &CredentialHandler{service: service}
}

// CreateCredentialRequest 创建凭证请求
type CreateCredentialRequest struct {
	Provider    string            `json:"provider" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Credentials map[string]string `json:"credentials" binding:"required"`
	BaseURL     string            `json:"base_url"`
	IsDefault   bool              `json:"is_default"`
	QuotaConfig *types.QuotaConfig `json:"quota_config"`
}

// CreateCredential godoc
// @Summary      创建凭证
// @Description  创建新的厂商凭证
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCredentialRequest  true  "凭证信息"
// @Success      201      {object}  map[string]interface{}  "创建的凭证"
// @Failure      400      {object}  errors.AppError         "请求参数错误"
// @Security     Bearer
// @Router       /credentials [post]
func (h *CredentialHandler) CreateCredential(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Start creating credential")

	var req CreateCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	// 验证 provider 是否存在
	if _, ok := provider.Get(provider.ProviderName(req.Provider)); !ok {
		c.Error(errors.NewBadRequestError("Unknown provider: " + req.Provider))
		return
	}

	credential := &types.ProviderCredential{
		Provider:    req.Provider,
		Name:        secutils.SanitizeForLog(req.Name),
		Credentials: req.Credentials,
		BaseURL:     req.BaseURL,
		IsDefault:   req.IsDefault,
		QuotaConfig: req.QuotaConfig,
		Status:      types.CredentialStatusActive,
	}

	if err := h.service.CreateCredential(ctx, credential); err != nil {
		if err == service.ErrCredentialExists {
			c.Error(errors.NewBadRequestError("Credential with the same name already exists"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Credential created successfully: %s", credential.ID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    credential.HideSensitiveInfo(),
	})
}

// GetCredential godoc
// @Summary      获取凭证详情
// @Description  根据ID获取凭证详情
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "凭证ID"
// @Success      200  {object}  map[string]interface{}  "凭证详情"
// @Failure      404  {object}  errors.AppError         "凭证不存在"
// @Security     Bearer
// @Router       /credentials/{id} [get]
func (h *CredentialHandler) GetCredential(c *gin.Context) {
	ctx := c.Request.Context()

	id := secutils.SanitizeForLog(c.Param("id"))
	if id == "" {
		c.Error(errors.NewBadRequestError("Credential ID cannot be empty"))
		return
	}

	credential, err := h.service.GetCredentialByID(ctx, id)
	if err != nil {
		if err == service.ErrCredentialNotFound {
			c.Error(errors.NewNotFoundError("Credential not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    credential.HideSensitiveInfo(),
	})
}

// ListCredentials godoc
// @Summary      获取凭证列表
// @Description  获取当前租户的凭证列表
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        provider  query     string  false  "厂商名称"
// @Success      200       {object}  map[string]interface{}  "凭证列表"
// @Security     Bearer
// @Router       /credentials [get]
func (h *CredentialHandler) ListCredentials(c *gin.Context) {
	ctx := c.Request.Context()

	providerName := c.Query("provider")
	logger.Infof(ctx, "Listing credentials for provider: %s", secutils.SanitizeForLog(providerName))

	credentials, err := h.service.ListCredentials(ctx, providerName)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	// 隐藏敏感信息
	result := make([]*types.ProviderCredential, len(credentials))
	for i, cred := range credentials {
		result[i] = cred.HideSensitiveInfo()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateCredentialRequest 更新凭证请求
type UpdateCredentialRequest struct {
	Name        string            `json:"name"`
	Credentials map[string]string `json:"credentials"`
	BaseURL     string            `json:"base_url"`
	IsDefault   bool              `json:"is_default"`
	QuotaConfig *types.QuotaConfig `json:"quota_config"`
}

// UpdateCredential godoc
// @Summary      更新凭证
// @Description  更新凭证配置信息
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "凭证ID"
// @Param        request  body      UpdateCredentialRequest  true  "更新信息"
// @Success      200      {object}  map[string]interface{}  "更新后的凭证"
// @Failure      404      {object}  errors.AppError         "凭证不存在"
// @Security     Bearer
// @Router       /credentials/{id} [put]
func (h *CredentialHandler) UpdateCredential(c *gin.Context) {
	ctx := c.Request.Context()

	id := secutils.SanitizeForLog(c.Param("id"))
	if id == "" {
		c.Error(errors.NewBadRequestError("Credential ID cannot be empty"))
		return
	}

	var req UpdateCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	// 获取现有凭证
	credential, err := h.service.GetCredentialByID(ctx, id)
	if err != nil {
		if err == service.ErrCredentialNotFound {
			c.Error(errors.NewNotFoundError("Credential not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	// 更新字段
	if req.Name != "" {
		credential.Name = req.Name
	}
	if req.Credentials != nil {
		// 合并凭证，只更新提供的字段
		for k, v := range req.Credentials {
			if v != "" {
				credential.Credentials[k] = v
			}
		}
	}
	if req.BaseURL != "" {
		credential.BaseURL = req.BaseURL
	}
	credential.IsDefault = req.IsDefault
	if req.QuotaConfig != nil {
		credential.QuotaConfig = req.QuotaConfig
	}

	if err := h.service.UpdateCredential(ctx, credential); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Credential updated successfully: %s", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    credential.HideSensitiveInfo(),
	})
}

// DeleteCredential godoc
// @Summary      删除凭证
// @Description  删除指定的凭证
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "凭证ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Failure      404  {object}  errors.AppError         "凭证不存在"
// @Security     Bearer
// @Router       /credentials/{id} [delete]
func (h *CredentialHandler) DeleteCredential(c *gin.Context) {
	ctx := c.Request.Context()

	id := secutils.SanitizeForLog(c.Param("id"))
	if id == "" {
		c.Error(errors.NewBadRequestError("Credential ID cannot be empty"))
		return
	}

	if err := h.service.DeleteCredential(ctx, id); err != nil {
		if err == service.ErrCredentialNotFound {
			c.Error(errors.NewNotFoundError("Credential not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Credential deleted successfully: %s", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Credential deleted",
	})
}

// TestCredential godoc
// @Summary      测试凭证连接
// @Description  测试凭证是否可用
// @Tags         凭证管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "凭证ID"
// @Success      200  {object}  map[string]interface{}  "测试结果"
// @Failure      404  {object}  errors.AppError         "凭证不存在"
// @Security     Bearer
// @Router       /credentials/{id}/test [post]
func (h *CredentialHandler) TestCredential(c *gin.Context) {
	ctx := c.Request.Context()

	id := secutils.SanitizeForLog(c.Param("id"))
	if id == "" {
		c.Error(errors.NewBadRequestError("Credential ID cannot be empty"))
		return
	}

	if err := h.service.TestCredential(ctx, id); err != nil {
		if err == service.ErrCredentialNotFound {
			c.Error(errors.NewNotFoundError("Credential not found"))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection test passed",
	})
}
