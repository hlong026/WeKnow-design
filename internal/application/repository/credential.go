package repository

import (
	"context"
	"errors"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"gorm.io/gorm"
)

// credentialRepository 凭证仓库实现
type credentialRepository struct {
	db *gorm.DB
}

// NewCredentialRepository 创建凭证仓库
func NewCredentialRepository(db *gorm.DB) interfaces.CredentialRepository {
	return &credentialRepository{db: db}
}

// Create 创建凭证
func (r *credentialRepository) Create(ctx context.Context, credential *types.ProviderCredential) error {
	return r.db.WithContext(ctx).Create(credential).Error
}

// GetByID 根据ID获取凭证
func (r *credentialRepository) GetByID(ctx context.Context, tenantID uint64, id string) (*types.ProviderCredential, error) {
	var credential types.ProviderCredential
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &credential, nil
}

// List 获取凭证列表
func (r *credentialRepository) List(ctx context.Context, tenantID uint64, provider string) ([]*types.ProviderCredential, error) {
	var credentials []*types.ProviderCredential
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	
	if err := query.Order("created_at DESC").Find(&credentials).Error; err != nil {
		return nil, err
	}
	return credentials, nil
}

// Update 更新凭证
func (r *credentialRepository) Update(ctx context.Context, credential *types.ProviderCredential) error {
	return r.db.WithContext(ctx).Model(&types.ProviderCredential{}).
		Where("id = ? AND tenant_id = ?", credential.ID, credential.TenantID).
		Select("*").Updates(credential).Error
}

// Delete 删除凭证
func (r *credentialRepository) Delete(ctx context.Context, tenantID uint64, id string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&types.ProviderCredential{}).Error
}

// GetDefault 获取指定厂商的默认凭证
func (r *credentialRepository) GetDefault(ctx context.Context, tenantID uint64, provider string) (*types.ProviderCredential, error) {
	var credential types.ProviderCredential
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND provider = ? AND is_default = ?", tenantID, provider, true).
		First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &credential, nil
}

// ClearDefault 清除指定厂商的默认凭证标记
func (r *credentialRepository) ClearDefault(ctx context.Context, tenantID uint64, provider string, excludeID string) error {
	query := r.db.WithContext(ctx).Model(&types.ProviderCredential{}).
		Where("tenant_id = ? AND provider = ? AND is_default = ?", tenantID, provider, true)
	
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	
	return query.Update("is_default", false).Error
}

// GetByProviderAndName 根据厂商和名称获取凭证
func (r *credentialRepository) GetByProviderAndName(ctx context.Context, tenantID uint64, provider, name string) (*types.ProviderCredential, error) {
	var credential types.ProviderCredential
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND provider = ? AND name = ?", tenantID, provider, name).
		First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &credential, nil
}
