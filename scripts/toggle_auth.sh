#!/bin/bash

# 切换认证功能的脚本
# 用法: ./scripts/toggle_auth.sh [enable|disable]

set -e

ENV_FILE=".env"

if [ ! -f "$ENV_FILE" ]; then
    echo "错误: .env 文件不存在"
    echo "请先复制 .env.example 为 .env"
    exit 1
fi

case "$1" in
    disable)
        echo "禁用用户认证..."
        if grep -q "^DISABLE_AUTH=" "$ENV_FILE"; then
            sed -i.bak 's/^DISABLE_AUTH=.*/DISABLE_AUTH=true/' "$ENV_FILE"
        else
            echo "DISABLE_AUTH=true" >> "$ENV_FILE"
        fi
        
        # 同时隐藏 Ollama 设置
        if grep -q "^HIDE_OLLAMA=" "$ENV_FILE"; then
            sed -i.bak 's/^HIDE_OLLAMA=.*/HIDE_OLLAMA=true/' "$ENV_FILE"
        else
            echo "HIDE_OLLAMA=true" >> "$ENV_FILE"
        fi
        
        echo "✓ 已禁用用户认证"
        echo "✓ 已隐藏 Ollama 相关设置"
        echo "✓ 登录页面将被完全隐藏，直接进入主页"
        echo ""
        echo "提示: 需要重启服务才能生效"
        echo "运行: ./scripts/start_all.sh --stop && ./scripts/start_all.sh"
        ;;
    enable)
        echo "启用用户认证..."
        if grep -q "^DISABLE_AUTH=" "$ENV_FILE"; then
            sed -i.bak 's/^DISABLE_AUTH=.*/DISABLE_AUTH=false/' "$ENV_FILE"
        else
            echo "DISABLE_AUTH=false" >> "$ENV_FILE"
        fi
        
        # 同时显示 Ollama 设置
        if grep -q "^HIDE_OLLAMA=" "$ENV_FILE"; then
            sed -i.bak 's/^HIDE_OLLAMA=.*/HIDE_OLLAMA=false/' "$ENV_FILE"
        else
            echo "HIDE_OLLAMA=false" >> "$ENV_FILE"
        fi
        
        echo "✓ 已启用用户认证"
        echo "✓ 已显示 Ollama 相关设置"
        echo "提示: 需要重启服务才能生效"
        echo "运行: ./scripts/start_all.sh --stop && ./scripts/start_all.sh"
        ;;
    status)
        if grep -q "^DISABLE_AUTH=true" "$ENV_FILE"; then
            echo "当前状态: 认证已禁用"
        elif grep -q "^DISABLE_AUTH=false" "$ENV_FILE"; then
            echo "当前状态: 认证已启用"
        else
            echo "当前状态: 未配置 (默认启用认证)"
        fi
        
        if grep -q "^HIDE_OLLAMA=true" "$ENV_FILE"; then
            echo "Ollama 设置: 已隐藏"
        elif grep -q "^HIDE_OLLAMA=false" "$ENV_FILE"; then
            echo "Ollama 设置: 已显示"
        else
            echo "Ollama 设置: 未配置 (默认显示)"
        fi
        ;;
    *)
        echo "用法: $0 [enable|disable|status]"
        echo ""
        echo "命令:"
        echo "  enable   - 启用用户认证"
        echo "  disable  - 禁用用户认证"
        echo "  status   - 查看当前状态"
        echo ""
        echo "示例:"
        echo "  $0 disable  # 禁用认证，适合内网部署"
        echo "  $0 enable   # 启用认证，适合生产环境"
        echo "  $0 status   # 查看当前认证状态"
        exit 1
        ;;
esac
