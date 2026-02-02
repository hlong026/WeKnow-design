package types

import (
	"time"
)

// ModelUsageStats 模型使用统计
type ModelUsageStats struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID     uint64    `json:"tenant_id" gorm:"index"`
	ModelID      string    `json:"model_id" gorm:"type:varchar(36);index"`
	CredentialID string    `json:"credential_id" gorm:"type:varchar(36)"`
	Date         time.Time `json:"date" gorm:"type:date;index"`
	RequestCount int       `json:"request_count" gorm:"default:0"`
	InputTokens  int64     `json:"input_tokens" gorm:"default:0"`
	OutputTokens int64     `json:"output_tokens" gorm:"default:0"`
	TotalCost    float64   `json:"total_cost" gorm:"type:decimal(10,4);default:0"`
	ErrorCount   int       `json:"error_count" gorm:"default:0"`
	AvgLatencyMs int       `json:"avg_latency_ms" gorm:"default:0"`
}

func (ModelUsageStats) TableName() string {
	return "model_usage_stats"
}

// UsageSummary 使用统计汇总
type UsageSummary struct {
	TotalRequests  int64   `json:"total_requests"`
	TotalInputTokens  int64   `json:"total_input_tokens"`
	TotalOutputTokens int64   `json:"total_output_tokens"`
	TotalCost      float64 `json:"total_cost"`
	TotalErrors    int64   `json:"total_errors"`
	AvgLatencyMs   int     `json:"avg_latency_ms"`
}

// DailyUsage 每日使用统计
type DailyUsage struct {
	Date         string  `json:"date"`
	RequestCount int     `json:"request_count"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	TotalCost    float64 `json:"total_cost"`
}
