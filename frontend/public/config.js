// 运行时配置（本地开发默认值，Docker 环境会被 entrypoint 脚本覆盖）
window.__RUNTIME_CONFIG__ = {
  MAX_FILE_SIZE_MB: 50,
  DISABLE_AUTH: true,  // 是否禁用用户认证
  HIDE_OLLAMA: false,   // 是否隐藏 Ollama 相关设置
  PROTECTED_BUILTIN_AGENTS: "builtin-quick-answer"  // 受保护的内置智能体（逗号分隔）
};
