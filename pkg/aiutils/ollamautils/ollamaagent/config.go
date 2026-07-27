package ollamaagent

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type McpServerConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type OllamaConfig struct {
	URL   string `yaml:"url"`
	Model string `yaml:"model"`
}

type Config struct {
	Ollama     OllamaConfig      `yaml:"ollama"`
	McpServers []McpServerConfig `yaml:"mcp_servers"`
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Load config from '%s' started.", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read config file '%s': %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse config file '%s': %w", path, err)
	}

	logging.LogInfoByCtxf(ctx, "Load config from '%s' finished.", path)

	return &cfg, nil
}
