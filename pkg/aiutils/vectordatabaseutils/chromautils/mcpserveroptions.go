package chromautils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type McpServerOptions struct {
	Port                 int
	OllamaUrl            string
	ChromaUrl            string
	ChromaCollectionName string
}

func (o *McpServerOptions) GetPort() (int, error) {
	if o.Port <= 0 {
		return 0, tracederrors.TracedErrorf("invalid Port='%d'", o.Port)
	}
	return o.Port, nil
}

func (o *McpServerOptions) GetOllamaUrl() (string, error) {
	if o.OllamaUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("OllamaUrl")
	}
	return o.OllamaUrl, nil
}

func (o *McpServerOptions) GetChromaUrl() (string, error) {
	if o.ChromaUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("ChromaUrl")
	}
	return o.ChromaUrl, nil
}

func (o *McpServerOptions) GetChromaCollectionName() (string, error) {
	if o.ChromaCollectionName == "" {
		return "", tracederrors.TracedErrorEmptyString("ChromaCollectionName")
	}
	return o.ChromaCollectionName, nil
}
