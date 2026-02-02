#!/bin/bash

# WeKnora Docker 镜像重构和重启脚本
# 用法: ./scripts/rebuild_and_restart.sh [选项]
# 选项:
#   --dev         使用开发环境配置 (docker-compose.dev.yml)
#   --prod        使用生产环境配置 (docker-compose.yml) [默认]
#   --service     只重构指定服务 (app|frontend|docreader|all) [默认: all]
#   --no-cache    不使用缓存构建
#   --profile     指定 profile (minio|neo4j|qdrant|jaeger|full)

set -e

# 默认配置
COMPOSE_FILE="docker-compose.yml"
SERVICE="all"
NO_CACHE=""
PROFILE=""
ENV_MODE="生产"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --dev)
            COMPOSE_FILE="docker-compose.dev.yml"
            ENV_MODE="开发"
            shift
            ;;
        --prod)
            COMPOSE_FILE="docker-compose.yml"
            ENV_MODE="生产"
            shift
            ;;
        --service)
            SERVICE="$2"
            shift 2
            ;;
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        --profile)
            PROFILE="--profile $2"
            shift 2
            ;;
        *)
            echo "未知选项: $1"
            echo "用法: $0 [--dev|--prod] [--service app|frontend|docreader|all] [--no-cache] [--profile minio|neo4j|qdrant|jaeger|full]"
            exit 1
            ;;
    esac
done

echo "=========================================="
echo "WeKnora Docker 镜像重构和重启"
echo "=========================================="
echo "环境模式: $ENV_MODE"
echo "配置文件: $COMPOSE_FILE"
echo "服务范围: $SERVICE"
echo "=========================================="

# 停止并删除容器
echo ""
echo "步骤 1/4: 停止现有容器..."
if [ "$SERVICE" = "all" ]; then
    docker compose -f $COMPOSE_FILE $PROFILE down
else
    docker compose -f $COMPOSE_FILE $PROFILE stop $SERVICE
    docker compose -f $COMPOSE_FILE $PROFILE rm -f $SERVICE
fi

# 删除旧镜像（可选）
echo ""
echo "步骤 2/4: 清理旧镜像..."
if [ "$SERVICE" = "all" ]; then
    echo "清理所有 WeKnora 相关镜像..."
    docker images | grep "wechatopenai/weknora" | awk '{print $3}' | xargs -r docker rmi -f || true
else
    echo "清理 $SERVICE 镜像..."
    docker images | grep "wechatopenai/weknora-$SERVICE" | awk '{print $3}' | xargs -r docker rmi -f || true
fi

# 重新构建镜像
echo ""
echo "步骤 3/4: 重新构建镜像..."
if [ "$SERVICE" = "all" ]; then
    docker compose -f $COMPOSE_FILE $PROFILE build $NO_CACHE
else
    docker compose -f $COMPOSE_FILE $PROFILE build $NO_CACHE $SERVICE
fi

# 启动服务
echo ""
echo "步骤 4/4: 启动服务..."
if [ "$SERVICE" = "all" ]; then
    docker compose -f $COMPOSE_FILE $PROFILE up -d
else
    docker compose -f $COMPOSE_FILE $PROFILE up -d $SERVICE
fi

# 显示运行状态
echo ""
echo "=========================================="
echo "服务状态:"
echo "=========================================="
docker compose -f $COMPOSE_FILE $PROFILE ps

echo ""
echo "=========================================="
echo "重构和重启完成！"
echo "=========================================="
echo ""
echo "查看日志: docker compose -f $COMPOSE_FILE logs -f [服务名]"
echo "停止服务: docker compose -f $COMPOSE_FILE down"
echo ""
