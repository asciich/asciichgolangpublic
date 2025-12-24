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

func sendPromptToSingleimage(ctx context.Context, prompt string, imagePath string, options *PromptOptions) (string, error) {
	if imagePath == "" {
		return "", tracederrors.TracedErrorEmptyString("imagePath")
	}

	if prompt == "" {
		return "", tracederrors.TracedErrorEmptyString("prompt")
	}

	if options == nil {
		options = new(PromptOptions)
	}

	optionsToUse := options.GetDeepCopy()

	if optionsToUse.ModelName == "" {
		optionsToUse.ModelName = GetImageProcessingModelName()
		logging.LogInfoByCtxf(ctx, "Set LLM model for image description to default '%s'.", optionsToUse.ModelName)
	}

	optionsToUse.ImagePaths = []string{imagePath}

	response, err := SendPrompt(
		ctx,
		prompt,
		optionsToUse,
	)
	if err != nil {
		return "", err
	}

	return response, nil
}

// Sends a prompt do tescribe an image to the ollama server and returns the complete description/ response.
func DescribeImage(ctx context.Context, imagePath string, options *PromptOptions) (string, error) {
	logging.LogInfoByCtxf(ctx, "Describe image using ollama started.")

	description, err := sendPromptToSingleimage(ctx, "Describe the image", imagePath, options)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Describe image using ollama finished.")

	return description, nil
}

// Sends a prompt do detect the characters in an image to the ollama server and returns them.
func OpticalCharacterRecognition(ctx context.Context, imagePath string, options *PromptOptions) (string, error) {
	logging.LogInfoByCtxf(ctx, "OCR using ollama started.")

	description, err := sendPromptToSingleimage(ctx, "Perform an optical character recognition by only outputting the detected characters.", imagePath, options)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "OCR image using ollama finished.")

	return description, nil
}

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

	// OllamaOptions contains model parameters
	type OllamaOptions struct {
		Temperature float64 `json:"temperature"`
	}

	// Request structure for Ollama API
	type OllamaRequest struct {
		Model   string        `json:"model"`
		Prompt  string        `json:"prompt"`
		System  string        `json:"system"`
		Stream  bool          `json:"stream"`
		Images  []string      `json:"images"`
		Options OllamaOptions `json:"options,omitempty"`
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

	if len(options.ImagePaths) >= 1 {
		requestBody.Images, err = options.GetImagesAsBase64Slice(ctx)
		if err != nil {
			return "", err
		}
		logging.LogInfoByCtxf(ctx, "Added '%d' images to the request.", len(options.ImagePaths))
	} else {
		logging.LogInfoByCtxf(ctx, "No images are loaded and send as part of the reuqest.")
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
