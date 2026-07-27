package ollamaagent

import (
	"context"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ToolRegistration struct {
	Tool       mcp.Tool
	Client     client.MCPClient
	ServerName string
}

type McpManager struct {
	clients map[string]client.MCPClient
	tools   map[string]*ToolRegistration
}

func NewMcpManager() *McpManager {
	return &McpManager{
		clients: make(map[string]client.MCPClient),
		tools:   make(map[string]*ToolRegistration),
	}
}

func (m *McpManager) Connect(ctx context.Context, serverCfg McpServerConfig) error {
	if serverCfg.Name == "" {
		return tracederrors.TracedErrorEmptyString("serverCfg.Name")
	}

	if serverCfg.URL == "" {
		return tracederrors.TracedErrorEmptyString("serverCfg.URL")
	}

	logging.LogInfoByCtxf(ctx, "Connect to MCP server '%s' at '%s' started.", serverCfg.Name, serverCfg.URL)

	c, err := client.NewSSEMCPClient(serverCfg.URL)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create SSE client for '%s': %w", serverCfg.Name, err)
	}

	if err := c.Start(ctx); err != nil {
		return tracederrors.TracedErrorf("Failed to start SSE client for '%s': %w", serverCfg.Name, err)
	}

	initReq := mcp.InitializeRequest{}
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "ollamaagent",
		Version: "1.0.0",
	}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION

	_, err = c.Initialize(ctx, initReq)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to initialize '%s': %w", serverCfg.Name, err)
	}

	toolsResult, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to list tools from '%s': %w", serverCfg.Name, err)
	}

	m.clients[serverCfg.Name] = c

	for _, tool := range toolsResult.Tools {
		toolName := tool.Name
		if _, exists := m.tools[toolName]; exists {
			toolName = serverCfg.Name + "." + toolName
		}
		m.tools[toolName] = &ToolRegistration{
			Tool:       tool,
			Client:     c,
			ServerName: serverCfg.Name,
		}
	}

	logging.LogInfoByCtxf(ctx, "Connect to MCP server '%s' at '%s' finished. Discovered %d tools.", serverCfg.Name, serverCfg.URL, len(toolsResult.Tools))

	return nil
}

func (m *McpManager) GetAllTools() []ToolRegistration {
	tools := make([]ToolRegistration, 0, len(m.tools))
	for _, reg := range m.tools {
		tools = append(tools, *reg)
	}
	return tools
}

func (m *McpManager) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	if toolName == "" {
		return nil, tracederrors.TracedErrorEmptyString("toolName")
	}

	logging.LogInfoByCtxf(ctx, "Call MCP tool '%s' started.", toolName)

	reg, ok := m.tools[toolName]
	if !ok {
		return nil, tracederrors.TracedErrorf("Tool '%s' not found", toolName)
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = reg.Tool.Name
	req.Params.Arguments = arguments

	result, err := reg.Client.CallTool(ctx, req)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Tool '%s' call failed: %w", toolName, err)
	}

	logging.LogInfoByCtxf(ctx, "Call MCP tool '%s' finished.", toolName)

	return result, nil
}

func (m *McpManager) ServerCount() int {
	return len(m.clients)
}

func (m *McpManager) ToolCount() int {
	return len(m.tools)
}

func (m *McpManager) Close() {
	for _, c := range m.clients {
		c.Close()
	}
}
