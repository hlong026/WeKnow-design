package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

// CredentialService 凭证服务接口
type CredentialService interface {
	// CreateCredential 创建凭证
	CreateCredential(ctx context.Context, credential *types.ProviderCredential) error
	// GetCredentialByID 根据ID获取凭证
	GetCredentialByID(ctx context.Context, id string) (*types.ProviderCredential, error)
	// ListCredentials 获取凭证列表
	ListCredentials(ctx context.Context, provider string) ([]*types.ProviderCredential, error)
	// UpdateCredential 更新凭证
	UpdateCredential(ctx context.Context, credential *types.ProviderCredential) error
	// DeleteCredential 删除凭证
	DeleteCredential(ctx context.Context, id string) error
	// TestCredential 测试凭证连接
	TestCredential(ctx context.Context, id string) error
	// GetDefaultCredential 获取指定厂商的默认凭证
	GetDefaultCredential(ctx context.Context, provider string) (*types.ProviderCredential, error)
}

// CredentialRepository 凭证仓库接口
type CredentialRepository interface {
	// Create 创建凭证
	Create(ctx context.Context, credential *types.ProviderCredential) error
	// GetByID 根据ID获取凭证
	GetByID(ctx context.Context, tenantID uint64, id string) (*types.ProviderCredential, error)
	// List 获取凭证列表
	List(ctx context.Context, tenantID uint64, provider string) ([]*types.ProviderCredential, error)
	// Update 更新凭证
	Update(ctx context.Context, credential *types.ProviderCredential) error
	// Delete 删除凭证
	Delete(ctx context.Context, tenantID uint64, id string) error
	// GetDefault 获取指定厂商的默认凭证
	GetDefault(ctx context.Context, tenantID uint64, provider string) (*types.ProviderCredential, error)
	// ClearDefault 清除指定厂商的默认凭证标记
	ClearDefault(ctx context.Context, tenantID uint64, provider string, excludeID string) error
	// GetByProviderAndName 根据厂商和名称获取凭证
	GetByProviderAndName(ctx context.Context, tenantID uint64, provider, name string) (*types.ProviderCredential, error)
}
