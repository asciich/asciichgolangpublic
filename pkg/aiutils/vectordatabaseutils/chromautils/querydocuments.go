package chromautils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// --- Ollama Chat ---

type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Message OllamaMessage `json:"message"`
}

func askOllama(ctx context.Context, ollamaUrl, model, question, ragContext string) (string, error) {
	logging.LogInfoByCtxf(ctx, "Ask ollama model '%s' on '%s' started.", model, ollamaUrl)

	prompt := fmt.Sprintf(`Answer the question based only on the following context. If the answer is not in the context, say "I don't have enough information to answer this question."

---
Context:
%s
---

Question: %s`, ragContext, question)

	reqBody := OllamaChatRequest{
		Model: model,
		Messages: []OllamaMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant. Answer questions based only on the provided context. Be concise and accurate.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to marshal chat request: %w", err)
	}

	resp, err := http.Post(ollamaUrl+"/api/chat", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", tracederrors.TracedErrorf("Ollama chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", tracederrors.TracedErrorf("Ollama chat returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", tracederrors.TracedErrorf("Failed to decode chat response: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Ask ollama model '%s' on '%s' finished.", model, ollamaUrl)

	return chatResp.Message.Content, nil
}

// --- Query Result Display ---

type QueryDocumentsResult struct {
	Answer         string
	SourceChunks   []string
	SourceFiles    []string
	Distances      []float32
}


func QueryDocuments(ctx context.Context, options *QueryOptions) (*QueryDocumentsResult, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	question, err := options.GetQuestion()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Query documents started. Question: '%s'", question)

	// 1. Get embedding for the question
	ollamaUrl, err := options.GetOllamaUrl()
	if err != nil {
		return nil, err
	}

	embeddingModel := options.GetEmbeddingModelOrDefault()

	logging.LogInfoByCtxf(ctx, "Generating embedding for query using model '%s'...", embeddingModel)

	reqBody := OllamaEmbedRequest{
		Model: embeddingModel,
		Input: []string{question},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal embed request: %w", err)
	}

	resp, err := http.Post(ollamaUrl+"/api/embed", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Ollama embed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("Ollama embed returned status %d: %s", resp.StatusCode, string(body))
	}

	var embedResp OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode embed response: %w", err)
	}

	if len(embedResp.Embeddings) == 0 {
		return nil, tracederrors.TracedError("No embeddings returned for query")
	}

	queryEmbedding := embedResp.Embeddings[0]
	logging.LogInfoByCtxf(ctx, "Query embedding generated (dimension: %d).", len(queryEmbedding))

	// 2. Search in Chroma
	chromaUrl, err := options.GetChromaUrl()
	if err != nil {
		return nil, err
	}

	client := NewClient(chromaUrl)

	if err := client.CheckReachable(ctx); err != nil {
		return nil, err
	}

	collectionName, err := options.GetCollectionName()
	if err != nil {
		return nil, err
	}

	collection, err := client.GetCollectionByName(ctx, collectionName)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get collection '%s': %w", collectionName, err)
	}

	logging.LogInfoByCtxf(ctx, "Querying collection '%s' (id: %s) for top %d results...", collectionName, collection.ID, options.GetNResultsOrDefault())

	results, err := client.Query(collection.ID, [][]float32{queryEmbedding}, options.GetNResultsOrDefault())
	if err != nil {
		return nil, tracederrors.TracedErrorf("Chroma query failed: %w", err)
	}

	if len(results.Documents) == 0 || len(results.Documents[0]) == 0 {
		return nil, tracederrors.TracedError("No results found in collection")
	}

	// 3. Build context from results
	var contextParts []string
	var sourceFiles []string
	var distances []float32

	for i, doc := range results.Documents[0] {
		contextParts = append(contextParts, doc)

		if i < len(results.Distances[0]) {
			distances = append(distances, results.Distances[0][i])
		}

		if i < len(results.Metadatas[0]) {
			if source, ok := results.Metadatas[0][i]["source"].(string); ok {
				sourceFiles = append(sourceFiles, source)
			}
		}

		logging.LogInfoByCtxf(ctx, "  Result[%d]: distance=%.4f, source=%s, preview=%s...",
			i,
			results.Distances[0][i],
			sourceFiles[len(sourceFiles)-1],
			truncateString(doc, 80),
		)
	}

	ragContext := strings.Join(contextParts, "\n---\n")

	// 4. Ask LLM with context
	llmModel := options.GetLlmModelOrDefault()
	logging.LogInfoByCtxf(ctx, "Asking LLM model '%s' with %d context chunks...", llmModel, len(contextParts))

	answer, err := askOllama(ctx, ollamaUrl, llmModel, question, ragContext)
	if err != nil {
		return nil, tracederrors.TracedErrorf("LLM query failed: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Query documents finished.")

	return &QueryDocumentsResult{
		Answer:       answer,
		SourceChunks: contextParts,
		SourceFiles:  sourceFiles,
		Distances:    distances,
	}, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
