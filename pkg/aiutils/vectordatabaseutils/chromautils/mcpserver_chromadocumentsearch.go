package chromautils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type McpQueryResult struct {
	IDs       [][]string                 `json:"ids"`
	Documents [][]string                 `json:"documents"`
	Metadatas [][]map[string]interface{} `json:"metadatas"`
	Distances [][]float32                `json:"distances"`
}

type McpGetResult struct {
	IDs       []string                 `json:"ids"`
	Documents []string                 `json:"documents"`
	Metadatas []map[string]interface{} `json:"metadatas"`
}

func mcpGetEmbeddings(ollamaURL string, texts []string) ([][]float32, error) {
	reqBody, err := json.Marshal(map[string]interface{}{
		"model": DEFAULT_EMBEDDING_MODEL_NAME,
		"input": texts,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal embed request: %w", err)
	}
	resp, err := http.Post(ollamaURL+"/api/embed", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to post to ollama: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("ollama embed failed (%d): %s", resp.StatusCode, body)
	}
	var result struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode ollama response: %w", err)
	}
	return result.Embeddings, nil
}

func mcpQueryChroma(chromaURL string, collectionID string, embeddings [][]float32, nResults int) (*McpQueryResult, error) {
	body, err := json.Marshal(map[string]interface{}{
		"query_embeddings": embeddings,
		"n_results":        nResults,
		"include":          []string{"documents", "metadatas", "distances"},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal query request: %w", err)
	}
	url := fmt.Sprintf("%s%s/collections/%s/query", chromaURL, v2BasePath, collectionID)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to post query to chroma: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("query failed (%d): %s", resp.StatusCode, respBody)
	}
	var result McpQueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode query result: %w", err)
	}
	return &result, nil
}

func mcpGetBySource(chromaURL string, collectionID string, source string) (*McpGetResult, error) {
	body, err := json.Marshal(map[string]interface{}{
		"where": map[string]interface{}{
			"source": source,
		},
		"include": []string{"documents", "metadatas"},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal get request: %w", err)
	}
	url := fmt.Sprintf("%s%s/collections/%s/get", chromaURL, v2BasePath, collectionID)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to post get to chroma: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("get by source failed (%d): %s", resp.StatusCode, respBody)
	}
	var result McpGetResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode get result: %w", err)
	}
	return &result, nil
}

func RunDocumentSearchMcpServer(ctx context.Context, options *McpServerOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	port, err := options.GetPort()
	if err != nil {
		return err
	}

	chromaURL, err := options.GetChromaUrl()
	if err != nil {
		return err
	}

	ollamaURL, err := options.GetOllamaUrl()
	if err != nil {
		return err
	}

	collectionName, err := options.GetChromaCollectionName()
	if err != nil {
		return err
	}

	s := server.NewMCPServer(
		"Chroma Document Search",
		"1.0.0",
	)

	chromaClient := NewClient(chromaURL)

	// Capture the outer context for logging inside handlers
	serverCtx := ctx

	// Tool: search_documents
	searchTool := mcp.NewTool("search_documents",
		mcp.WithDescription("Search indexed documents in ChromaDB using semantic similarity"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Natural language search query")),
		mcp.WithNumber("n_results", mcp.Description("Number of results (default 5)")),
		mcp.WithString("collection", mcp.Description("Collection name (optional)")),
	)

	s.AddTool(searchTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		logging.LogInfoByCtxf(serverCtx, "Tool 'search_documents' call started.")

		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, tracederrors.TracedErrorf("Invalid arguments type: %T", req.Params.Arguments)
		}

		query, ok := args["query"].(string)
		if !ok || query == "" {
			return nil, tracederrors.TracedErrorf("Missing or invalid 'query' argument")
		}

		nResults := 5
		if n, ok := args["n_results"].(float64); ok {
			nResults = int(n)
		}

		colName := collectionName
		if c, ok := args["collection"].(string); ok && c != "" {
			colName = c
		}

		// 1. Get collection
		col, err := chromaClient.GetCollectionByName(serverCtx, colName)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get collection: %w", err)
		}

		// 2. Embed query
		embeddings, err := mcpGetEmbeddings(ollamaURL, []string{query})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get embeddings: %w", err)
		}

		// 3. Query
		results, err := mcpQueryChroma(chromaURL, col.ID, embeddings, nResults)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Query failed: %w", err)
		}

		// 4. Format output
		type ResultEntry struct {
			ID       string                 `json:"id"`
			Document string                 `json:"document"`
			Metadata map[string]interface{} `json:"metadata"`
			Distance float32                `json:"distance"`
		}

		var entries []ResultEntry
		if len(results.IDs) > 0 {
			for i, id := range results.IDs[0] {
				entry := ResultEntry{ID: id}
				if len(results.Documents) > 0 && len(results.Documents[0]) > i {
					entry.Document = results.Documents[0][i]
				}
				if len(results.Metadatas) > 0 && len(results.Metadatas[0]) > i {
					entry.Metadata = results.Metadatas[0][i]
				}
				if len(results.Distances) > 0 && len(results.Distances[0]) > i {
					entry.Distance = results.Distances[0][i]
				}
				entries = append(entries, entry)
			}
		}

		output, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to marshal result entries: %w", err)
		}

		duration := time.Since(startTime)
		logging.LogInfoByCtxf(serverCtx, "Tool 'search_documents' call finished. Returned %d bytes in %s.", len(output), duration)

		return mcp.NewToolResultText(string(output)), nil
	})

	// Tool: get_full_document
	getDocTool := mcp.NewTool("get_full_document",
		mcp.WithDescription("Retrieve the full document by its source path. All chunks belonging to the same source file are concatenated in order to reconstruct the original document."),
		mcp.WithString("source", mcp.Required(), mcp.Description("The source file path as stored in the metadata (e.g. '/path/to/file.md')")),
		mcp.WithString("collection", mcp.Description("Collection name (optional)")),
	)

	s.AddTool(getDocTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		logging.LogInfoByCtxf(serverCtx, "Tool 'get_full_document' call started.")

		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return nil, tracederrors.TracedErrorf("Invalid arguments type: %T", req.Params.Arguments)
		}

		source, ok := args["source"].(string)
		if !ok || source == "" {
			return nil, tracederrors.TracedErrorf("Missing or invalid 'source' argument")
		}

		colName := collectionName
		if c, ok := args["collection"].(string); ok && c != "" {
			colName = c
		}

		// 1. Get collection
		col, err := chromaClient.GetCollectionByName(serverCtx, colName)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get collection: %w", err)
		}

		// 2. Get all chunks for this source
		results, err := mcpGetBySource(chromaURL, col.ID, source)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to get document by source: %w", err)
		}

		if len(results.IDs) == 0 {
			duration := time.Since(startTime)
			logging.LogInfoByCtxf(serverCtx, "Tool 'get_full_document' call finished. No chunks found for source '%s'. Took %s.", source, duration)
			return mcp.NewToolResultText(fmt.Sprintf("No chunks found for source: %s", source)), nil
		}

		// 3. Sort chunks by ID to reconstruct original order (IDs are "doc_0", "doc_1", etc.)
		type indexedChunk struct {
			Index int
			Text  string
		}

		chunks := make([]indexedChunk, 0, len(results.IDs))
		for i, id := range results.IDs {
			idx := i
			if len(id) > 4 && id[:4] == "doc_" {
				if parsed, err := strconv.Atoi(id[4:]); err == nil {
					idx = parsed
				}
			}
			text := ""
			if i < len(results.Documents) {
				text = results.Documents[i]
			}
			chunks = append(chunks, indexedChunk{Index: idx, Text: text})
		}

		sort.Slice(chunks, func(i, j int) bool {
			return chunks[i].Index < chunks[j].Index
		})

		// 4. Concatenate all chunks
		var fullDocument string
		for _, chunk := range chunks {
			fullDocument += chunk.Text
		}

		type FullDocumentResult struct {
			Source     string `json:"source"`
			ChunkCount int    `json:"chunk_count"`
			Document   string `json:"document"`
		}

		output, err := json.MarshalIndent(FullDocumentResult{
			Source:     source,
			ChunkCount: len(chunks),
			Document:   fullDocument,
		}, "", "  ")
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to marshal full document result: %w", err)
		}

		duration := time.Since(startTime)
		logging.LogInfoByCtxf(serverCtx, "Tool 'get_full_document' call finished. Returned %d bytes (%d chunks) in %s.", len(output), len(chunks), duration)

		return mcp.NewToolResultText(string(output)), nil
	})

	// Tool: list_collections
	listTool := mcp.NewTool("list_collections",
		mcp.WithDescription("List all ChromaDB collections"),
	)

	s.AddTool(listTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		logging.LogInfoByCtxf(serverCtx, "Tool 'list_collections' call started.")

		url := fmt.Sprintf("%s%s/collections", chromaURL, v2BasePath)
		resp, err := http.Get(url)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to list collections: %w", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		duration := time.Since(startTime)
		logging.LogInfoByCtxf(serverCtx, "Tool 'list_collections' call finished. Returned %d bytes in %s.", len(body), duration)

		return mcp.NewToolResultText(string(body)), nil
	})

	// Start SSE HTTP transport on specified port
	sseServer := server.NewSSEServer(s)

	logging.LogInfoByCtxf(ctx, "MCP server listening on :%d", port)
	if err := sseServer.Start(":" + strconv.Itoa(port)); err != nil {
		return tracederrors.TracedErrorf("MCP SSE server failed on port %d: %w", port, err)
	}

	return nil
}
