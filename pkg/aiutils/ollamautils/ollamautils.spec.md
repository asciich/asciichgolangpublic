# ollamautils specifications

## Implementation

- To run Ollama:
    - `RunCpuOnly(ctx context.Context) error` runs ollama in docker with CPU support only.
    - `RunGPUNvidia(ctx context.Context) error` runs ollama in docker with nvidia GPU support.
    - `RunGPUAmd(ctx context.Context) error` runs ollama in docker with amd GPU support.
        - Uses the dedicated `ollama/ollama:rocm` image for ROCm/AMD support, whereas all other run commands use the default `ollama/ollama` image.
    - `RunGPU(ctx context.Context) error` runs in docker with autodetected GPU support.
        - If no GPU detected return an error. Falling back to CPU is no option.
    - All run commands use the same volume mount `ollama:/root/.ollama`.
