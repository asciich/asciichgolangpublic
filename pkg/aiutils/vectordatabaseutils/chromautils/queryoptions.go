package chromautils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type QueryOptions struct {
	OllamaUrl      string
	ChromaUrl      string
	CollectionName string
	EmbeddingModel string
	LlmModel       string
	Question       string
	NResults       int
}

func (o *QueryOptions) GetOllamaUrl() (string, error) {
	if o.OllamaUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("OllamaUrl")
	}
	return o.OllamaUrl, nil
}

func (o *QueryOptions) GetChromaUrl() (string, error) {
	if o.ChromaUrl == "" {
		return "", tracederrors.TracedErrorEmptyString("ChromaUrl")
	}
	return o.ChromaUrl, nil
}

func (o *QueryOptions) GetCollectionName() (string, error) {
	if o.CollectionName == "" {
		return "", tracederrors.TracedErrorEmptyString("CollectionName")
	}
	return o.CollectionName, nil
}

func (o *QueryOptions) GetQuestion() (string, error) {
	if o.Question == "" {
		return "", tracederrors.TracedErrorEmptyString("Question")
	}
	return o.Question, nil
}

func (o *QueryOptions) GetEmbeddingModelOrDefault() string {
	if o.EmbeddingModel == "" {
		return "nomic-embed-text"
	}
	return o.EmbeddingModel
}

func (o *QueryOptions) GetLlmModelOrDefault() string {
	if o.LlmModel == "" {
		return "llama3"
	}
	return o.LlmModel
}

func (o *QueryOptions) GetNResultsOrDefault() int {
	if o.NResults <= 0 {
		return 4
	}
	return o.NResults
}
