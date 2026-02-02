package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CredentialStatus 凭证状态
type CredentialStatus string

const (
	CredentialStatusActive        CredentialStatus = "active"
	CredentialStatusInvalid       CredentialStatus = "invalid"
	CredentialStatusQuotaExceeded CredentialStatus = "quota_exceeded"
)

// QuotaConfig 配额配置
type QuotaConfig struct {
	DailyLimit     int64   `json:"daily_limit,omitempty"`
	MonthlyLimit   int64   `json:"monthly_limit,omitempty"`
	TokenLimit     int64   `json:"token_limit,omitempty"`
	AlertThreshold float64 `json:"alert_threshold,omitempty"`
}

// Value implements driver.Valuer for QuotaConfig
func (q QuotaConfig) Value() (driver.Value, error) {
	return json.Marshal(q)
}

// Scan implements sql.Scanner for QuotaConfig
func (q *QuotaConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, q)
}

// ProviderCredential 厂商凭证
type ProviderCredential struct {
	ID          string            `json:"id" gorm:"type:varchar(36);primaryKey"`
	TenantID    uint64            `json:"tenant_id" gorm:"index"`
	Provider    string            `json:"provider" gorm:"type:varchar(50);index"`
	Name        string            `json:"name" gorm:"type:varchar(255)"`
	Credentials map[string]string `json:"credentials" gorm:"type:jsonb;serializer:json"`
	BaseURL     string            `json:"base_url" gorm:"type:varchar(500)"`
	IsDefault   bool              `json:"is_default" gorm:"default:false"`
	Status      CredentialStatus  `json:"status" gorm:"type:varchar(50);default:'active'"`
	QuotaConfig *QuotaConfig      `json:"quota_config,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index"`
}

func (ProviderCredential) TableName() string {
	return "provider_credentials"
}

// BeforeCreate GORM hook - 自动生成UUID
func (c *ProviderCredential) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Status == "" {
		c.Status = CredentialStatusActive
	}
	return nil
}

// HideSensitiveInfo 隐藏敏感信息
func (c *ProviderCredential) HideSensitiveInfo() *ProviderCredential {
	copy := *c
	copy.Credentials = make(map[string]string)
	for k, v := range c.Credentials {
		if k == "api_key" || k == "secret_key" {
			if len(v) > 8 {
				copy.Credentials[k] = v[:4] + "****" + v[len(v)-4:]
			} else {
				copy.Credentials[k] = "****"
			}
		} else {
			copy.Credentials[k] = v
		}
	}
	return &copy
}
