package handler

import (
	"net/http"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/Tencent/WeKnora/internal/types"
	secutils "github.com/Tencent/WeKnora/internal/utils"
	"github.com/gin-gonic/gin"
)

// ProviderHandler 厂商处理器
type ProviderHandler struct{}

// NewProviderHandler 创建厂商处理器
func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{}
}

// ProviderDTO 厂商信息 DTO
type ProviderDTO struct {
	Name           string                 `json:"name"`
	DisplayName    string                 `json:"display_name"`
	Description    string                 `json:"description"`
	Icon           string                 `json:"icon,omitempty"`
	Website        string                 `json:"website,omitempty"`
	DocsURL        string                 `json:"docs_url,omitempty"`
	AuthConfig     AuthConfigDTO          `json:"auth_config"`
	SupportedTypes []string               `json:"supported_types"`
	PresetModels   []PresetModelDTO       `json:"preset_models,omitempty"`
	Endpoints      map[string]string      `json:"endpoints"`
	Features       ProviderFeaturesDTO    `json:"features"`
}

// AuthConfigDTO 认证配置 DTO
type AuthConfigDTO struct {
	Type     string         `json:"type"`
	Fields   []AuthFieldDTO `json:"fields"`
	HelpText string         `json:"help_text,omitempty"`
	HelpURL  string         `json:"help_url,omitempty"`
}

// AuthFieldDTO 认证字段 DTO
type AuthFieldDTO struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
	HelpText    string `json:"help_text,omitempty"`
}

// PresetModelDTO 预置模型 DTO
type PresetModelDTO struct {
	ModelID      string         `json:"model_id"`
	DisplayName  string         `json:"display_name"`
	ModelType    string         `json:"model_type"`
	Capabilities []string       `json:"capabilities"`
	ContextSize  int            `json:"context_size"`
	Pricing      *PricingDTO    `json:"pricing,omitempty"`
	Deprecated   bool           `json:"deprecated,omitempty"`
}

// PricingDTO 价格信息 DTO
type PricingDTO struct {
	InputPrice  float64 `json:"input_price"`
	OutputPrice float64 `json:"output_price"`
	Currency    string  `json:"currency"`
}

// ProviderFeaturesDTO 厂商特性 DTO
type ProviderFeaturesDTO struct {
	SupportsStreaming    bool `json:"supports_streaming"`
	SupportsFunctionCall bool `json:"supports_function_call"`
	SupportsVision       bool `json:"supports_vision"`
	SupportsJSON         bool `json:"supports_json_mode"`
	SupportsCustomModel  bool `json:"supports_custom_model"`
}

// modelTypeToFrontendStr 将后端 ModelType 转换为前端字符串
func modelTypeToFrontendStr(mt types.ModelType) string {
	switch mt {
	case types.ModelTypeKnowledgeQA:
		return "chat"
	case types.ModelTypeEmbedding:
		return "embedding"
	case types.ModelTypeRerank:
		return "rerank"
	case types.ModelTypeVLLM:
		return "vllm"
	default:
		return string(mt)
	}
}

// convertProviderInfo 转换 ProviderInfo 为 DTO
func convertProviderInfo(info provider.ProviderInfo) ProviderDTO {
	// 转换端点
	endpoints := make(map[string]string)
	for mt, url := range info.DefaultURLs {
		endpoints[modelTypeToFrontendStr(mt)] = url
	}

	// 转换支持的类型
	supportedTypes := make([]string, len(info.ModelTypes))
	for i, mt := range info.ModelTypes {
		supportedTypes[i] = modelTypeToFrontendStr(mt)
	}

	// 获取预置模型
	presetModels := provider.GetPresetModels(info.Name)
	presetModelDTOs := make([]PresetModelDTO, len(presetModels))
	for i, pm := range presetModels {
		caps := make([]string, len(pm.Capabilities))
		for j, c := range pm.Capabilities {
			caps[j] = string(c)
		}
		
		var pricing *PricingDTO
		if pm.Pricing != nil {
			pricing = &PricingDTO{
				InputPrice:  pm.Pricing.InputPrice,
				OutputPrice: pm.Pricing.OutputPrice,
				Currency:    pm.Pricing.Currency,
			}
		}
		
		presetModelDTOs[i] = PresetModelDTO{
			ModelID:      pm.ModelID,
			DisplayName:  pm.DisplayName,
			ModelType:    modelTypeToFrontendStr(pm.ModelType),
			Capabilities: caps,
			ContextSize:  pm.ContextSize,
			Pricing:      pricing,
			Deprecated:   pm.Deprecated,
		}
	}

	// 构建认证配置
	authConfig := AuthConfigDTO{
		Type: "api_key",
		Fields: []AuthFieldDTO{
			{
				Key:         "api_key",
				Label:       "API Key",
				Type:        "password",
				Required:    info.RequiresAuth,
				Placeholder: "请输入 API Key",
			},
		},
	}

	// 根据厂商添加额外字段
	switch info.Name {
	case provider.ProviderAliyun:
		authConfig.HelpURL = "https://bailian.console.aliyun.com/#/api-key"
		authConfig.HelpText = "在阿里云百炼平台获取 API Key"
	case provider.ProviderOpenAI:
		authConfig.HelpURL = "https://platform.openai.com/api-keys"
		authConfig.HelpText = "在 OpenAI 平台获取 API Key"
	case provider.ProviderZhipu:
		authConfig.HelpURL = "https://open.bigmodel.cn/usercenter/apikeys"
		authConfig.HelpText = "在智谱 AI 开放平台获取 API Key"
	case provider.ProviderDeepSeek:
		authConfig.HelpURL = "https://platform.deepseek.com/api_keys"
		authConfig.HelpText = "在 DeepSeek 平台获取 API Key"
	case provider.ProviderSiliconFlow:
		authConfig.HelpURL = "https://cloud.siliconflow.cn/account/ak"
		authConfig.HelpText = "在硅基流动控制台获取 API Key"
	}

	// 添加额外配置字段
	for _, field := range info.ExtraFields {
		authConfig.Fields = append(authConfig.Fields, AuthFieldDTO{
			Key:         field.Key,
			Label:       field.Label,
			Type:        field.Type,
			Required:    field.Required,
			Placeholder: field.Placeholder,
		})
	}

	// 构建特性
	features := ProviderFeaturesDTO{
		SupportsStreaming:    true,
		SupportsCustomModel:  true,
	}
	
	switch info.Name {
	case provider.ProviderOpenAI:
		features.SupportsFunctionCall = true
		features.SupportsVision = true
		features.SupportsJSON = true
	case provider.ProviderAliyun:
		features.SupportsFunctionCall = true
		features.SupportsVision = true
	case provider.ProviderZhipu:
		features.SupportsFunctionCall = true
		features.SupportsVision = true
	case provider.ProviderDeepSeek:
		features.SupportsFunctionCall = true
	}

	return ProviderDTO{
		Name:           string(info.Name),
		DisplayName:    info.DisplayName,
		Description:    info.Description,
		AuthConfig:     authConfig,
		SupportedTypes: supportedTypes,
		PresetModels:   presetModelDTOs,
		Endpoints:      endpoints,
		Features:       features,
	}
}

// ListProviders godoc
// @Summary      获取厂商列表
// @Description  获取所有支持的模型厂商列表
// @Tags         厂商管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "厂商列表"
// @Security     Bearer
// @Router       /providers [get]
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Listing all providers")

	providers := provider.List()
	result := make([]ProviderDTO, len(providers))
	for i, p := range providers {
		result[i] = convertProviderInfo(p)
	}

	logger.Infof(ctx, "Retrieved %d providers", len(result))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetProvider godoc
// @Summary      获取厂商详情
// @Description  获取指定厂商的详细信息
// @Tags         厂商管理
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "厂商名称"
// @Success      200       {object}  map[string]interface{}  "厂商详情"
// @Failure      404       {object}  errors.AppError         "厂商不存在"
// @Security     Bearer
// @Router       /providers/{provider} [get]
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	ctx := c.Request.Context()

	providerName := secutils.SanitizeForLog(c.Param("provider"))
	logger.Infof(ctx, "Getting provider: %s", providerName)

	p, ok := provider.Get(provider.ProviderName(providerName))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Provider not found",
		})
		return
	}

	result := convertProviderInfo(p.Info())
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetProviderModels godoc
// @Summary      获取厂商预置模型
// @Description  获取指定厂商的预置模型列表
// @Tags         厂商管理
// @Accept       json
// @Produce      json
// @Param        provider    path      string  true   "厂商名称"
// @Param        model_type  query     string  false  "模型类型 (chat, embedding, rerank)"
// @Success      200         {object}  map[string]interface{}  "预置模型列表"
// @Security     Bearer
// @Router       /providers/{provider}/models [get]
func (h *ProviderHandler) GetProviderModels(c *gin.Context) {
	ctx := c.Request.Context()

	providerName := secutils.SanitizeForLog(c.Param("provider"))
	modelType := c.Query("model_type")
	logger.Infof(ctx, "Getting preset models for provider: %s, type: %s", providerName, modelType)

	var presetModels []provider.PresetModel
	if modelType != "" {
		// 转换前端类型到后端类型
		var backendType types.ModelType
		switch modelType {
		case "chat":
			backendType = types.ModelTypeKnowledgeQA
		case "embedding":
			backendType = types.ModelTypeEmbedding
		case "rerank":
			backendType = types.ModelTypeRerank
		case "vllm":
			backendType = types.ModelTypeVLLM
		default:
			backendType = types.ModelType(modelType)
		}
		presetModels = provider.GetPresetModelsByType(provider.ProviderName(providerName), backendType)
	} else {
		presetModels = provider.GetPresetModels(provider.ProviderName(providerName))
	}

	// 转换为 DTO
	result := make([]PresetModelDTO, len(presetModels))
	for i, pm := range presetModels {
		caps := make([]string, len(pm.Capabilities))
		for j, c := range pm.Capabilities {
			caps[j] = string(c)
		}
		
		var pricing *PricingDTO
		if pm.Pricing != nil {
			pricing = &PricingDTO{
				InputPrice:  pm.Pricing.InputPrice,
				OutputPrice: pm.Pricing.OutputPrice,
				Currency:    pm.Pricing.Currency,
			}
		}
		
		result[i] = PresetModelDTO{
			ModelID:      pm.ModelID,
			DisplayName:  pm.DisplayName,
			ModelType:    modelTypeToFrontendStr(pm.ModelType),
			Capabilities: caps,
			ContextSize:  pm.ContextSize,
			Pricing:      pricing,
			Deprecated:   pm.Deprecated,
		}
	}

	logger.Infof(ctx, "Retrieved %d preset models", len(result))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
