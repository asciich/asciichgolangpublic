# Openhandscmd specifications

This are the specifications for the [`openhandscmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../../defaultclicommands.spec.md).

## Implementation

- The general `configure` command can be used to configure the LLM `--name`:
    - It takes an `--openai-base-url` which can and should include a port. Alternative this info can be read from the env var `OPENAI_BASE_URL`.
    - It reads the API token from `OPENAI_API_KEY`.
