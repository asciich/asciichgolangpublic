package ollamautils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Send a single prompt (no conversation) to a ollama server and return the complete response.
//
// This function is not streaming the result, it waits until the answer is complete and returns it as a whole.
func SendPrompt(ctx context.Context, prompt string, options *PromptOptions) (string, error) {
	if prompt == "" {
		return "", tracederrors.TracedErrorEmptyString("prompt")
	}

	if options == nil {
		options = &PromptOptions{}
	}

	logging.LogInfoByCtxf(ctx, "Send ollama prompt started.")

	url, err := options.GetGenerateUrl(ctx)
	if err != nil {
		return "", err
	}

	modelName, err := options.GetModelNameOrDefault(ctx)
	if err != nil {
		return "", err
	}

	// Request structure for Ollama API
	type OllamaRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		System string `json:"system"`
		Stream bool   `json:"stream"`
	}

	// Response structure for Ollama API
	type OllamaResponse struct {
		Model    string `json:"model"`
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	requestBody := OllamaRequest{
		Model:  modelName,
		Prompt: prompt,
		System: "", // can be used to add agent instructions
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", tracederrors.TracedErrorf("Error marshaling request: %v\n", err)
	}

	tStart := time.Now()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", tracederrors.TracedErrorf("Error sending request: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", tracederrors.TracedErrorf("Error reading response: %v\n", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", tracederrors.TracedErrorf("Error unmarshaling response: %v\n", err)
	}

	duration := time.Since(tStart)

	response := strings.TrimSpace(ollamaResp.Response)

	logging.LogInfoByCtxf(ctx, "Send ollama prompt finished. The prompt execution took '%v'", duration)

	return response, nil
}
