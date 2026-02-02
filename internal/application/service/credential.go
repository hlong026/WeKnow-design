package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// ErrCredentialNotFound 凭证未找到错误
var ErrCredentialNotFound = errors.New("credential not found")

// ErrCredentialExists 凭证已存在错误
var ErrCredentialExists = errors.New("credential with the same name already exists")

// credentialService 凭证服务实现
type credentialService struct {
	repo interfaces.CredentialRepository
}

// NewCredentialService 创建凭证服务
func NewCredentialService(repo interfaces.CredentialRepository) interfaces.CredentialService {
	return &credentialService{repo: repo}
}

// CreateCredential 创建凭证
func (s *credentialService) CreateCredential(ctx context.Context, credential *types.ProviderCredential) error {
	logger.Infof(ctx, "Creating credential: %s for provider: %s", credential.Name, credential.Provider)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)
	credential.TenantID = tenantID

	// 检查是否已存在同名凭证
	existing, err := s.repo.GetByProviderAndName(ctx, tenantID, credential.Provider, credential.Name)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"provider": credential.Provider,
			"name":     credential.Name,
		})
		return err
	}
	if existing != nil {
		return ErrCredentialExists
	}

	// 如果设置为默认，清除其他默认凭证
	if credential.IsDefault {
		if err := s.repo.ClearDefault(ctx, tenantID, credential.Provider, ""); err != nil {
			logger.ErrorWithFields(ctx, err, map[string]interface{}{
				"provider": credential.Provider,
			})
			return err
		}
	}

	// 创建凭证
	if err := s.repo.Create(ctx, credential); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"provider": credential.Provider,
			"name":     credential.Name,
		})
		return err
	}

	logger.Infof(ctx, "Credential created successfully: %s", credential.ID)
	return nil
}

// GetCredentialByID 根据ID获取凭证
func (s *credentialService) GetCredentialByID(ctx context.Context, id string) (*types.ProviderCredential, error) {
	if id == "" {
		return nil, errors.New("credential ID cannot be empty")
	}

	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)
	credential, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"credential_id": id,
		})
		return nil, err
	}
	if credential == nil {
		return nil, ErrCredentialNotFound
	}

	return credential, nil
}

// ListCredentials 获取凭证列表
func (s *credentialService) ListCredentials(ctx context.Context, provider string) ([]*types.ProviderCredential, error) {
	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)
	logger.Infof(ctx, "Listing credentials for tenant: %d, provider: %s", tenantID, provider)

	credentials, err := s.repo.List(ctx, tenantID, provider)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": tenantID,
			"provider":  provider,
		})
		return nil, err
	}

	logger.Infof(ctx, "Retrieved %d credentials", len(credentials))
	return credentials, nil
}

// UpdateCredential 更新凭证
func (s *credentialService) UpdateCredential(ctx context.Context, credential *types.ProviderCredential) error {
	logger.Infof(ctx, "Updating credential: %s", credential.ID)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)

	// 检查凭证是否存在
	existing, err := s.repo.GetByID(ctx, tenantID, credential.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCredentialNotFound
	}

	credential.TenantID = tenantID

	// 如果设置为默认，清除其他默认凭证
	if credential.IsDefault {
		if err := s.repo.ClearDefault(ctx, tenantID, credential.Provider, credential.ID); err != nil {
			logger.ErrorWithFields(ctx, err, map[string]interface{}{
				"provider": credential.Provider,
			})
			return err
		}
	}

	if err := s.repo.Update(ctx, credential); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"credential_id": credential.ID,
		})
		return err
	}

	logger.Infof(ctx, "Credential updated successfully: %s", credential.ID)
	return nil
}

// DeleteCredential 删除凭证
func (s *credentialService) DeleteCredential(ctx context.Context, id string) error {
	logger.Infof(ctx, "Deleting credential: %s", id)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)

	// 检查凭证是否存在
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCredentialNotFound
	}

	// TODO: 检查是否有模型在使用此凭证

	if err := s.repo.Delete(ctx, tenantID, id); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"credential_id": id,
		})
		return err
	}

	logger.Infof(ctx, "Credential deleted successfully: %s", id)
	return nil
}

// TestCredential 测试凭证连接
func (s *credentialService) TestCredential(ctx context.Context, id string) error {
	logger.Infof(ctx, "Testing credential: %s", id)

	credential, err := s.GetCredentialByID(ctx, id)
	if err != nil {
		return err
	}

	// 获取 provider 并验证配置
	p, ok := provider.Get(provider.ProviderName(credential.Provider))
	if !ok {
		return fmt.Errorf("unknown provider: %s", credential.Provider)
	}

	apiKey := credential.Credentials["api_key"]
	baseURL := credential.BaseURL
	if baseURL == "" {
		// 使用默认 URL
		info := p.Info()
		baseURL = info.GetDefaultURL(types.ModelTypeKnowledgeQA)
	}

	config := &provider.Config{
		Provider:  provider.ProviderName(credential.Provider),
		BaseURL:   baseURL,
		APIKey:    apiKey,
		ModelName: "test",
	}

	if err := p.ValidateConfig(config); err != nil {
		// 更新凭证状态为无效
		credential.Status = types.CredentialStatusInvalid
		s.repo.Update(ctx, credential)
		return fmt.Errorf("credential validation failed: %w", err)
	}

	// 实际调用 API 测试连接
	testErr := s.testAPIConnection(ctx, credential.Provider, baseURL, apiKey)
	if testErr != nil {
		credential.Status = types.CredentialStatusInvalid
		s.repo.Update(ctx, credential)
		return fmt.Errorf("connection test failed: %w", testErr)
	}

	// 更新凭证状态为有效
	credential.Status = types.CredentialStatusActive
	s.repo.Update(ctx, credential)

	logger.Infof(ctx, "Credential test passed: %s", id)
	return nil
}

// testAPIConnection 实际测试 API 连接
func (s *credentialService) testAPIConnection(ctx context.Context, providerName, baseURL, apiKey string) error {
	// 根据不同厂商调用对应的测试接口
	switch provider.ProviderName(providerName) {
	case provider.ProviderOpenAI:
		return s.testOpenAIConnection(ctx, baseURL, apiKey)
	case provider.ProviderAliyun:
		return s.testAliyunConnection(ctx, baseURL, apiKey)
	case provider.ProviderZhipu:
		return s.testZhipuConnection(ctx, baseURL, apiKey)
	case provider.ProviderDeepSeek:
		return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
	case provider.ProviderSiliconFlow:
		return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
	default:
		// 通用 OpenAI 兼容接口测试
		return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
	}
}

// testOpenAIConnection 测试 OpenAI API 连接
func (s *credentialService) testOpenAIConnection(ctx context.Context, baseURL, apiKey string) error {
	return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
}

// testAliyunConnection 测试阿里云 DashScope API 连接
func (s *credentialService) testAliyunConnection(ctx context.Context, baseURL, apiKey string) error {
	// DashScope 使用 OpenAI 兼容接口
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
}

// testZhipuConnection 测试智谱 API 连接
func (s *credentialService) testZhipuConnection(ctx context.Context, baseURL, apiKey string) error {
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn/api/paas/v4"
	}
	return s.testOpenAICompatibleConnection(ctx, baseURL, apiKey)
}

// testOpenAICompatibleConnection 测试 OpenAI 兼容接口连接
func (s *credentialService) testOpenAICompatibleConnection(ctx context.Context, baseURL, apiKey string) error {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// 调用 /models 接口测试连接
	url := baseURL + "/models"
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("access forbidden, check API key permissions")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API returned error status: %d", resp.StatusCode)
	}

	logger.Infof(ctx, "API connection test successful, status: %d", resp.StatusCode)
	return nil
}

// GetDefaultCredential 获取指定厂商的默认凭证
func (s *credentialService) GetDefaultCredential(ctx context.Context, providerName string) (*types.ProviderCredential, error) {
	tenantID := ctx.Value(types.TenantIDContextKey).(uint64)
	
	credential, err := s.repo.GetDefault(ctx, tenantID, providerName)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		// 如果没有默认凭证，返回第一个
		credentials, err := s.repo.List(ctx, tenantID, providerName)
		if err != nil {
			return nil, err
		}
		if len(credentials) > 0 {
			return credentials[0], nil
		}
		return nil, ErrCredentialNotFound
	}
	
	return credential, nil
}
