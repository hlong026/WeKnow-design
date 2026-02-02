package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BuiltinAgentID constants for built-in agents
const (
	// BuiltinQuickAnswerID is the ID for the built-in quick answer (RAG) agent
	BuiltinQuickAnswerID = "builtin-quick-answer"
	// BuiltinSmartReasoningID is the ID for the built-in smart reasoning (ReAct) agent
	BuiltinSmartReasoningID = "builtin-smart-reasoning"
	// BuiltinDeepResearcherID is the ID for the built-in deep researcher agent
	BuiltinDeepResearcherID = "builtin-deep-researcher"
	// BuiltinDataAnalystID is the ID for the built-in data analyst agent
	BuiltinDataAnalystID = "builtin-data-analyst"
	// BuiltinKnowledgeGraphExpertID is the ID for the built-in knowledge graph expert agent
	BuiltinKnowledgeGraphExpertID = "builtin-knowledge-graph-expert"
	// BuiltinDocumentAssistantID is the ID for the built-in document assistant agent
	BuiltinDocumentAssistantID = "builtin-document-assistant"
	// BuiltinKnowledgeRefinerID is the ID for the built-in knowledge refiner agent
	BuiltinKnowledgeRefinerID = "builtin-knowledge-refiner"
)

// AgentMode constants for agent running mode
const (
	// AgentModeQuickAnswer is the RAG mode for quick Q&A
	AgentModeQuickAnswer = "quick-answer"
	// AgentModeSmartReasoning is the ReAct mode for multi-step reasoning
	AgentModeSmartReasoning = "smart-reasoning"
)

// CustomAgent represents a configurable AI agent (similar to GPTs)
type CustomAgent struct {
	// Unique identifier of the agent (composite primary key with TenantID)
	// For built-in agents, this is 'builtin-quick-answer' or 'builtin-smart-reasoning'
	// For custom agents, this is a UUID
	ID string `yaml:"id" json:"id" gorm:"type:varchar(36);primaryKey"`
	// Name of the agent
	Name string `yaml:"name" json:"name" gorm:"type:varchar(255);not null"`
	// Description of the agent
	Description string `yaml:"description" json:"description" gorm:"type:text"`
	// Avatar/Icon of the agent (emoji or icon name)
	Avatar string `yaml:"avatar" json:"avatar" gorm:"type:varchar(64)"`
	// Whether this is a built-in agent (normal mode / agent mode)
	IsBuiltin bool `yaml:"is_builtin" json:"is_builtin" gorm:"default:false"`
	// Tenant ID (composite primary key with ID)
	TenantID uint64 `yaml:"tenant_id" json:"tenant_id" gorm:"primaryKey"`
	// Created by user ID
	CreatedBy string `yaml:"created_by" json:"created_by" gorm:"type:varchar(36)"`

	// Agent configuration
	Config CustomAgentConfig `yaml:"config" json:"config" gorm:"type:json"`

	// Timestamps
	CreatedAt time.Time      `yaml:"created_at" json:"created_at"`
	UpdatedAt time.Time      `yaml:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `yaml:"deleted_at" json:"deleted_at" gorm:"index"`
}

// CustomAgentConfig represents the configuration of a custom agent
type CustomAgentConfig struct {
	// ===== Basic Settings =====
	// Agent mode: "quick-answer" for RAG mode, "smart-reasoning" for ReAct agent mode
	AgentMode string `yaml:"agent_mode" json:"agent_mode"`
	// System prompt for the agent (unified prompt, uses {{web_search_status}} placeholder for dynamic behavior)
	SystemPrompt string `yaml:"system_prompt" json:"system_prompt"`
	// Context template for normal mode (how to format retrieved chunks)
	ContextTemplate string `yaml:"context_template" json:"context_template"`

	// ===== Model Settings =====
	// Model ID to use for conversations
	ModelID string `yaml:"model_id" json:"model_id"`
	// ReRank model ID for retrieval
	RerankModelID string `yaml:"rerank_model_id" json:"rerank_model_id"`
	// Temperature for LLM (0-1)
	Temperature float64 `yaml:"temperature" json:"temperature"`
	// Maximum completion tokens (only for normal mode)
	MaxCompletionTokens int `yaml:"max_completion_tokens" json:"max_completion_tokens"`

	// ===== Agent Mode Settings =====
	// Maximum iterations for ReAct loop (only for agent type)
	MaxIterations int `yaml:"max_iterations" json:"max_iterations"`
	// Allowed tools (only for agent type)
	AllowedTools []string `yaml:"allowed_tools" json:"allowed_tools"`
	// Whether reflection is enabled (only for agent type)
	ReflectionEnabled bool `yaml:"reflection_enabled" json:"reflection_enabled"`
	// MCP service selection mode: "all" = all enabled MCP services, "selected" = specific services, "none" = no MCP
	MCPSelectionMode string `yaml:"mcp_selection_mode" json:"mcp_selection_mode"`
	// Selected MCP service IDs (only used when MCPSelectionMode is "selected")
	MCPServices []string `yaml:"mcp_services" json:"mcp_services"`

	// ===== Knowledge Base Settings =====
	// Knowledge base selection mode: "all" = all KBs, "selected" = specific KBs, "none" = no KB
	KBSelectionMode string `yaml:"kb_selection_mode" json:"kb_selection_mode"`
	// Associated knowledge base IDs (only used when KBSelectionMode is "selected")
	KnowledgeBases []string `yaml:"knowledge_bases" json:"knowledge_bases"`

	// ===== File Type Restriction Settings =====
	// Supported file types for this agent (e.g., ["csv", "xlsx", "xls"])
	// Empty means all file types are supported
	// When set, only files with matching extensions can be used with this agent
	SupportedFileTypes []string `yaml:"supported_file_types" json:"supported_file_types"`

	// ===== FAQ Strategy Settings =====
	// Whether FAQ priority strategy is enabled (FAQ answers prioritized over document chunks)
	FAQPriorityEnabled bool `yaml:"faq_priority_enabled" json:"faq_priority_enabled"`
	// FAQ direct answer threshold - if similarity > this value, use FAQ answer directly
	FAQDirectAnswerThreshold float64 `yaml:"faq_direct_answer_threshold" json:"faq_direct_answer_threshold"`
	// FAQ score boost multiplier - FAQ results score multiplied by this factor
	FAQScoreBoost float64 `yaml:"faq_score_boost" json:"faq_score_boost"`

	// ===== Web Search Settings =====
	// Whether web search is enabled
	WebSearchEnabled bool `yaml:"web_search_enabled" json:"web_search_enabled"`
	// Maximum web search results
	WebSearchMaxResults int `yaml:"web_search_max_results" json:"web_search_max_results"`

	// ===== Multi-turn Conversation Settings =====
	// Whether multi-turn conversation is enabled
	MultiTurnEnabled bool `yaml:"multi_turn_enabled" json:"multi_turn_enabled"`
	// Number of history turns to keep in context
	HistoryTurns int `yaml:"history_turns" json:"history_turns"`

	// ===== Retrieval Strategy Settings (for both modes) =====
	// Embedding/Vector retrieval top K
	EmbeddingTopK int `yaml:"embedding_top_k" json:"embedding_top_k"`
	// Keyword retrieval threshold
	KeywordThreshold float64 `yaml:"keyword_threshold" json:"keyword_threshold"`
	// Vector retrieval threshold
	VectorThreshold float64 `yaml:"vector_threshold" json:"vector_threshold"`
	// Rerank top K
	RerankTopK int `yaml:"rerank_top_k" json:"rerank_top_k"`
	// Rerank threshold
	RerankThreshold float64 `yaml:"rerank_threshold" json:"rerank_threshold"`

	// ===== Advanced Settings (mainly for normal mode) =====
	// Whether to enable query expansion
	EnableQueryExpansion bool `yaml:"enable_query_expansion" json:"enable_query_expansion"`
	// Whether to enable query rewrite for multi-turn conversations
	EnableRewrite bool `yaml:"enable_rewrite" json:"enable_rewrite"`
	// Rewrite prompt system message
	RewritePromptSystem string `yaml:"rewrite_prompt_system" json:"rewrite_prompt_system"`
	// Rewrite prompt user message template
	RewritePromptUser string `yaml:"rewrite_prompt_user" json:"rewrite_prompt_user"`
	// Fallback strategy: "fixed" for fixed response, "model" for model generation
	FallbackStrategy string `yaml:"fallback_strategy" json:"fallback_strategy"`
	// Fixed fallback response (when FallbackStrategy is "fixed")
	FallbackResponse string `yaml:"fallback_response" json:"fallback_response"`
	// Fallback prompt (when FallbackStrategy is "model")
	FallbackPrompt string `yaml:"fallback_prompt" json:"fallback_prompt"`
}

// Value implements driver.Valuer interface for CustomAgentConfig
func (c CustomAgentConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner interface for CustomAgentConfig
func (c *CustomAgentConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, c)
}

// TableName returns the table name for CustomAgent
func (CustomAgent) TableName() string {
	return "custom_agents"
}

// EnsureDefaults sets default values for the agent
func (a *CustomAgent) EnsureDefaults() {
	if a == nil {
		return
	}
	if a.Config.Temperature == 0 {
		a.Config.Temperature = 0.7
	}
	if a.Config.MaxIterations == 0 {
		a.Config.MaxIterations = 10
	}
	if a.Config.WebSearchMaxResults == 0 {
		a.Config.WebSearchMaxResults = 5
	}
	if a.Config.HistoryTurns == 0 {
		a.Config.HistoryTurns = 5
	}
	// Retrieval strategy defaults
	if a.Config.EmbeddingTopK == 0 {
		a.Config.EmbeddingTopK = 10
	}
	if a.Config.KeywordThreshold == 0 {
		a.Config.KeywordThreshold = 0.3
	}
	if a.Config.VectorThreshold == 0 {
		a.Config.VectorThreshold = 0.5
	}
	if a.Config.RerankTopK == 0 {
		a.Config.RerankTopK = 5
	}
	if a.Config.RerankThreshold == 0 {
		a.Config.RerankThreshold = 0.5
	}
	// Advanced settings defaults
	if a.Config.FallbackStrategy == "" {
		a.Config.FallbackStrategy = "model"
	}
	if a.Config.MaxCompletionTokens == 0 {
		a.Config.MaxCompletionTokens = 2048
	}
	// Agent mode should always enable multi-turn conversation
	if a.Config.AgentMode == AgentModeSmartReasoning {
		a.Config.MultiTurnEnabled = true
	}
}

// IsAgentMode returns true if this agent uses ReAct agent mode
func (a *CustomAgent) IsAgentMode() bool {
	return a.Config.AgentMode == AgentModeSmartReasoning
}

// GetBuiltinQuickAnswerAgent returns the built-in quick answer (RAG) mode agent
func GetBuiltinQuickAnswerAgent(tenantID uint64) *CustomAgent {
	return &CustomAgent{
		ID:          BuiltinQuickAnswerID,
		Name:        "快速问答",
		Description: "基于知识库的 RAG 问答，快速准确地回答问题",
		IsBuiltin:   true,
		TenantID:    tenantID,
		Config: CustomAgentConfig{
			AgentMode:    AgentModeQuickAnswer,
			SystemPrompt: "",
			ContextTemplate: `请根据以下参考资料回答用户问题。

参考资料：
{{contexts}}

用户问题：{{query}}`,
			Temperature:         0.7,
			MaxCompletionTokens: 2048,
			WebSearchEnabled:    true,
			WebSearchMaxResults: 5,
			MultiTurnEnabled:    true,
			HistoryTurns:        5,
			KBSelectionMode:     "all",
			// FAQ strategy
			FAQPriorityEnabled:       true,
			FAQDirectAnswerThreshold: 0.9,
			FAQScoreBoost:            1.2,
			// Retrieval strategy
			EmbeddingTopK:    10,
			KeywordThreshold: 0.3,
			VectorThreshold:  0.5,
			RerankTopK:       10,
			RerankThreshold:  0.3,
			// Advanced settings
			EnableQueryExpansion: true,
			EnableRewrite:        true,
			FallbackStrategy:     "model",
		},
	}
}

// GetBuiltinSmartReasoningAgent returns the built-in smart reasoning (ReAct) mode agent
func GetBuiltinSmartReasoningAgent(tenantID uint64) *CustomAgent {
	return &CustomAgent{
		ID:          BuiltinSmartReasoningID,
		Name:        "智能推理",
		Description: "ReAct 推理框架，支持多步思考和工具调用",
		IsBuiltin:   true,
		TenantID:    tenantID,
		Config: CustomAgentConfig{
			AgentMode:           AgentModeSmartReasoning,
			SystemPrompt:        "",
			Temperature:         0.7,
			MaxCompletionTokens: 2048,
			MaxIterations:       50,
			KBSelectionMode:     "all",
			AllowedTools:        []string{"thinking", "todo_write", "knowledge_search", "grep_chunks", "list_knowledge_chunks", "query_knowledge_graph", "get_document_info"},
			WebSearchEnabled:    true,
			WebSearchMaxResults: 5,
			ReflectionEnabled:   false,
			MultiTurnEnabled:    true,
			HistoryTurns:        5,
			// FAQ strategy
			FAQPriorityEnabled:       true,
			FAQDirectAnswerThreshold: 0.9,
			FAQScoreBoost:            1.2,
			// Retrieval strategy
			EmbeddingTopK:    10,
			KeywordThreshold: 0.3,
			VectorThreshold:  0.5,
			RerankTopK:       10,
			RerankThreshold:  0.3,
		},
	}
}

// GetBuiltinDataAnalystAgent returns the built-in data analyst agent
// This agent specializes in analyzing CSV/Excel data using SQL queries via DuckDB
func GetBuiltinDataAnalystAgent(tenantID uint64) *CustomAgent {
	return &CustomAgent{
		ID:          BuiltinDataAnalystID,
		Name:        "数据分析师",
		Description: "专业数据分析智能体，支持 CSV/Excel 文件的 SQL 查询与统计分析",
		Avatar:      "📊",
		IsBuiltin:   true,
		TenantID:    tenantID,
		Config: CustomAgentConfig{
			AgentMode:           AgentModeSmartReasoning,
			SystemPrompt: `### Role
You are WeKnora Data Analyst, an intelligent data analysis assistant powered by DuckDB. You specialize in analyzing structured data from CSV and Excel files using SQL queries.

### Mission
Help users explore, analyze, and derive insights from their tabular data through intelligent SQL query generation and execution.

### Critical Constraints
1. **Schema First:** ALWAYS call data_schema before writing any SQL query to understand the table structure.
2. **Read-Only:** Only SELECT queries allowed. INSERT, UPDATE, DELETE, CREATE, DROP are forbidden.
3. **Iterative Refinement:** If a query fails, analyze the error and refine your approach.

### Workflow
1. **Understand:** Call data_schema to get table name, columns, types, and row count.
2. **Plan:** For complex questions, use todo_write to break into sub-queries.
3. **Query:** Call data_analysis with the knowledge_id and SQL query.
4. **Analyze:** Interpret results and provide insights.

### SQL Best Practices for DuckDB
- Use double quotes for identifiers: SELECT "Column Name" FROM "table_name"
- Aggregate functions: COUNT(*), SUM(), AVG(), MIN(), MAX(), MEDIAN(), STDDEV()
- String matching: LIKE, ILIKE (case-insensitive), REGEXP
- Use LIMIT to prevent overwhelming output (default to 100 rows max)

### Tool Guidelines
- **data_schema:** ALWAYS use first. Required before any query.
- **data_analysis:** Execute SQL queries. Only SELECT queries allowed.
- **thinking:** Plan complex analyses, debug query issues.
- **todo_write:** Track multi-step analysis tasks.

### Output Standards
- Present results in well-formatted tables or summaries
- Provide actionable insights, not just raw numbers
- Relate findings back to the user's original question

Current Time: {{current_time}}
`,
			Temperature:         0.3, // Lower temperature for precise SQL generation
			MaxCompletionTokens: 4096,
			MaxIterations:       30,
			KBSelectionMode:     "all",
			// Only support CSV and Excel files for data analysis
			// Use standard values (xlsx), backend will auto-include xls via alias
			SupportedFileTypes: []string{"csv", "xlsx"},
			// Core tools for data analysis
			AllowedTools: []string{
				"thinking",
				"todo_write",
				"data_schema",   // Get table schema information
				"data_analysis", // Execute SQL queries on data
			},
			WebSearchEnabled:    false, // Data analysis doesn't need web search
			WebSearchMaxResults: 0,
			ReflectionEnabled:   true, // Enable reflection for query optimization
			MultiTurnEnabled:    true,
			HistoryTurns:        10, // More history for iterative analysis
			// Retrieval strategy (minimal, as we focus on data tools)
			EmbeddingTopK:    5,
			KeywordThreshold: 0.3,
			VectorThreshold:  0.5,
			RerankTopK:       5,
			RerankThreshold:  0.3,
		},
	}
}

// Deprecated: Use GetBuiltinQuickAnswerAgent instead
func GetBuiltinNormalAgent(tenantID uint64) *CustomAgent {
	return GetBuiltinQuickAnswerAgent(tenantID)
}

// Deprecated: Use GetBuiltinSmartReasoningAgent instead
func GetBuiltinAgentAgent(tenantID uint64) *CustomAgent {
	return GetBuiltinSmartReasoningAgent(tenantID)
}

// GetBuiltinKnowledgeRefinerAgent returns the built-in knowledge refiner agent
// This agent specializes in extracting and refining knowledge from knowledge bases
func GetBuiltinKnowledgeRefinerAgent(tenantID uint64) *CustomAgent {
	return &CustomAgent{
		ID:          BuiltinKnowledgeRefinerID,
		Name:        "知识提炼师",
		Description: "从知识库中提炼关键信息，并可将提炼结果添加到指定知识库",
		Avatar:      "💎",
		IsBuiltin:   true,
		TenantID:    tenantID,
		Config: CustomAgentConfig{
			AgentMode:           AgentModeSmartReasoning,
			SystemPrompt: `### Role
你是 WeKnora 知识提炼师，一个专业的知识提炼和整理助手。你擅长从大量文档中提取核心信息、总结要点，并将其整理成结构化的知识内容。

### Mission
帮助用户从知识库中提炼关键信息，生成高质量的知识摘要，并支持将提炼的内容添加到指定的知识库中。

### Core Capabilities
1. **知识检索与分析**：使用 knowledge_search 和 grep_chunks 工具深入检索和分析知识库内容
2. **信息提炼**：从检索结果中提取关键信息、核心观点和重要细节
3. **结构化整理**：将提炼的信息组织成清晰、结构化的 Markdown 格式
4. **知识添加**：使用 add_knowledge_to_kb 工具将提炼的内容添加到指定知识库

### Workflow
1. **理解需求**：明确用户想要提炼什么类型的信息
2. **检索知识**：使用 knowledge_search 检索相关文档和内容
3. **深度分析**：使用 grep_chunks 和 list_knowledge_chunks 获取详细信息
4. **提炼整理**：
   - 提取核心观点和关键信息
   - 去除冗余和重复内容
   - 组织成清晰的结构
   - 使用 Markdown 格式化
5. **确认添加**：询问用户是否需要将提炼的内容添加到知识库
6. **执行添加**：如果用户确认，使用 add_knowledge_to_kb 工具添加到指定知识库

### Output Standards
- 使用清晰的 Markdown 格式
- 包含标题、要点、详细说明
- 保持信息的准确性和完整性
- 适当使用列表、表格等结构化元素
- 标注信息来源（如果需要）

### Tool Guidelines
- **knowledge_search**：检索相关知识内容
- **grep_chunks**：搜索特定关键词或模式
- **list_knowledge_chunks**：获取文档的所有分块
- **get_document_info**：获取文档元信息
- **thinking**：规划提炼策略，分析信息结构
- **todo_write**：跟踪多步骤的提炼任务
- **add_knowledge_to_kb**：将提炼的内容添加到知识库

### Important Notes
- 在添加知识到知识库前，务必先向用户展示提炼的内容并征得确认
- 添加时需要用户提供目标知识库 ID
- 提炼的内容应该是高质量、结构化的 Markdown 格式
- 保持客观中立，准确传达原始信息的含义

当前时间：{{current_time}}
`,
			Temperature:         0.5, // 适中的温度，保持创造性和准确性的平衡
			MaxCompletionTokens: 4096,
			MaxIterations:       30,
			KBSelectionMode:     "all",
			// 核心工具：知识检索、分析和添加
			AllowedTools: []string{
				"thinking",
				"todo_write",
				"knowledge_search",      // 检索知识
				"grep_chunks",           // 搜索特定内容
				"list_knowledge_chunks", // 列出文档分块
				"get_document_info",     // 获取文档信息
				"add_knowledge_to_kb",   // 添加知识到知识库（新工具）
			},
			WebSearchEnabled:    false, // 专注于内部知识库，不需要网络搜索
			WebSearchMaxResults: 0,
			ReflectionEnabled:   true, // 启用反思以优化提炼质量
			MultiTurnEnabled:    true,
			HistoryTurns:        10, // 更多历史记录以支持迭代提炼
			// FAQ 策略
			FAQPriorityEnabled:       true,
			FAQDirectAnswerThreshold: 0.9,
			FAQScoreBoost:            1.2,
			// 检索策略
			EmbeddingTopK:    15, // 更多检索结果以获得全面信息
			KeywordThreshold: 0.3,
			VectorThreshold:  0.5,
			RerankTopK:       10,
			RerankThreshold:  0.3,
		},
	}
}

// GetBuiltinKnowledgeGraphExpertAgent returns the built-in knowledge graph expert agent
// This agent specializes in exploring entity relationships and knowledge networks
func GetBuiltinKnowledgeGraphExpertAgent(tenantID uint64) *CustomAgent {
	return &CustomAgent{
		ID:          BuiltinKnowledgeGraphExpertID,
		Name:        "知识图谱专家",
		Description: "探索实体关系和知识网络，深度分析知识图谱中的关联信息",
		Avatar:      "🕸️",
		IsBuiltin:   true,
		TenantID:    tenantID,
		Config: CustomAgentConfig{
			AgentMode:           AgentModeSmartReasoning,
			SystemPrompt: `### Role
你是 WeKnora 知识图谱专家，一个专业的知识网络分析助手。你擅长探索实体之间的关系、分析知识图谱结构，帮助用户理解复杂的知识关联网络。

### Mission
帮助用户探索和理解知识库中的实体关系、概念关联和知识网络结构，提供深度的关系分析和网络洞察。

### Core Capabilities
1. **实体关系探索**：使用 query_knowledge_graph 工具查询实体之间的关系
2. **知识网络分析**：分析实体的关联网络和语义连接
3. **关系可视化解释**：清晰解释实体之间的关系类型和连接路径
4. **深度关联挖掘**：发现隐藏的知识关联和间接关系

### When to Use Knowledge Graph
✅ **适合使用图谱查询的场景**：
- 理解实体之间的关系（如"Docker 和 Kubernetes 的关系"）
- 探索知识网络和概念关联
- 查找特定实体的相关信息
- 理解技术架构和系统关系
- 分析概念依赖和影响范围

❌ **不适合的场景**：
- 一般文本搜索 → 使用 knowledge_search
- 需要精确文档内容 → 使用 knowledge_search
- 知识库未配置图谱提取

### Workflow
1. **理解查询意图**：识别用户想要探索的实体或关系
2. **图谱查询**：使用 query_knowledge_graph 查询相关实体和关系
3. **结果分析**：
   - 分析实体类型和属性
   - 识别关系类型和方向
   - 评估关系强度和相关性
4. **关系解释**：
   - 清晰解释实体之间的连接
   - 说明关系的语义含义
   - 提供关系网络的整体视图
5. **深度探索**：
   - 使用 list_knowledge_chunks 获取详细内容
   - 使用 knowledge_search 补充上下文信息
   - 发现间接关系和隐藏连接

### Output Standards
- 清晰展示实体关系网络
- 使用图形化的文字描述（如树状结构、网络图）
- 解释关系的语义含义和重要性
- 提供关系强度和相关度评估
- 标注信息来源和图谱配置状态

### Tool Guidelines
- **query_knowledge_graph**：核心工具，查询实体和关系
- **knowledge_search**：补充文本搜索，获取上下文
- **list_knowledge_chunks**：获取详细文档内容
- **get_document_info**：了解文档元信息
- **thinking**：规划查询策略，分析关系网络
- **todo_write**：跟踪多步骤的图谱探索任务

### Graph Configuration Awareness
- 检查知识库是否配置了图谱提取
- 了解配置的实体类型（Nodes）和关系类型（Relations）
- 根据图谱配置调整查询策略
- 如果未配置图谱，引导用户配置或使用其他工具

### Important Notes
- 优先使用 query_knowledge_graph 工具进行图谱查询
- 关注实体之间的关系类型和语义连接
- 提供清晰的关系网络可视化描述
- 解释图谱配置对查询结果的影响
- 结合文本搜索提供全面的知识理解

当前时间：{{current_time}}
`,
			Temperature:         0.5, // 适中的温度，保持分析的准确性
			MaxCompletionTokens: 4096,
			MaxIterations:       30,
			KBSelectionMode:     "all",
			// 核心工具：图谱查询和知识检索
			AllowedTools: []string{
				"thinking",
				"todo_write",
				"query_knowledge_graph", // 核心工具：查询知识图谱
				"knowledge_search",      // 补充工具：文本搜索
				"list_knowledge_chunks", // 获取详细内容
				"get_document_info",     // 获取文档信息
			},
			WebSearchEnabled:    false, // 专注于内部知识图谱，不需要网络搜索
			WebSearchMaxResults: 0,
			ReflectionEnabled:   true, // 启用反思以优化分析质量
			MultiTurnEnabled:    true,
			HistoryTurns:        10, // 更多历史记录以支持迭代探索
			// FAQ 策略
			FAQPriorityEnabled:       false, // 图谱查询不需要 FAQ 优先
			FAQDirectAnswerThreshold: 0.9,
			FAQScoreBoost:            1.0,
			// 检索策略
			EmbeddingTopK:    10,
			KeywordThreshold: 0.3,
			VectorThreshold:  0.5,
			RerankTopK:       10,
			RerankThreshold:  0.3,
		},
	}
}

// BuiltinAgentRegistry provides a registry of all built-in agents for easy extension
var BuiltinAgentRegistry = map[string]func(uint64) *CustomAgent{
	BuiltinQuickAnswerID:          GetBuiltinQuickAnswerAgent,
	BuiltinSmartReasoningID:       GetBuiltinSmartReasoningAgent,
	BuiltinDataAnalystID:          GetBuiltinDataAnalystAgent,
	BuiltinKnowledgeRefinerID:     GetBuiltinKnowledgeRefinerAgent,
	BuiltinKnowledgeGraphExpertID: GetBuiltinKnowledgeGraphExpertAgent,
}

// builtinAgentIDsOrdered defines the fixed display order of built-in agents
var builtinAgentIDsOrdered = []string{
	BuiltinQuickAnswerID,
	BuiltinSmartReasoningID,
	BuiltinDataAnalystID,
	BuiltinKnowledgeGraphExpertID,
	BuiltinKnowledgeRefinerID,
}

// GetBuiltinAgentIDs returns all built-in agent IDs in fixed order
func GetBuiltinAgentIDs() []string {
	return builtinAgentIDsOrdered
}

// IsBuiltinAgentID checks if the given ID is a built-in agent ID
func IsBuiltinAgentID(id string) bool {
	_, exists := BuiltinAgentRegistry[id]
	return exists
}

// GetBuiltinAgent returns a built-in agent by ID, or nil if not found
func GetBuiltinAgent(id string, tenantID uint64) *CustomAgent {
	if factory, exists := BuiltinAgentRegistry[id]; exists {
		return factory(tenantID)
	}
	return nil
}
