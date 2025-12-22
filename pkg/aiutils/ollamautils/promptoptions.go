package ollamautils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

type PromptOptions struct {
	Hostname string

	Port int

	// Name of the model to use
	ModelName string
}

func (p *PromptOptions) GetHostnameOrDefault() (string, error) {
	if p.Hostname != "" {
		return p.Hostname, nil
	}

	return "localhost", nil
}

func (p *PromptOptions) GetPortOrDefault() (int, error) {
	if p.Port == 0 {
		return GetDefaultPort(), nil
	}

	return p.Port, nil
}

func (p *PromptOptions) GetGenerateUrl(ctx context.Context) (string, error) {
	hostname, err := p.GetHostnameOrDefault()
	if err != nil {
		return "", err
	}

	port, err := p.GetPortOrDefault()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://%s:%d/api/generate", hostname, port)

	logging.LogInfoByCtxf(ctx, "Generate URL is %s .", url)

	return url, nil
}

func (p *PromptOptions) GetModelNameOrDefault(ctx context.Context) (string, error) {
	if p.ModelName == "" {
		modelName := GetFastModelName()
		logging.LogInfoByCtxf(ctx, "Model name is not explicitly set. Use default ollama model name '%s'.", modelName)
		return modelName, nil
	}

	modelName := p.ModelName
	logging.LogInfoByCtxf(ctx, "Use explicit set ollama model name: '%s'.", modelName)
	return modelName, nil
}
