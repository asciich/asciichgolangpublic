package chromautils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/aiutils/vectordatabaseutils/vectordatabasegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type OllamaEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type OllamaEmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func getEmbeddings(ctx context.Context, texts []string, options *IndexOptions) ([][]float32, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	embeddingModel := options.GetEmbeddingModelNameOrDefault()

	ollamaUrl, err := options.GetOllamaUrl()
	if err != nil {
		return nil, err
	}

	var totalLen int64
	for _, t := range texts {
		totalLen += int64(len([]byte(t)))
	}

	logging.LogInfoByCtxf(ctx, "Get embeddings for %d texts with a total size of %d bytes using the model '%s' on ollama '%s' started.", len(texts), totalLen, embeddingModel, ollamaUrl)

	reqBody := OllamaEmbedRequest{
		Model: embeddingModel,
		Input: texts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ollamaUrl+"/api/embed", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, tracederrors.TracedErrorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("ollama returned status %d: %s", resp.StatusCode, body)
	}

	var embedResp OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode ollama response: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Get embeddings for %d texts with a total size of %d bytes using the model '%s' on ollama '%s' finished.", len(texts), totalLen, embeddingModel, ollamaUrl)

	return embedResp.Embeddings, nil
}

func loadDocuments(ctx context.Context, dir string, options *IndexOptions) ([]string, []string, error) {
	if dir == "" {
		return nil, nil, tracederrors.TracedErrorEmptyString("dir")
	}

	if options == nil {
		return nil, nil, tracederrors.TracedError(options)
	}

	logging.LogInfoByCtxf(ctx, "Load documents in in directory '%s' started.", dir)

	var contents []string
	var filenames []string

	basenameRegex, err := options.GetBaseNameRegex()
	if err != nil {
		return nil, nil, err
	}

	regex, err := regexp.Compile(basenameRegex)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to compile regex '%s': %w", basenameRegex, err)
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return tracederrors.TracedErrorf("Failed in WalkDirs: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		baseName := filepath.Base(path)

		if regex.Match([]byte(baseName)) {
			data, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
			if err != nil {
				return err
			}
			contents = append(contents, string(data))
			filenames = append(filenames, path)
		}
		return nil
	})

	logging.LogInfoByCtxf(ctx, "Load documents in in directory '%s' finished. Loaded %d documents.", dir, len(contents))

	return contents, filenames, err
}

type Chunk struct {
	Text   string
	Source string
}

func splitIntoChunks(ctx context.Context, contents []string, filenames []string, options *IndexOptions) ([]*Chunk, error) {
	if contents == nil {
		return nil, tracederrors.TracedErrorNil("contents")
	}

	if filenames == nil {
		return nil, tracederrors.TracedErrorNil("filenames")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if len(contents) <= 0 {
		return nil, tracederrors.TracedErrorf("Contents has no elements, nothing to split into chunks")
	}

	if len(contents) != len(filenames) {
		return nil, tracederrors.TracedErrorf("Length mismatch: len(contents)=%d != len(filenames)=%d", len(contents), len(filenames))
	}

	chunkSize, chunkOverlap, err := options.GetChunkSizeAndOverlapOrDefault()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Split %d files into chunks with size=%d and overlap=%d started.", len(contents), chunkSize, chunkOverlap)

	var chunks = []*Chunk{}

	for i, content := range contents {
		parts := vectordatabasegeneric.SplitText(content, chunkSize, chunkOverlap)
		for _, p := range parts {
			if strings.TrimSpace(p) != "" {
				chunks = append(chunks, &Chunk{Text: p, Source: filenames[i]})
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Split %d files into chunks with size=%d and overlap=%d finished. Generated %d chunks.", len(contents), chunkSize, chunkOverlap, len(chunks))

	return chunks, nil
}

func IndexDocuments(ctx context.Context, options *IndexOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtxf(ctx, "Index documents into chroma started.")

	docsDir, err := options.GetDocumentsDirectory()
	if err != nil {
		return err
	}

	contents, filenames, err := loadDocuments(ctx, docsDir, options)
	if err != nil {
		return err
	}

	chunks, err := splitIntoChunks(ctx, contents, filenames, options)
	if err != nil {
		return err
	}

	// 3. Generate embeddings (batch)
	fmt.Println("Generating embeddings via Ollama...")
	batchSize := 32
	var allEmbeddings [][]float32

	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		batch := make([]string, end-i)
		for j, c := range chunks[i:end] {
			batch[j] = c.Text
		}

		embeddings, err := getEmbeddings(ctx, batch, options)
		if err != nil {
			return err
		}
		allEmbeddings = append(allEmbeddings, embeddings...)
		logging.LogInfoByCtxf(ctx, "Embedded batch %d/%d\n", i/batchSize+1, (len(chunks)+batchSize-1)/batchSize)
	}

	// 4. Store in Chroma
	logging.LogInfoByCtxf(ctx, "Storing in Chroma...")
	chromaUrl, err := options.GetChromaUrl()
	if err != nil {
		return err
	}

	client := NewClient(chromaUrl)

	// Check connectivity
	if err := client.CheckReachable(ctx); err != nil {
		return err
	}

	collectionName, err := options.GetChromaCollectionName()
	if err != nil {
		return err
	}

	// Delete collection if it exists (fresh start)
	_ = client.DeleteCollection(ctx, collectionName)

	// Create collection
	collection, err := client.CreateCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	// Add in batches
	addBatchSize := 100
	for i := 0; i < len(chunks); i += addBatchSize {
		end := i + addBatchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		ids := make([]string, end-i)
		docs := make([]string, end-i)
		metas := make([]map[string]any, end-i)
		embeds := make([][]float32, end-i)

		for j := 0; j < end-i; j++ {
			ids[j] = fmt.Sprintf("doc_%d", i+j)
			docs[j] = chunks[i+j].Text
			metas[j] = map[string]any{"source": chunks[i+j].Source}
			embeds[j] = allEmbeddings[i+j]
		}

		err = client.Add(ctx, collection.ID, ids, embeds, docs, metas)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to add batch %d to Chroma: %w", i/addBatchSize, err)
		}

		logging.LogInfoByCtxf(ctx, "Added batch %d/%d to Chroma", i/addBatchSize+1, (len(chunks)+addBatchSize-1)/addBatchSize)
	}

	logging.LogInfoByCtxf(ctx, "Index '%d' documents into chroma finished.", len(filenames))

	return nil
}
