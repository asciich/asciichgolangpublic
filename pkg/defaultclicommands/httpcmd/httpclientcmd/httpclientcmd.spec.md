# httpclientcmd specifications

This are the specifications for the [`httpclientcmd` package](README.md).

This document extends the [constitution.md](/constitution.md).

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
