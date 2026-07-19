package chromautils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const DEFAULT_EMBEDDING_MODEL_NAME = "nomic-embed-text"
const DEFAULT_CHUNK_SIZE = 256
const DEFAULT_CHUNK_OVERLAP = 30

type IndexOptions struct {
	// Full URL to the running Ollama instance.
	//  E.g. http://192.168.1.x:11434
	OllamaUrl string

	// Full URL to the running Chroma instance
	ChromaUrl string

	ChromaCollectionName string

	// The root directory of the documents to index.
	DocumentsDirectory string

	EmbeddingModuelName string
	ChunkSize           int
	ChunkOverlap        int

	BaseNameRegex string
}

func (i *IndexOptions) GetEmbeddingModelNameOrDefault() string {
	if i.EmbeddingModuelName == "" {
		return DEFAULT_EMBEDDING_MODEL_NAME
	}

	return i.EmbeddingModuelName
}

func (i *IndexOptions) GetOllamaUrl() (string, error) {
	if i.OllamaUrl == "" {
		return "", tracederrors.TracedError("OllamaUrl not set")
	}

	return i.OllamaUrl, nil
}

func (i *IndexOptions) GetDocumentsDirectory() (string, error) {
	if i.DocumentsDirectory == "" {
		return "", tracederrors.TracedError("DocumentsDirectory not set")
	}

	return i.DocumentsDirectory, nil
}

func (i *IndexOptions) GetChunkSizeOrDefault() (int, error) {
	if i.ChunkSize == 0 {
		return DEFAULT_CHUNK_SIZE, nil
	}

	if i.ChunkSize < 0 {
		return 0, tracederrors.TracedErrorf("Invalid ChunkSize '%d'", i.ChunkSize)
	}

	return i.ChunkSize, nil
}

func (i *IndexOptions) GetChunkOverlapOrDefault() (int, error) {
	if i.ChunkOverlap == 0 {
		return DEFAULT_CHUNK_OVERLAP, nil
	}

	if i.ChunkOverlap < 0 {
		return 0, tracederrors.TracedErrorf("Invalid ChunkOverlap '%d'", i.ChunkOverlap)
	}

	return i.ChunkOverlap, nil
}

func (i *IndexOptions) GetChunkSizeAndOverlapOrDefault() (int, int, error) {
	chunkSize, err := i.GetChunkSizeOrDefault()
	if err != nil {
		return 0, 0, err
	}

	chunkOverlap, err := i.GetChunkOverlapOrDefault()
	if err != nil {
		return 0, 0, err
	}

	return chunkSize, chunkOverlap, nil
}

func (i *IndexOptions) GetChromaUrl() (string, error) {
	if i.ChromaUrl == "" {
		return "", tracederrors.TracedError("ChromaUrl not set")
	}

	return i.ChromaUrl, nil
}

func (i *IndexOptions) GetChromaCollectionName() (string, error) {
	if i.ChromaCollectionName == "" {
		return "", tracederrors.TracedError("ChromaCollectionName not set")
	}

	return i.ChromaCollectionName, nil
}

func (o *IndexOptions) GetBaseNameRegex() (string, error) {
	if o.BaseNameRegex == "" {
		return "", tracederrors.TracedError("BaseNameRegex not set")
	}

	return o.BaseNameRegex, nil
}
