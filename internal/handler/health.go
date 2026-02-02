package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
	neo4jDriver neo4j.Driver
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB, redisClient *redis.Client, neo4jDriver neo4j.Driver) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		neo4jDriver: neo4jDriver,
	}
}

// HealthStatus 健康状态
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// ComponentHealth 组件健康状态
type ComponentHealth struct {
	Status  HealthStatus `json:"status"`
	Latency string       `json:"latency,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status     HealthStatus               `json:"status"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
}

// HealthCheck godoc
// @Summary      健康检查
// @Description  检查服务及其依赖组件的健康状态
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse  "服务健康"
// @Failure      503  {object}  HealthResponse  "服务不健康"
// @Router       /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	components := make(map[string]ComponentHealth)
	overallStatus := HealthStatusHealthy

	// 检查数据库
	dbHealth := h.checkDatabase(ctx)
	components["database"] = dbHealth
	if dbHealth.Status == HealthStatusUnhealthy {
		overallStatus = HealthStatusUnhealthy
	}

	// 检查 Redis（如果配置了）
	if h.redisClient != nil {
		redisHealth := h.checkRedis(ctx)
		components["redis"] = redisHealth
		if redisHealth.Status == HealthStatusUnhealthy {
			// Redis 不可用时降级为 degraded，因为可能使用内存流管理器
			if overallStatus == HealthStatusHealthy {
				overallStatus = HealthStatusDegraded
			}
		}
	}

	// 检查 Neo4j（如果启用）
	if h.neo4jDriver != nil && os.Getenv("NEO4J_ENABLE") == "true" {
		neo4jHealth := h.checkNeo4j(ctx)
		components["neo4j"] = neo4jHealth
		if neo4jHealth.Status == HealthStatusUnhealthy {
			// Neo4j 不可用时降级为 degraded
			if overallStatus == HealthStatusHealthy {
				overallStatus = HealthStatusDegraded
			}
		}
	}

	response := HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: components,
	}

	// 根据状态返回不同的 HTTP 状态码
	httpStatus := http.StatusOK
	if overallStatus == HealthStatusUnhealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, response)
}

// LivenessCheck 存活检查（Kubernetes liveness probe）
// 只检查服务本身是否存活，不检查依赖
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck 就绪检查（Kubernetes readiness probe）
// 检查服务是否准备好接收流量
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// 只检查关键依赖：数据库
	dbHealth := h.checkDatabase(ctx)
	if dbHealth.Status == HealthStatusUnhealthy {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not_ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"reason":    "database unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentHealth {
	if h.db == nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  "database not configured",
		}
	}

	start := time.Now()
	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  err.Error(),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  err.Error(),
		}
	}

	return ComponentHealth{
		Status:  HealthStatusHealthy,
		Latency: time.Since(start).String(),
	}
}

// checkRedis 检查 Redis 连接
func (h *HealthHandler) checkRedis(ctx context.Context) ComponentHealth {
	if h.redisClient == nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  "redis not configured",
		}
	}

	start := time.Now()
	if _, err := h.redisClient.Ping(ctx).Result(); err != nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  err.Error(),
		}
	}

	return ComponentHealth{
		Status:  HealthStatusHealthy,
		Latency: time.Since(start).String(),
	}
}

// checkNeo4j 检查 Neo4j 连接
func (h *HealthHandler) checkNeo4j(ctx context.Context) ComponentHealth {
	if h.neo4jDriver == nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  "neo4j not configured",
		}
	}

	start := time.Now()
	if err := h.neo4jDriver.VerifyConnectivity(ctx); err != nil {
		return ComponentHealth{
			Status: HealthStatusUnhealthy,
			Error:  err.Error(),
		}
	}

	return ComponentHealth{
		Status:  HealthStatusHealthy,
		Latency: time.Since(start).String(),
	}
}
