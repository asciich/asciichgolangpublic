# OpenWebUI CLI Commands

This package provides CLI commands for managing OpenWebUI.

## Commands

### run-as-docker-container

Starts OpenWebUI as a Docker container on the local machine.

**Usage:**
```bash
ai openwebui run-as-docker-container --port 8080 --container-name openwebui
```

**Options:**
- `--port`: Port for OpenWebUI to listen to (required)
- `--container-name`: Name of the Docker container (required)
- `--reachable-by-other-machines`: If set, binds to 0.0.0.0 instead of 127.0.0.1
- `--data-path`: Path to persist OpenWebUI data (optional)
- `--ollama-base-url`: Base URL of Ollama instance to connect to (optional)

**Examples:**

Basic usage:
```bash
ai openwebui run-as-docker-container --port 8080 --container-name openwebui
```

With data persistence:
```bash
ai openwebui run-as-docker-container \
  --port 8080 \
  --container-name openwebui \
  --data-path /var/lib/openwebui
```

With Ollama integration:
```bash
ai openwebui run-as-docker-container \
  --port 8080 \
  --container-name openwebui \
  --ollama-base-url http://localhost:11434
```

Accessible from other machines:
```bash
ai openwebui run-as-docker-container \
  --port 8080 \
  --container-name openwebui \
  --reachable-by-other-machines
```
