# openhandscmd

OpenHands related commands.

## Specifications

For specifications see [openhandscmd.spec.md](openhandscmd.spec.md)

## Commands

### `configure`

Configure openhands for a generic OpenAI-compatible LLM.

```
Usage:
    <command> ai openhands configure --name <profile-name> --model <model-name> --openai-base-url <url>
```

Reads the API token from `OPENAI_API_KEY` env var.
The base URL can be specified via `--openai-base-url` or the `OPENAI_BASE_URL` env var.

Flags:
- `--url` - URL to openhands instance. E.g: http://localhost:8000 (required)
- `--name` - Name for the LLM profile (required)
- `--model` - Model name to use. E.g: gpt-4, claude-3, etc. (required)
- `--openai-base-url` - Base URL for the OpenAI-compatible API. Can also be read from OPENAI_BASE_URL env var.

### `configure-google-ai-studio`

Configure openhands for Google AI Studio LLM.

```
Usage:
    <command> ai openhands configure-google-ai-studio --url <url>
```

Needs the env var `GOOGLE_AI_STUDIO_API_KEY` set.

Flags:
- `--url` - URL to openhands. E.g: http://localhost:8000 (required)

### `configure-infomaniak`

Configure openhands for Infomaniak LLM.

```
Usage:
    <command> ai openhands configure-infomaniak --url <url> --product-id <id>
```

Needs the env var `INFOMANIAK_API_KEY` set.

Flags:
- `--url` - URL to openhands. E.g: http://localhost:8000 (required)
- `--product-id` - Infomaniak product ID. E.g: 12345 (required)

### `configure-swisscom-myai`

Configure openhands for myAI LLM of Swisscom.

```
Usage:
    <command> ai openhands configure-swisscom-myai --url <url>
```

Needs the env var `SWISSCOM_MYAI_API_KEY` set.

Flags:
- `--url` - URL to openhands. E.g: http://localhost:8000 (required)

### `run-as-docker-container`

Runs openhands as docker container on the local machine.

```
Usage:
    <command> ai openhands run-as-docker-container --container-name <name> --workspace-path <path> --port <port>
```

Flags:
- `--port` - Port for openhands to listen to (required)
- `--container-name` - Name of the docker container (required)
- `--reachable-by-other-machines` - If set, binds to 0.0.0.0 instead of 127.0.0.1 to allow access from other machines in the network
- `--workspace-path` - Path to the workspace directory to mount into the container (default: .)
