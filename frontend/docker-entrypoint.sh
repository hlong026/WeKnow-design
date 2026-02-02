#!/bin/sh

# 生成运行时配置文件，注入环境变量到前端
cat > /usr/share/nginx/html/config.js << EOF
window.__RUNTIME_CONFIG__ = {
  MAX_FILE_SIZE_MB: ${MAX_FILE_SIZE_MB:-50},
  DISABLE_AUTH: ${DISABLE_AUTH:-false},
  HIDE_OLLAMA: ${HIDE_OLLAMA:-false},
  PROTECTED_BUILTIN_AGENTS: "${PROTECTED_BUILTIN_AGENTS:-builtin-quick-answer}"
};
EOF

# 处理 nginx 配置
export MAX_FILE_SIZE=${MAX_FILE_SIZE_MB}M
envsubst '${MAX_FILE_SIZE}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

# 启动 nginx
exec nginx -g 'daemon off;'
