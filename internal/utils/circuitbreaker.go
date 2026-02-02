package utils

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
)

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int

const (
	// StateClosed 关闭状态（正常工作）
	StateClosed CircuitBreakerState = iota
	// StateOpen 打开状态（熔断中）
	StateOpen
	// StateHalfOpen 半开状态（尝试恢复）
	StateHalfOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	// Name 熔断器名称
	Name string
	// MaxFailures 触发熔断的最大失败次数
	MaxFailures int
	// Timeout 熔断超时时间（熔断后多久尝试恢复）
	Timeout time.Duration
	// MaxHalfOpenRequests 半开状态下允许的最大请求数
	MaxHalfOpenRequests int
	// OnStateChange 状态变化回调
	OnStateChange func(name string, from, to CircuitBreakerState)
}

// DefaultCircuitBreakerConfig 默认配置
func DefaultCircuitBreakerConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:                name,
		MaxFailures:         5,
		Timeout:             30 * time.Second,
		MaxHalfOpenRequests: 3,
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	config CircuitBreakerConfig

	mu                  sync.RWMutex
	state               CircuitBreakerState
	failures            int
	successes           int
	lastFailureTime     time.Time
	halfOpenRequests    int
}

// ErrCircuitBreakerOpen 熔断器打开错误
var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.MaxFailures <= 0 {
		config.MaxFailures = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxHalfOpenRequests <= 0 {
		config.MaxHalfOpenRequests = 3
	}

	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute 执行函数，带熔断保护
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if !cb.allowRequest() {
		logger.Warnf(ctx, "[CircuitBreaker] %s: request rejected, circuit is open", cb.config.Name)
		return ErrCircuitBreakerOpen
	}

	err := fn()

	cb.recordResult(ctx, err)
	return err
}

// allowRequest 检查是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// 检查是否超过熔断超时时间
		if time.Since(cb.lastFailureTime) > cb.config.Timeout {
			cb.toHalfOpen()
			return true
		}
		return false

	case StateHalfOpen:
		// 半开状态下限制请求数
		if cb.halfOpenRequests < cb.config.MaxHalfOpenRequests {
			cb.halfOpenRequests++
			return true
		}
		return false

	default:
		return false
	}
}

// recordResult 记录请求结果
func (cb *CircuitBreaker) recordResult(ctx context.Context, err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure(ctx)
	} else {
		cb.onSuccess(ctx)
	}
}

// onSuccess 处理成功
func (cb *CircuitBreaker) onSuccess(ctx context.Context) {
	switch cb.state {
	case StateClosed:
		// 重置失败计数
		cb.failures = 0

	case StateHalfOpen:
		cb.successes++
		// 如果半开状态下成功次数达到阈值，关闭熔断器
		if cb.successes >= cb.config.MaxHalfOpenRequests {
			cb.toClosed(ctx)
		}
	}
}

// onFailure 处理失败
func (cb *CircuitBreaker) onFailure(ctx context.Context) {
	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		// 如果失败次数达到阈值，打开熔断器
		if cb.failures >= cb.config.MaxFailures {
			cb.toOpen(ctx)
		}

	case StateHalfOpen:
		// 半开状态下失败，重新打开熔断器
		cb.toOpen(ctx)
	}
}

// toOpen 切换到打开状态
func (cb *CircuitBreaker) toOpen(ctx context.Context) {
	if cb.state == StateOpen {
		return
	}

	oldState := cb.state
	cb.state = StateOpen
	cb.halfOpenRequests = 0
	cb.successes = 0

	logger.Warnf(ctx, "[CircuitBreaker] %s: state changed from %s to open, failures=%d",
		cb.config.Name, oldState.String(), cb.failures)

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, oldState, StateOpen)
	}
}

// toHalfOpen 切换到半开状态
func (cb *CircuitBreaker) toHalfOpen() {
	if cb.state == StateHalfOpen {
		return
	}

	oldState := cb.state
	cb.state = StateHalfOpen
	cb.halfOpenRequests = 0
	cb.successes = 0
	cb.failures = 0

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, oldState, StateHalfOpen)
	}
}

// toClosed 切换到关闭状态
func (cb *CircuitBreaker) toClosed(ctx context.Context) {
	if cb.state == StateClosed {
		return
	}

	oldState := cb.state
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0

	logger.Infof(ctx, "[CircuitBreaker] %s: state changed from %s to closed, circuit recovered",
		cb.config.Name, oldState.String())

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, oldState, StateClosed)
	}
}

// State 获取当前状态
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats 获取统计信息
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":              cb.config.Name,
		"state":             cb.state.String(),
		"failures":          cb.failures,
		"successes":         cb.successes,
		"half_open_requests": cb.halfOpenRequests,
		"last_failure_time": cb.lastFailureTime,
	}
}

// CircuitBreakerRegistry 熔断器注册表
type CircuitBreakerRegistry struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
}

// NewCircuitBreakerRegistry 创建熔断器注册表
func NewCircuitBreakerRegistry() *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Get 获取或创建熔断器
func (r *CircuitBreakerRegistry) Get(name string) *CircuitBreaker {
	r.mu.RLock()
	if cb, ok := r.breakers[name]; ok {
		r.mu.RUnlock()
		return cb
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// 双重检查
	if cb, ok := r.breakers[name]; ok {
		return cb
	}

	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig(name))
	r.breakers[name] = cb
	return cb
}

// GetWithConfig 获取或创建带配置的熔断器
func (r *CircuitBreakerRegistry) GetWithConfig(config CircuitBreakerConfig) *CircuitBreaker {
	r.mu.RLock()
	if cb, ok := r.breakers[config.Name]; ok {
		r.mu.RUnlock()
		return cb
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// 双重检查
	if cb, ok := r.breakers[config.Name]; ok {
		return cb
	}

	cb := NewCircuitBreaker(config)
	r.breakers[config.Name] = cb
	return cb
}

// AllStats 获取所有熔断器统计
func (r *CircuitBreakerRegistry) AllStats() map[string]map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for name, cb := range r.breakers {
		stats[name] = cb.Stats()
	}
	return stats
}

// 全局熔断器注册表
var globalCircuitBreakerRegistry = NewCircuitBreakerRegistry()

// GetCircuitBreaker 获取全局熔断器
func GetCircuitBreaker(name string) *CircuitBreaker {
	return globalCircuitBreakerRegistry.Get(name)
}

// GetCircuitBreakerWithConfig 获取带配置的全局熔断器
func GetCircuitBreakerWithConfig(config CircuitBreakerConfig) *CircuitBreaker {
	return globalCircuitBreakerRegistry.GetWithConfig(config)
}

// GetAllCircuitBreakerStats 获取所有熔断器统计
func GetAllCircuitBreakerStats() map[string]map[string]interface{} {
	return globalCircuitBreakerRegistry.AllStats()
}
