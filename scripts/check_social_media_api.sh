#!/bin/bash

# 检查自媒体 API 是否可用

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"

echo "========================================="
echo "检查自媒体 API 状态"
echo "========================================="
echo ""

# 1. 检查健康状态
echo "1. 检查后端健康状态..."
HEALTH_RESPONSE=$(curl -s "${API_BASE_URL}/health")
echo "响应: ${HEALTH_RESPONSE}"

if echo "${HEALTH_RESPONSE}" | grep -q '"status":"ok"'; then
    echo "✓ 后端服务正常运行"
else
    echo "✗ 后端服务未运行或异常"
    exit 1
fi
echo ""

# 2. 检查 API 路由（不需要认证的测试）
echo "2. 检查社媒 API 路由..."
echo "尝试访问: ${API_BASE_URL}/api/v1/social-media/extract"

# 发送一个空请求，看是否返回 401（需要认证）而不是 404
ROUTE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "${API_BASE_URL}/api/v1/social-media/extract" \
  -H "Content-Type: application/json" \
  -d '{}')

HTTP_STATUS=$(echo "${ROUTE_RESPONSE}" | grep "HTTP_STATUS" | cut -d':' -f2)

echo "HTTP 状态码: ${HTTP_STATUS}"

if [ "${HTTP_STATUS}" = "404" ]; then
    echo "✗ 路由不存在 (404) - 后端可能需要重新编译和启动"
    echo ""
    echo "解决方案："
    echo "1. 停止后端服务"
    echo "2. 重新编译: go build -o weknora ./cmd/server"
    echo "3. 重新启动: ./weknora"
    exit 1
elif [ "${HTTP_STATUS}" = "401" ] || [ "${HTTP_STATUS}" = "400" ]; then
    echo "✓ 路由存在 (返回 ${HTTP_STATUS}，这是正常的，因为我们没有提供认证信息)"
else
    echo "? 返回状态码: ${HTTP_STATUS}"
fi
echo ""

# 3. 检查 Swagger 文档
echo "3. 检查 Swagger 文档..."
if [ "$(curl -s -o /dev/null -w "%{http_code}" "${API_BASE_URL}/swagger/index.html")" = "200" ]; then
    echo "✓ Swagger 文档可访问: ${API_BASE_URL}/swagger/index.html"
    echo "  可以在浏览器中查看 API 文档"
else
    echo "ℹ Swagger 文档不可用（可能在生产模式下禁用）"
fi
echo ""

echo "========================================="
echo "检查完成"
echo "========================================="
