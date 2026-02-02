package stream

import (
	"context"
	"sync"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// 默认配置常量
const (
	// DefaultCleanupInterval 默认清理间隔
	DefaultCleanupInterval = 5 * time.Minute
	// DefaultMaxAge 默认数据最大保留时间
	DefaultMaxAge = 30 * time.Minute
	// DefaultMaxEventsPerStream 每个流的最大事件数
	DefaultMaxEventsPerStream = 10000
)

// memoryStreamData holds stream events in memory
type memoryStreamData struct {
	events      []interfaces.StreamEvent
	lastUpdated time.Time
	createdAt   time.Time
	mu          sync.RWMutex
}

// MemoryStreamManager implements StreamManager using in-memory storage
type MemoryStreamManager struct {
	// Map: sessionID -> messageID -> stream data
	streams         map[string]map[string]*memoryStreamData
	mu              sync.RWMutex
	cleanupInterval time.Duration
	maxAge          time.Duration
	stopCleanup     chan struct{}
	cleanupStopped  chan struct{}
}

// MemoryStreamManagerOption 配置选项函数类型
type MemoryStreamManagerOption func(*MemoryStreamManager)

// WithCleanupInterval 设置清理间隔
func WithCleanupInterval(interval time.Duration) MemoryStreamManagerOption {
	return func(m *MemoryStreamManager) {
		m.cleanupInterval = interval
	}
}

// WithMaxAge 设置数据最大保留时间
func WithMaxAge(maxAge time.Duration) MemoryStreamManagerOption {
	return func(m *MemoryStreamManager) {
		m.maxAge = maxAge
	}
}

// NewMemoryStreamManager creates a new in-memory stream manager
func NewMemoryStreamManager(opts ...MemoryStreamManagerOption) *MemoryStreamManager {
	m := &MemoryStreamManager{
		streams:         make(map[string]map[string]*memoryStreamData),
		cleanupInterval: DefaultCleanupInterval,
		maxAge:          DefaultMaxAge,
		stopCleanup:     make(chan struct{}),
		cleanupStopped:  make(chan struct{}),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(m)
	}

	// 启动后台清理协程
	go m.startCleanupRoutine()

	return m
}

// startCleanupRoutine 启动定期清理过期数据的后台协程
func (m *MemoryStreamManager) startCleanupRoutine() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	defer close(m.cleanupStopped)

	ctx := context.Background()
	logger.Infof(ctx, "[MemoryStreamManager] Cleanup routine started, interval: %v, maxAge: %v",
		m.cleanupInterval, m.maxAge)

	for {
		select {
		case <-m.stopCleanup:
			logger.Info(ctx, "[MemoryStreamManager] Cleanup routine stopped")
			return
		case <-ticker.C:
			m.cleanup()
		}
	}
}

// cleanup 清理过期的流数据
func (m *MemoryStreamManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	now := time.Now()
	expiredStreams := 0
	expiredMessages := 0

	for sessionID, messages := range m.streams {
		for msgID, data := range messages {
			data.mu.RLock()
			isExpired := now.Sub(data.lastUpdated) > m.maxAge
			data.mu.RUnlock()

			if isExpired {
				delete(messages, msgID)
				expiredMessages++
			}
		}

		// 如果 session 下没有消息了，删除整个 session
		if len(messages) == 0 {
			delete(m.streams, sessionID)
			expiredStreams++
		}
	}

	if expiredStreams > 0 || expiredMessages > 0 {
		logger.Infof(ctx, "[MemoryStreamManager] Cleanup completed: removed %d sessions, %d messages, remaining sessions: %d",
			expiredStreams, expiredMessages, len(m.streams))
	}
}

// Close 关闭 MemoryStreamManager，停止清理协程
func (m *MemoryStreamManager) Close() error {
	close(m.stopCleanup)
	<-m.cleanupStopped
	return nil
}

// Stats 返回当前内存使用统计
func (m *MemoryStreamManager) Stats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalMessages := 0
	totalEvents := 0
	for _, messages := range m.streams {
		totalMessages += len(messages)
		for _, data := range messages {
			data.mu.RLock()
			totalEvents += len(data.events)
			data.mu.RUnlock()
		}
	}

	return map[string]interface{}{
		"sessions":       len(m.streams),
		"total_messages": totalMessages,
		"total_events":   totalEvents,
	}
}

// getOrCreateStream gets or creates stream data
func (m *MemoryStreamManager) getOrCreateStream(sessionID, messageID string) *memoryStreamData {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.streams[sessionID]; !exists {
		m.streams[sessionID] = make(map[string]*memoryStreamData)
	}

	if _, exists := m.streams[sessionID][messageID]; !exists {
		now := time.Now()
		m.streams[sessionID][messageID] = &memoryStreamData{
			events:      make([]interfaces.StreamEvent, 0),
			lastUpdated: now,
			createdAt:   now,
		}
	}

	return m.streams[sessionID][messageID]
}

// getStream gets existing stream data (returns nil if not found)
func (m *MemoryStreamManager) getStream(sessionID, messageID string) *memoryStreamData {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if sessionMap, exists := m.streams[sessionID]; exists {
		return sessionMap[messageID]
	}
	return nil
}

// AppendEvent appends a single event to the stream
func (m *MemoryStreamManager) AppendEvent(
	ctx context.Context,
	sessionID, messageID string,
	event interfaces.StreamEvent,
) error {
	stream := m.getOrCreateStream(sessionID, messageID)

	stream.mu.Lock()
	defer stream.mu.Unlock()

	// Set timestamp if not already set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// 防止单个流事件过多导致内存溢出
	if len(stream.events) >= DefaultMaxEventsPerStream {
		// 保留后半部分事件，丢弃旧事件
		halfSize := DefaultMaxEventsPerStream / 2
		stream.events = stream.events[halfSize:]
		logger.Warnf(ctx, "[MemoryStreamManager] Stream %s/%s exceeded max events, truncated to %d",
			sessionID, messageID, len(stream.events))
	}

	// Append event
	stream.events = append(stream.events, event)
	stream.lastUpdated = time.Now()

	return nil
}

// GetEvents gets events starting from offset
// Returns: events slice, next offset, error
func (m *MemoryStreamManager) GetEvents(
	ctx context.Context,
	sessionID, messageID string,
	fromOffset int,
) ([]interfaces.StreamEvent, int, error) {
	stream := m.getStream(sessionID, messageID)
	if stream == nil {
		// Stream doesn't exist yet
		return []interfaces.StreamEvent{}, fromOffset, nil
	}

	stream.mu.RLock()
	defer stream.mu.RUnlock()

	// Check if offset is beyond current events
	if fromOffset >= len(stream.events) {
		return []interfaces.StreamEvent{}, fromOffset, nil
	}

	// Get events from offset to end
	events := stream.events[fromOffset:]
	nextOffset := len(stream.events)

	// Return copy of events to avoid race conditions
	eventsCopy := make([]interfaces.StreamEvent, len(events))
	copy(eventsCopy, events)

	return eventsCopy, nextOffset, nil
}

// Ensure MemoryStreamManager implements StreamManager interface
var _ interfaces.StreamManager = (*MemoryStreamManager)(nil)
