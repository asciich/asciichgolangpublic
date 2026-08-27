# copilotutils specifications

This are the specifications for the [`copilotutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

- `RunCopilotCliContainer` can be used to spin a container to running copilot CLI as agent in an isolated way.
    - The argument of type `RunCopilotCliContainerOptions` provides at least:
        - `WorkspacePath`. This path is mounted as workspace and is optional. If not set the whole agent runs on a temporary workspace inside the container.
        - `SSHAgent`. If set to true the `SSH_AUTH_SOCK` env var is passed into the container. Furthermore the socket is mounted inside the container as well, this way the ssh-agent can be used inside the container and is available for clone git repos.
        - `SystemPrompt`: The system prompt. This is an optional value.
        - `Prompt`: The prompt to execute.
        - `ContainerImage`: The container image to use. Defaults to `ghcr.io/henrybravo/docker-sandbox-run-copilot`.
        - `GitUserName`: Optionally your git username.
        - `GitUserEmail`: Optionally your git email.
        - `GitConfigPath`: Optionally the full git config to mount. Is mounted read only. Can point to your global one `~/.gitconfig`.
        - `SshConfigPath`: Optionally the full ssh config to mount. Is mounted read only. Can point to your global one `~/.ssh/config`
    - Pass `COPILOT_GITHUB_TOKEN` for authentication. If not set return an error and do not even start the container.
    - Use a generated container name: `copilot-cli-<8randomchars>`.
    - Default user to use is `agent`, the default home directory is `/home/agent`.
