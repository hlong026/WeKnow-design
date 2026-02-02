package provider

import "github.com/Tencent/WeKnora/internal/types"

// PresetModel 预置模型定义
type PresetModel struct {
	ModelID      string                  `json:"model_id"`
	DisplayName  string                  `json:"display_name"`
	ModelType    types.ModelType         `json:"model_type"`
	Capabilities []types.ModelCapability `json:"capabilities"`
	ContextSize  int                     `json:"context_size"`
	Pricing      *types.PricingInfo      `json:"pricing,omitempty"`
	Deprecated   bool                    `json:"deprecated"`
}

// ProviderFeatures 厂商特性支持
type ProviderFeatures struct {
	SupportsStreaming    bool `json:"supports_streaming"`
	SupportsFunctionCall bool `json:"supports_function_call"`
	SupportsVision       bool `json:"supports_vision"`
	SupportsJSON         bool `json:"supports_json_mode"`
	SupportsCustomModel  bool `json:"supports_custom_model"`
}

// AuthField 认证字段定义
type AuthField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"` // text, password, select
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
	HelpText    string `json:"help_text"`
}

// AuthConfiguration 认证配置
type AuthConfiguration struct {
	Type     string      `json:"type"` // api_key, oauth, custom
	Fields   []AuthField `json:"fields"`
	HelpText string      `json:"help_text"`
	HelpURL  string      `json:"help_url"`
}

// ProviderMetadata 厂商完整元数据
type ProviderMetadata struct {
	Name           ProviderName               `json:"name"`
	DisplayName    string                     `json:"display_name"`
	Description    string                     `json:"description"`
	Icon           string                     `json:"icon"`
	Website        string                     `json:"website"`
	DocsURL        string                     `json:"docs_url"`
	AuthConfig     AuthConfiguration          `json:"auth_config"`
	SupportedTypes []types.ModelType          `json:"supported_types"`
	PresetModels   []PresetModel              `json:"preset_models"`
	Endpoints      map[types.ModelType]string `json:"endpoints"`
	Features       ProviderFeatures           `json:"features"`
}

// GetProviderMetadata 获取厂商完整元数据
func GetProviderMetadata(name ProviderName) *ProviderMetadata {
	p, ok := Get(name)
	if !ok {
		return nil
	}
	info := p.Info()
	
	return &ProviderMetadata{
		Name:           info.Name,
		DisplayName:    info.DisplayName,
		Description:    info.Description,
		SupportedTypes: info.ModelTypes,
		Endpoints:      info.DefaultURLs,
		Features: ProviderFeatures{
			SupportsStreaming:    true,
			SupportsFunctionCall: name == ProviderOpenAI || name == ProviderAliyun,
			SupportsVision:       name == ProviderOpenAI || name == ProviderAliyun,
			SupportsJSON:         name == ProviderOpenAI,
			SupportsCustomModel:  true,
		},
	}
}
