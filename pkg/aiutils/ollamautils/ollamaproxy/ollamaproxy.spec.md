# ollamaproxy specifications

## Implementation

This is a proxy in front of Ollama:
- Authentication: 
    - Uses the env var OLLAMA_API_KEY as the shared API key; clients send it OpenAI-style as Authorization: Bearer <key>, and any missing or wrong key is answered with 401 in OpenAI’s error format.
- Allow only one request to Ollama at a time.
- A Prometheus `/metrics` page needs to be exposed:
    - All metrics have the prefix `ollama_proxy_`
    - Number of currently waiting prompts.
    - `ollama_proxy_promt_running` is 1 if there is an active promt currently processed, 0 otherwise.
    - `ollama_proxy_promt_counter` counting all promts performed.

## Testing

- Use unit tests to validate the implementation.
