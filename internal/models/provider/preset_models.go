package provider

import "github.com/Tencent/WeKnora/internal/types"

// OpenAI 预置模型
var OpenAIPresetModels = []PresetModel{
	// Chat Models
	{
		ModelID: "gpt-4o", DisplayName: "GPT-4o", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityVision, types.CapabilityFunctionCall, types.CapabilityJSON},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 2.5, OutputPrice: 10, Currency: "USD"},
	},
	{
		ModelID: "gpt-4o-mini", DisplayName: "GPT-4o Mini", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityVision, types.CapabilityFunctionCall},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 0.15, OutputPrice: 0.6, Currency: "USD"},
	},
	{
		ModelID: "gpt-4-turbo", DisplayName: "GPT-4 Turbo", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityVision, types.CapabilityFunctionCall},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 10, OutputPrice: 30, Currency: "USD"},
	},
	{
		ModelID: "o1", DisplayName: "o1", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 200000, Pricing: &types.PricingInfo{InputPrice: 15, OutputPrice: 60, Currency: "USD"},
	},
	// Embedding Models
	{
		ModelID: "text-embedding-3-large", DisplayName: "Text Embedding 3 Large", ModelType: types.ModelTypeEmbedding,
		Capabilities: []types.ModelCapability{types.CapabilityEmbedding},
		ContextSize: 8191, Pricing: &types.PricingInfo{InputPrice: 0.13, Currency: "USD"},
	},
	{
		ModelID: "text-embedding-3-small", DisplayName: "Text Embedding 3 Small", ModelType: types.ModelTypeEmbedding,
		Capabilities: []types.ModelCapability{types.CapabilityEmbedding},
		ContextSize: 8191, Pricing: &types.PricingInfo{InputPrice: 0.02, Currency: "USD"},
	},
}

// 阿里云 DashScope 预置模型
var AliyunPresetModels = []PresetModel{
	{
		ModelID: "qwen-max", DisplayName: "通义千问-Max", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityFunctionCall},
		ContextSize: 32000, Pricing: &types.PricingInfo{InputPrice: 0.02, OutputPrice: 0.06, Currency: "CNY"},
	},
	{
		ModelID: "qwen-plus", DisplayName: "通义千问-Plus", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityFunctionCall},
		ContextSize: 131072, Pricing: &types.PricingInfo{InputPrice: 0.0008, OutputPrice: 0.002, Currency: "CNY"},
	},
	{
		ModelID: "qwen-turbo", DisplayName: "通义千问-Turbo", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 131072, Pricing: &types.PricingInfo{InputPrice: 0.0003, OutputPrice: 0.0006, Currency: "CNY"},
	},
	{
		ModelID: "text-embedding-v3", DisplayName: "通用文本向量-v3", ModelType: types.ModelTypeEmbedding,
		Capabilities: []types.ModelCapability{types.CapabilityEmbedding},
		ContextSize: 8192, Pricing: &types.PricingInfo{InputPrice: 0.0007, Currency: "CNY"},
	},
	{
		ModelID: "gte-rerank", DisplayName: "GTE Rerank", ModelType: types.ModelTypeRerank,
		Capabilities: []types.ModelCapability{types.CapabilityRerank},
		Pricing: &types.PricingInfo{InputPrice: 0.001, Currency: "CNY"},
	},
}


// 智谱 AI 预置模型
var ZhipuPresetModels = []PresetModel{
	{
		ModelID: "glm-4-plus", DisplayName: "GLM-4-Plus", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityFunctionCall},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 0.05, OutputPrice: 0.05, Currency: "CNY"},
	},
	{
		ModelID: "glm-4-air", DisplayName: "GLM-4-Air", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 0.001, OutputPrice: 0.001, Currency: "CNY"},
	},
	{
		ModelID: "glm-4-flash", DisplayName: "GLM-4-Flash", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 128000, Pricing: &types.PricingInfo{InputPrice: 0.0001, OutputPrice: 0.0001, Currency: "CNY"},
	},
	{
		ModelID: "embedding-3", DisplayName: "Embedding-3", ModelType: types.ModelTypeEmbedding,
		Capabilities: []types.ModelCapability{types.CapabilityEmbedding},
		ContextSize: 8192, Pricing: &types.PricingInfo{InputPrice: 0.0005, Currency: "CNY"},
	},
}

// DeepSeek 预置模型
var DeepSeekPresetModels = []PresetModel{
	{
		ModelID: "deepseek-chat", DisplayName: "DeepSeek Chat", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat, types.CapabilityFunctionCall},
		ContextSize: 64000, Pricing: &types.PricingInfo{InputPrice: 0.001, OutputPrice: 0.002, Currency: "CNY"},
	},
	{
		ModelID: "deepseek-reasoner", DisplayName: "DeepSeek Reasoner (R1)", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 64000, Pricing: &types.PricingInfo{InputPrice: 0.004, OutputPrice: 0.016, Currency: "CNY"},
	},
}

// SiliconFlow 预置模型
var SiliconFlowPresetModels = []PresetModel{
	{
		ModelID: "deepseek-ai/DeepSeek-V3", DisplayName: "DeepSeek-V3", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 64000, Pricing: &types.PricingInfo{InputPrice: 0.001, OutputPrice: 0.002, Currency: "CNY"},
	},
	{
		ModelID: "Qwen/Qwen2.5-72B-Instruct", DisplayName: "Qwen2.5-72B", ModelType: types.ModelTypeKnowledgeQA,
		Capabilities: []types.ModelCapability{types.CapabilityChat},
		ContextSize: 32000, Pricing: &types.PricingInfo{InputPrice: 0.004, OutputPrice: 0.004, Currency: "CNY"},
	},
	{
		ModelID: "BAAI/bge-m3", DisplayName: "BGE-M3", ModelType: types.ModelTypeEmbedding,
		Capabilities: []types.ModelCapability{types.CapabilityEmbedding},
		ContextSize: 8192, Pricing: &types.PricingInfo{InputPrice: 0.0001, Currency: "CNY"},
	},
	{
		ModelID: "BAAI/bge-reranker-v2-m3", DisplayName: "BGE Reranker v2 M3", ModelType: types.ModelTypeRerank,
		Capabilities: []types.ModelCapability{types.CapabilityRerank},
		Pricing: &types.PricingInfo{InputPrice: 0.0001, Currency: "CNY"},
	},
}

// GetPresetModels 获取指定厂商的预置模型列表
func GetPresetModels(provider ProviderName) []PresetModel {
	switch provider {
	case ProviderOpenAI:
		return OpenAIPresetModels
	case ProviderAliyun:
		return AliyunPresetModels
	case ProviderZhipu:
		return ZhipuPresetModels
	case ProviderDeepSeek:
		return DeepSeekPresetModels
	case ProviderSiliconFlow:
		return SiliconFlowPresetModels
	default:
		return nil
	}
}

// GetPresetModelsByType 获取指定厂商和类型的预置模型
func GetPresetModelsByType(provider ProviderName, modelType types.ModelType) []PresetModel {
	all := GetPresetModels(provider)
	var result []PresetModel
	for _, m := range all {
		if m.ModelType == modelType {
			result = append(result, m)
		}
	}
	return result
}
