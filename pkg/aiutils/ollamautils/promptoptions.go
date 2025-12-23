package ollamautils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/encodingutils/base64utils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type PromptOptions struct {
	Hostname string

	Port int

	// Name of the model to use
	ModelName string

	// File paths to images/ pictures to embed into the request.
	ImagePaths []string
}

func (p *PromptOptions) GetDeepCopy() (*PromptOptions) {
	copy := new(PromptOptions)

	*copy = *p

	if p.ImagePaths != nil {
		copy.ImagePaths = slicesutils.GetDeepCopyOfStringsSlice(p.ImagePaths)
	}

	return copy
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

func (p *PromptOptions) GetImagesAsBase64Slice(ctx context.Context) ([]string, error) {
	if len(p.ImagePaths) == 0 {
		return nil, tracederrors.TracedError("No ImagePaths set to encode")
	}

	logging.LogInfoByCtxf(ctx, "Encode prompt image data into base64 slice started.")

	images := []string{}
	for _, ipath := range p.ImagePaths {
		data, err := nativefiles.ReadAsBytes(ctx, ipath)
		if err != nil {
			return nil, err
		}

		encoded, err := base64utils.EncodeBytesAsString(data)
		if err != nil {
			return nil, err
		}

		images = append(images, encoded)

		sha256 := checksumutils.GetSha256SumFromBytes(data)
		logging.LogInfoByCtxf(ctx, "Appended image '%s' with sha256sum '%s' to prompt request.", ipath, sha256)
	}

	logging.LogInfoByCtxf(ctx, "Encode prompt image data into base64 slice finished.")

	return images, nil
}
