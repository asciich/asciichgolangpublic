# httpclientcmd specifications

## Implementation

- The default scheme to use is `https://`:
    - This should behave the same way:
        ```bash
        asciichgolangpublic http client get example.com # https:// must be automatically used if not specified.
        ```
    - As this does:
        ```bash
        asciichgolangpublic http client get https://example.com
        ```
    - Applies to all commands, not only `get`.
