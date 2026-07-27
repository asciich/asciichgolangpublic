package ollamaagent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const maxToolCallIterations = 10

type Agent struct {
	mcpManager   *McpManager
	ollama       *OllamaClient
	messages     []OllamaMessage
	tools        []OllamaTool
	systemPrompt string
}

func NewAgent(ctx context.Context, cfg *Config) (*Agent, error) {
	if cfg == nil {
		return nil, tracederrors.TracedErrorNil("cfg")
	}

	logging.LogInfoByCtxf(ctx, "Create ollama agent started.")

	manager := NewMcpManager()
	for _, serverCfg := range cfg.McpServers {
		if err := manager.Connect(ctx, serverCfg); err != nil {
			manager.Close()
			return nil, tracederrors.TracedErrorf("Failed to connect to '%s': %w", serverCfg.Name, err)
		}
	}

	ollamaTools := convertMcpToolsToOllama(manager.GetAllTools())

	systemPrompt := `You are a helpful AI assistant with access to tools.
Use the available tools to search and retrieve information when needed to answer questions accurately.
Always cite the source of information when available in metadata.
If you cannot find relevant information using the tools, say so honestly.

IMPORTANT tool usage rules:
- Always use 'search_documents' FIRST to find relevant chunks and discover the correct source file paths.
- Only use 'get_full_document' AFTER you have found the exact source path from a search_documents result metadata.
- Never guess or invent source file paths.
- Never guess or invent collection names. Either omit the collection parameter to use the default, or use 'list_collections' first to discover available collections.
- When calling search_documents, do NOT set the 'collection' parameter unless you have confirmed it exists via list_collections.`

	ollama := NewOllamaClient(cfg.Ollama.URL, cfg.Ollama.Model)

	logging.LogInfoByCtxf(ctx, "Create ollama agent finished. Connected to %d server(s) with %d tools available.", manager.ServerCount(), manager.ToolCount())

	return &Agent{
		mcpManager: manager,
		ollama:     ollama,
		tools:      ollamaTools,
		messages: []OllamaMessage{
			{Role: "system", Content: systemPrompt},
		},
		systemPrompt: systemPrompt,
	}, nil
}

func (a *Agent) Chat(ctx context.Context, userInput string) (string, error) {
	if userInput == "" {
		return "", tracederrors.TracedErrorEmptyString("userInput")
	}

	logging.LogInfoByCtxf(ctx, "Chat with ollama agent started.")

	a.messages = append(a.messages, OllamaMessage{
		Role:    "user",
		Content: userInput,
	})

	for i := 0; i < maxToolCallIterations; i++ {
		resp, err := a.ollama.Chat(a.messages, a.tools)
		if err != nil {
			return "", tracederrors.TracedErrorf("Ollama chat failed: %w", err)
		}

		if len(resp.Message.ToolCalls) == 0 {
			a.messages = append(a.messages, resp.Message)

			logging.LogInfoByCtxf(ctx, "Chat with ollama agent finished.")

			return resp.Message.Content, nil
		}

		a.messages = append(a.messages, resp.Message)

		for _, toolCall := range resp.Message.ToolCalls {
			fmt.Printf("  🔧 Calling tool: %s(%v)\n", toolCall.Function.Name, toolCall.Function.Arguments)

			result, err := a.mcpManager.CallTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				a.messages = append(a.messages, OllamaMessage{
					Role:    "tool",
					Content: fmt.Sprintf("Error calling %s: %v", toolCall.Function.Name, err),
				})
				continue
			}

			toolOutput := extractToolResultText(result)
			fmt.Printf("  ✅ Tool '%s' returned %d bytes\n", toolCall.Function.Name, len(toolOutput))

			a.messages = append(a.messages, OllamaMessage{
				Role:    "tool",
				Content: toolOutput,
			})
		}
	}

	return "", tracederrors.TracedErrorf("Max tool call iterations (%d) reached", maxToolCallIterations)
}

func (a *Agent) ServerCount() int {
	return a.mcpManager.ServerCount()
}

func (a *Agent) ToolCount() int {
	return a.mcpManager.ToolCount()
}

func (a *Agent) Close() {
	a.mcpManager.Close()
}

func (a *Agent) ResetConversation() {
	a.messages = []OllamaMessage{
		{Role: "system", Content: a.systemPrompt},
	}
}

func convertMcpToolsToOllama(tools []ToolRegistration) []OllamaTool {
	ollamaTools := make([]OllamaTool, 0, len(tools))

	for _, reg := range tools {
		params := make(map[string]interface{})
		if reg.Tool.InputSchema.Properties != nil {
			params["type"] = "object"
			params["properties"] = reg.Tool.InputSchema.Properties
			if len(reg.Tool.InputSchema.Required) > 0 {
				params["required"] = reg.Tool.InputSchema.Required
			}
		}

		ollamaTools = append(ollamaTools, OllamaTool{
			Type: "function",
			Function: OllamaFunction{
				Name:        reg.Tool.Name,
				Description: reg.Tool.Description,
				Parameters:  params,
			},
		})
	}

	return ollamaTools
}

func extractToolResultText(result *mcp.CallToolResult) string {
	if result == nil {
		return ""
	}

	var parts []string
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			parts = append(parts, textContent.Text)
		} else {
			data, err := json.Marshal(content)
			if err == nil {
				parts = append(parts, string(data))
			}
		}
	}

	if len(parts) == 1 {
		return parts[0]
	}

	combined, _ := json.Marshal(parts)
	return string(combined)
}
