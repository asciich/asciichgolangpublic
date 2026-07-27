package ollamaagent

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunOllamaAgent(ctx context.Context, configPath string) error {
	if configPath == "" {
		return tracederrors.TracedErrorEmptyString("configPath")
	}

	logging.LogInfoByCtxf(ctx, "Run ollama agent with config '%s' started.", configPath)

	cfg, err := LoadConfig(ctx, configPath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to load config: %w", err)
	}

	return RunOllamaAgentWithConfig(ctx, cfg)
}

func RunOllamaAgentWithConfig(ctx context.Context, cfg *Config) error {
	if cfg == nil {
		return tracederrors.TracedErrorNil("cfg")
	}

	logging.LogInfoByCtxf(ctx, "Run ollama MCP agent started.")

	a, err := NewAgent(ctx, cfg)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create agent: %w", err)
	}
	defer a.Close()

	fmt.Println("🤖 Agent ready. Type your question (or 'exit' to quit):")
	fmt.Printf("   Connected to %d MCP server(s) with %d tools available.\n\n", a.ServerCount(), a.ToolCount())

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			break
		}

		response, err := a.Chat(ctx, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		fmt.Printf("\n%s\n\n", response)
	}

	logging.LogInfoByCtxf(ctx, "Run ollama MCP agent finished.")

	return nil
}
