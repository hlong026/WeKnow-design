package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ModelCapability 模型能力
type ModelCapability string

const (
	CapabilityChat         ModelCapability = "chat"
	CapabilityCompletion   ModelCapability = "completion"
	CapabilityEmbedding    ModelCapability = "embedding"
	CapabilityRerank       ModelCapability = "rerank"
	CapabilityVision       ModelCapability = "vision"
	CapabilityFunctionCall ModelCapability = "function_call"
	CapabilityStreaming    ModelCapability = "streaming"
	CapabilityJSON         ModelCapability = "json_mode"
)

// PricingInfo 价格信息
type PricingInfo struct {
	InputPrice  float64 `json:"input_price"`
	OutputPrice float64 `json:"output_price"`
	Currency    string  `json:"currency"`
}

// Value implements driver.Valuer for PricingInfo
func (p PricingInfo) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements sql.Scanner for PricingInfo
func (p *PricingInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, p)
}

// ModelConfiguration 模型配置 (重构后的模型表)
type ModelConfiguration struct {
	ID           string            `json:"id" gorm:"type:varchar(36);primaryKey"`
	TenantID     uint64            `json:"tenant_id" gorm:"index"`
	CredentialID string            `json:"credential_id" gorm:"type:varchar(36);index"`
	Provider     string            `json:"provider" gorm:"type:varchar(50);index"`
	ModelID      string            `json:"model_id" gorm:"type:varchar(255)"`
	DisplayName  string            `json:"display_name" gorm:"type:varchar(255)"`
	ModelType    ModelType         `json:"model_type" gorm:"type:varchar(50);index"`
	Capabilities []ModelCapability `json:"capabilities" gorm:"type:jsonb;serializer:json"`
	Parameters   ModelParameters   `json:"parameters" gorm:"type:jsonb"`
	PricingInfo  *PricingInfo      `json:"pricing_info" gorm:"type:jsonb"`
	IsEnabled    bool              `json:"is_enabled" gorm:"default:true"`
	IsDefault    bool              `json:"is_default" gorm:"default:false"`
	IsBuiltin    bool              `json:"is_builtin" gorm:"default:false"`
	Tags         []string          `json:"tags" gorm:"type:jsonb;serializer:json"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"deleted_at" gorm:"index"`

	// 关联
	Credential *ProviderCredential `json:"credential,omitempty" gorm:"foreignKey:CredentialID"`
}

func (ModelConfiguration) TableName() string {
	return "model_configurations"
}
