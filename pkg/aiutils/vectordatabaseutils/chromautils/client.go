package chromautils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const v2BasePath = "/api/v2/tenants/default_tenant/databases/default_database"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetBaseUrl() (string, error) {
	if c.baseURL == "" {
		return "", tracederrors.TracedError("baseURL not set")
	}

	return c.baseURL, nil
}

func (c *Client) CheckReachable(ctx context.Context) error {
	return c.Heartbeat(ctx)
}

// Heartbeat checks connectivity
func (c *Client) Heartbeat(ctx context.Context) error {
	url := c.baseURL + "/api/v2/heartbeat"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tracederrors.TracedErrorf("heartbeat on %s failed: status %d", url, resp.StatusCode)
	}

	logging.LogInfoByCtxf(ctx, "Chroma reachable. Heartbeat at %s requested successfully.", url)

	return nil
}

// CreateCollection creates or gets a collection
func (c *Client) CreateCollection(ctx context.Context, name string) (*Collection, error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create chroma collection '%s' on %s started.", name, baseUrl)

	body := map[string]interface{}{
		"name":          name,
		"get_or_create": true,
	}
	jsonData, _ := json.Marshal(body)

	fullUrl := baseUrl + v2BasePath + "/collections"

	resp, err := c.httpClient.Post(
		fullUrl,
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to post to create new collection '%s' on %s: %w", name, fullUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("create collection '%s' failed using %s (%d): %s", name, fullUrl, resp.StatusCode, respBody)
	}

	var col Collection
	if err := json.NewDecoder(resp.Body).Decode(&col); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode collection '%s' body returned by %s: %w", name, fullUrl, err)
	}

	logging.LogInfoByCtxf(ctx, "Create chroma collection '%s' on %s finised.", name, baseUrl)

	return &col, nil
}

// DeleteCollection deletes a collection by name
func (c *Client) DeleteCollection(ctx context.Context, name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete chroma collection '%s' on %s started.", name, baseUrl)

	fullUrl := baseUrl + v2BasePath + "/collections/" + name

	req, _ := http.NewRequest(http.MethodDelete, fullUrl, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logging.LogInfoByCtxf(ctx, "Delete chroma collection '%s' on %s finished.", name, baseUrl)

	return nil
}

// Add adds documents with embeddings to a collection
func (c *Client) Add(ctx context.Context, collectionID string, ids []string, embeddings [][]float32, documents []string, metadatas []map[string]interface{}) error {
	if collectionID == "" {
		return tracederrors.TracedErrorEmptyString("collectionID")
	}

	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add documents to with embeddings to the collection id='%s' on chroma %s started.", collectionID, baseUrl)

	body := map[string]any{
		"ids":        ids,
		"embeddings": embeddings,
		"documents":  documents,
		"metadatas":  metadatas,
	}
	jsonData, _ := json.Marshal(body)

	fullUrl := baseUrl + v2BasePath + "/collections/" + collectionID + "/add"

	resp, err := c.httpClient.Post(fullUrl, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return tracederrors.TracedErrorf("add failed (%d) using %s: %s", resp.StatusCode, fullUrl, respBody)
	}

	logging.LogInfoByCtxf(ctx, "Add documents to with embeddings to the collection id='%s' on chroma %s finished.", collectionID, baseUrl)

	return nil
}

// QueryResult holds query results
type QueryResult struct {
	IDs       [][]string                 `json:"ids"`
	Documents [][]string                 `json:"documents"`
	Metadatas [][]map[string]interface{} `json:"metadatas"`
	Distances [][]float32                `json:"distances"`
}

// Query searches the collection
func (c *Client) Query(collectionID string, queryEmbeddings [][]float32, nResults int) (*QueryResult, error) {
	if collectionID == "" {
		return nil, tracederrors.TracedErrorEmptyString("collectionID")
	}

	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"query_embeddings": queryEmbeddings,
		"n_results":        nResults,
		"include":          []string{"documents", "metadatas", "distances"},
	}
	jsonData, _ := json.Marshal(body)

	fullUrl := baseUrl + v2BasePath + "/collections/" + collectionID + "/query"
	resp, err := c.httpClient.Post(fullUrl, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to queryy collection using %s: %w", fullUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("query failed using %s (%d): %s", resp.StatusCode, fullUrl, respBody)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode QueryResult: %w", err)
	}
	return &result, nil
}

func (c *Client) GetCollectionByName(ctx context.Context, name string) (*Collection, error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/collections/%s", baseUrl, v2BasePath, name)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get collection '%s': %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("Get collection '%s' failed (%d): %s", name, resp.StatusCode, string(body))
	}

	var col Collection
	if err := json.NewDecoder(resp.Body).Decode(&col); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode collection response: %w", err)
	}

	return &col, nil
}
