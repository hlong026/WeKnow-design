package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

var addKnowledgeToKBTool = BaseTool{
	name: ToolAddKnowledgeToKB,
	description: `Add refined knowledge content to a specified knowledge base.

This tool allows you to add manually created or refined knowledge content (in Markdown format) to a knowledge base.

## Purpose
- Add extracted and refined information to knowledge bases
- Create structured knowledge entries from analysis results
- Save summarized or processed content for future retrieval

## Workflow
1. Present the refined content to the user for review
2. Ask the user to confirm and provide the target knowledge base ID
3. Use this tool to add the content to the specified knowledge base

## Required Input
- knowledge_base_id: The ID of the target knowledge base
- title: A descriptive title for the knowledge entry
- content: The knowledge content in Markdown format

## Important Notes
- Always show the content to the user before adding
- Get user confirmation and the target knowledge base ID
- Content should be well-formatted Markdown
- Title should be descriptive and concise

## Output
Returns success status and the created knowledge entry ID.`,
	schema: json.RawMessage(`{
  "type": "object",
  "properties": {
    "knowledge_base_id": {
      "type": "string",
      "description": "REQUIRED: The ID of the target knowledge base"
    },
    "title": {
      "type": "string",
      "description": "REQUIRED: A descriptive title for the knowledge entry"
    },
    "content": {
      "type": "string",
      "description": "REQUIRED: The knowledge content in Markdown format"
    }
  },
  "required": ["knowledge_base_id", "title", "content"]
}`),
}

// AddKnowledgeToKBInput defines the input parameters for add_knowledge_to_kb tool
type AddKnowledgeToKBInput struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
}

// AddKnowledgeToKBTool adds knowledge content to a knowledge base
type AddKnowledgeToKBTool struct {
	BaseTool
	knowledgeService interfaces.KnowledgeService
}

// NewAddKnowledgeToKBTool creates a new add_knowledge_to_kb tool
func NewAddKnowledgeToKBTool(
	knowledgeService interfaces.KnowledgeService,
) *AddKnowledgeToKBTool {
	return &AddKnowledgeToKBTool{
		BaseTool:         addKnowledgeToKBTool,
		knowledgeService: knowledgeService,
	}
}

// Execute executes the add_knowledge_to_kb tool
func (t *AddKnowledgeToKBTool) Execute(ctx context.Context, args json.RawMessage) (*types.ToolResult, error) {
	logger.Infof(ctx, "[Tool][AddKnowledgeToKB] Execute started")

	// Parse args
	var input AddKnowledgeToKBInput
	if err := json.Unmarshal(args, &input); err != nil {
		logger.Errorf(ctx, "[Tool][AddKnowledgeToKB] Failed to parse args: %v", err)
		return &types.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse args: %v", err),
		}, err
	}

	// Validate input
	if strings.TrimSpace(input.KnowledgeBaseID) == "" {
		err := fmt.Errorf("knowledge_base_id is required")
		logger.Errorf(ctx, "[Tool][AddKnowledgeToKB] %v", err)
		return &types.ToolResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if strings.TrimSpace(input.Title) == "" {
		err := fmt.Errorf("title is required")
		logger.Errorf(ctx, "[Tool][AddKnowledgeToKB] %v", err)
		return &types.ToolResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if strings.TrimSpace(input.Content) == "" {
		err := fmt.Errorf("content is required")
		logger.Errorf(ctx, "[Tool][AddKnowledgeToKB] %v", err)
		return &types.ToolResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	logger.Infof(ctx, "[Tool][AddKnowledgeToKB] Adding knowledge to KB: %s, title: %s, content length: %d",
		input.KnowledgeBaseID, input.Title, len(input.Content))

	// Create knowledge entry using manual knowledge payload
	payload := &types.ManualKnowledgePayload{
		Title:   input.Title,
		Content: input.Content,
		Status:  "published", // Set status to published
	}

	knowledge, err := t.knowledgeService.CreateKnowledgeFromManual(
		ctx,
		input.KnowledgeBaseID,
		payload,
	)
	if err != nil {
		logger.Errorf(ctx, "[Tool][AddKnowledgeToKB] Failed to create knowledge: %v", err)
		return &types.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to add knowledge to knowledge base: %v", err),
		}, err
	}

	logger.Infof(ctx, "[Tool][AddKnowledgeToKB] Successfully created knowledge: %s", knowledge.ID)

	// Format output
	output := fmt.Sprintf(`✅ Successfully added knowledge to knowledge base!

**Knowledge ID:** %s
**Title:** %s
**Knowledge Base ID:** %s
**Content Length:** %d characters

The knowledge has been processed and will be available for search shortly.`,
		knowledge.ID,
		input.Title,
		input.KnowledgeBaseID,
		len(input.Content),
	)

	return &types.ToolResult{
		Success: true,
		Output:  output,
	}, nil
}
