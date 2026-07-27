package ollamaagent

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type OllamaFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type OllamaTool struct {
	Type     string         `json:"type"`
	Function OllamaFunction `json:"function"`
}

type OllamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content"`
	ToolCalls []OllamaToolCall `json:"tool_calls,omitempty"`
}

type OllamaToolCall struct {
	Function OllamaToolCallFunction `json:"function"`
}

type OllamaToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Tools    []OllamaTool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
}

type OllamaChatResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

type OllamaClient struct {
	baseURL string
	model   string
	client  *http.Client
}

func NewOllamaClient(baseURL string, model string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

func (o *OllamaClient) Chat(messages []OllamaMessage, tools []OllamaTool) (*OllamaChatResponse, error) {
	reqBody := OllamaChatRequest{
		Model:    o.model,
		Messages: messages,
		Tools:    tools,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal chat request: %w", err)
	}

	resp, err := o.client.Post(o.baseURL+"/api/chat", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to post to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, tracederrors.TracedErrorf("Ollama chat failed (%d): %s", resp.StatusCode, body)
	}

	var chatResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decode ollama response: %w", err)
	}

	return &chatResp, nil
}
