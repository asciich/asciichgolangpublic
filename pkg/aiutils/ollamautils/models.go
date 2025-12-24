package ollamautils

func GetFastModelName() string {
	return "TinyLlama" // no useful results, just for testing
}

func GetModerateSpeedModelName() string {
	return "Mistral"
}

func GetSlowModelName() string {
	return "llama3"
}

func GetImageProcessingModelName() string {
	return "llava-llama3:8b"
}

func GetSlowImageProcessingModelName() string {
	return "llama3.2-vision"
}

func GetOCRModelName() string {
	return "deepseek-ocr"
}

// More models:
// - deepseek-v2: https://ollama.com/library/deepseek-v2 8.9GB
// - deepseek-r1:1.5b https://ollama.com/library/deepseek-r1 1.1GB
// - deepseek-llm
// - llava for image description
