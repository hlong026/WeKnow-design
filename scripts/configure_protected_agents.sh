#!/bin/bash

# 配置受保护的内置智能体
# 用法: ./scripts/configure_protected_agents.sh [default|all|none|custom]

set -e

ENV_FILE=".env"

if [ ! -f "$ENV_FILE" ]; then
    echo "错误: .env 文件不存在"
    echo "请先复制 .env.example 为 .env"
    exit 1
fi

case "$1" in
    default)
        # 默认：只保护快速问答
        if grep -q "^PROTECTED_BUILTIN_AGENTS=" "$ENV_FILE"; then
            sed -i.bak 's/^PROTECTED_BUILTIN_AGENTS=.*/PROTECTED_BUILTIN_AGENTS=builtin-quick-answer/' "$ENV_FILE"
        else
            echo "PROTECTED_BUILTIN_AGENTS=builtin-quick-answer" >> "$ENV_FILE"
        fi
        echo "✓ 已设置为默认配置（只保护快速问答）"
        echo ""
        echo "受保护的智能体:"
        echo "  - builtin-quick-answer (快速问答)"
        ;;
    all)
        # 保护所有内置智能体
        if grep -q "^PROTECTED_BUILTIN_AGENTS=" "$ENV_FILE"; then
            sed -i.bak 's/^PROTECTED_BUILTIN_AGENTS=.*/PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning,builtin-data-analyst/' "$ENV_FILE"
        else
            echo "PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning,builtin-data-analyst" >> "$ENV_FILE"
        fi
        echo "✓ 已保护所有内置智能体"
        echo ""
        echo "受保护的智能体:"
        echo "  - builtin-quick-answer (快速问答)"
        echo "  - builtin-smart-reasoning (智能推理)"
        echo "  - builtin-data-analyst (数据分析师)"
        ;;
    none)
        # 允许删除所有智能体
        if grep -q "^PROTECTED_BUILTIN_AGENTS=" "$ENV_FILE"; then
            sed -i.bak 's/^PROTECTED_BUILTIN_AGENTS=.*/PROTECTED_BUILTIN_AGENTS=/' "$ENV_FILE"
        else
            echo "PROTECTED_BUILTIN_AGENTS=" >> "$ENV_FILE"
        fi
        echo "✓ 已允许删除所有智能体"
        echo ""
        echo "⚠️  警告：所有内置智能体都可以被删除"
        echo "⚠️  请谨慎使用此配置"
        ;;
    custom)
        echo "自定义配置"
        echo ""
        echo "请手动编辑 .env 文件中的 PROTECTED_BUILTIN_AGENTS 变量"
        echo ""
        echo "可用的智能体 ID:"
        echo "  - builtin-quick-answer      (快速问答)"
        echo "  - builtin-smart-reasoning   (智能推理)"
        echo "  - builtin-data-analyst      (数据分析师)"
        echo ""
        echo "示例:"
        echo "  PROTECTED_BUILTIN_AGENTS=builtin-quick-answer,builtin-smart-reasoning"
        exit 0
        ;;
    status)
        if grep -q "^PROTECTED_BUILTIN_AGENTS=" "$ENV_FILE"; then
            value=$(grep "^PROTECTED_BUILTIN_AGENTS=" "$ENV_FILE" | cut -d'=' -f2)
            if [ -z "$value" ]; then
                echo "当前配置: 允许删除所有智能体"
            else
                echo "当前配置: 受保护的智能体"
                IFS=',' read -ra AGENTS <<< "$value"
                for agent in "${AGENTS[@]}"; do
                    agent=$(echo "$agent" | xargs)  # trim whitespace
                    case "$agent" in
                        builtin-quick-answer)
                            echo "  - builtin-quick-answer (快速问答)"
                            ;;
                        builtin-smart-reasoning)
                            echo "  - builtin-smart-reasoning (智能推理)"
                            ;;
                        builtin-data-analyst)
                            echo "  - builtin-data-analyst (数据分析师)"
                            ;;
                        *)
                            echo "  - $agent"
                            ;;
                    esac
                done
            fi
        else
            echo "当前配置: 未配置 (默认保护快速问答)"
        fi
        exit 0
        ;;
    *)
        echo "用法: $0 [default|all|none|custom|status]"
        echo ""
        echo "命令:"
        echo "  default  - 只保护快速问答（默认配置）"
        echo "  all      - 保护所有内置智能体"
        echo "  none     - 允许删除所有智能体（谨慎使用）"
        echo "  custom   - 显示自定义配置说明"
        echo "  status   - 查看当前配置"
        echo ""
        echo "示例:"
        echo "  $0 default  # 设置为默认配置"
        echo "  $0 all      # 保护所有内置智能体"
        echo "  $0 status   # 查看当前配置"
        exit 1
        ;;
esac

echo ""
echo "提示: 需要重启服务才能生效"
echo "运行: ./scripts/start_all.sh --stop && ./scripts/start_all.sh"
