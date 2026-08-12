package copilotcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/copilotutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunInCopilotCliContainerCmd() *cobra.Command {
	const short = "Run a prompt in an isolated copilot CLI container."

	cmd := &cobra.Command{
		Use:   "run-in-copilot-cli-container",
		Short: short,
		Long: short + `

Spins up a Docker container with GitHub Copilot CLI running as an agent.
The container provides an isolated environment for executing prompts.

Authentication is done via the COPILOT_GITHUB_TOKEN environment variable
which must be set before running this command.

Example:
  ` + os.Args[0] + ` ai copilot run-in-copilot-cli-container \
    --prompt='Refactor main.go to use dependency injection' \
    --system-prompt='You are a helpful coding assistant.' \
    --workspace-path='/home/user/project' \
    --git-user-name='John Doe' \
    --git-user-email='john@example.com' \
    --git-config-path='~/.gitconfig' \
    --ssh-config-path='~/.ssh/config' \
    --ssh-agent \
    --verbose

Example without workspace (uses temporary workspace inside container):
  ` + os.Args[0] + ` ai copilot run-in-copilot-cli-container \
    --prompt='Write a hello world in Go' \
    --verbose
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			prompt, err := cmd.Flags().GetString("prompt")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			systemPrompt, err := cmd.Flags().GetString("system-prompt")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			workspacePath, err := cmd.Flags().GetString("workspace-path")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			containerImage, err := cmd.Flags().GetString("container-image")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			gitUserName, err := cmd.Flags().GetString("git-user-name")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			gitUserEmail, err := cmd.Flags().GetString("git-user-email")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			gitConfigPath, err := cmd.Flags().GetString("git-config-path")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			sshConfigPath, err := cmd.Flags().GetString("ssh-config-path")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			sshAgent, err := cmd.Flags().GetBool("ssh-agent")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if prompt == "" {
				logging.LogFatal("Prompt is not set. Use --prompt to specify the prompt to execute.")
			}

			options := copilotutils.NewRunCopilotCliContainerOptions()
			mustutils.Must0(options.SetPrompt(prompt))

			if systemPrompt != "" {
				mustutils.Must0(options.SetSystemPrompt(systemPrompt))
			}

			if workspacePath != "" {
				mustutils.Must0(options.SetWorkspacePath(workspacePath))
			}

			if containerImage != "" {
				mustutils.Must0(options.SetContainerImage(containerImage))
			}

			if gitUserName != "" {
				mustutils.Must0(options.SetGitUserName(gitUserName))
			}

			if gitUserEmail != "" {
				mustutils.Must0(options.SetGitUserEmail(gitUserEmail))
			}

			if gitConfigPath != "" {
				mustutils.Must0(options.SetGitConfigPath(gitConfigPath))
			}

			if sshConfigPath != "" {
				mustutils.Must0(options.SetSshConfigPath(sshConfigPath))
			}

			options.SetSSHAgent(sshAgent)

			output := mustutils.Must(copilotutils.RunCopilotCliContainer(options))

			logging.LogInfoByCtxf(ctx, "Copilot CLI output:\n%s", output)
			logging.LogGoodByCtxf(ctx, "Copilot CLI container finished.")
		},
	}

	cmd.Flags().String("prompt", "", "The prompt to execute in the copilot CLI agent (required).")
	cmd.Flags().String("system-prompt", "", "The system prompt for the copilot CLI agent (optional).")
	cmd.Flags().String("workspace-path", "", "Path to mount as workspace inside the container. If not set a temporary workspace is used.")
	cmd.Flags().String("container-image", "", "Container image to use. Defaults to ghcr.io/henrybravo/docker-sandbox-run-copilot.")
	cmd.Flags().String("git-user-name", "", "Git username to use inside the container (optional).")
	cmd.Flags().String("git-user-email", "", "Git email to use inside the container (optional).")
	cmd.Flags().String("git-config-path", "", "Path to git config file to mount read-only inside the container. E.g. ~/.gitconfig (optional).")
	cmd.Flags().String("ssh-config-path", "", "Path to SSH config file to mount read-only inside the container. E.g. ~/.ssh/config (optional).")
	cmd.Flags().Bool("ssh-agent", false, "Pass SSH_AUTH_SOCK into the container for git clone operations.")

	return cmd
}
