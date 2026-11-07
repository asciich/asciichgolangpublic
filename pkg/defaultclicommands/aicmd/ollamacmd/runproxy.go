package ollamacmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils/ollamaproxy"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunProxyCmd() *cobra.Command {
	const short = "Run a proxy in front of an ollama server."

	cmd := &cobra.Command{
		Use:   "run-proxy",
		Short: short,
		Long: short + `

The proxy adds authentication, request serialization and monitoring on top of a
plain ollama server:
  - Authentication: the shared API key has to be sent by the clients OpenAI style
    as 'Authorization: Bearer <key>'. Every request with a missing or wrong key
    is answered with HTTP 401 using OpenAI's error format.
    The key is provided either directly by the environment variable
    ` + ollamaproxy.EnvVarNameOllamaApiKey + ` or by a file referenced by the environment variable
    ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + `. The key file path can also be given by '--api-key-file'.
    The key is read from one of the following sources (first match wins):
      1. The file given by '--api-key-file'.
      2. The environment variable ` + ollamaproxy.EnvVarNameOllamaApiKey + `.
      3. The file whose path is given by the environment variable
         ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + `.
    Leading and trailing whitespace (like a trailing newline) in the key file is
    ignored, so docker and kubernetes secrets can be used directly.
    Setting both ` + ollamaproxy.EnvVarNameOllamaApiKey + ` and ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + ` is ambiguous and
    results in an error (unless '--api-key-file' is set, which takes precedence).
    The key itself is intentionally not configurable by a flag so it does not
    show up in the process list or in the shell history.
  - Only one request at a time is forwarded to the ollama server. All other
    requests are queued and processed as soon as the ollama server is free
    again. Streaming responses ('"stream": true') are passed through without
    buffering, but occupy the single slot until the stream is finished.
  - Prometheus metrics are exposed on '/metrics' (configurable with
    '--metrics-path'). All metrics use the prefix 'ollama_proxy_', especially
    the number of currently waiting prompts, 'ollama_proxy_promt_running' and
    'ollama_proxy_promt_counter'. The metrics endpoint is not protected by the
    API key so prometheus can scrape it without additional configuration.

Environment variables:
  ` + ollamaproxy.EnvVarNameOllamaApiKey + `
      The shared API key the clients have to send as bearer token.
  ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + `
      Path to a file containing the shared API key. Alternative to
      ` + ollamaproxy.EnvVarNameOllamaApiKey + `, handy for docker and kubernetes secrets.
      Must not be set together with ` + ollamaproxy.EnvVarNameOllamaApiKey + `.

The proxy runs until it is terminated by CTRL-C or SIGTERM and is shut down
gracefully afterwards.

Usage:
    # Key directly from the environment variable ` + ollamaproxy.EnvVarNameOllamaApiKey + `:
    export ` + ollamaproxy.EnvVarNameOllamaApiKey + `='your-secret-key'
    ` + os.Args[0] + ` ai ollama run-proxy

    # Key from a file referenced by the environment variable ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + `:
    export ` + ollamaproxy.EnvVarNameOllamaApiKeyFile + `='/run/secrets/ollama_api_key'
    ` + os.Args[0] + ` ai ollama run-proxy

    # Key from a file given by flag (takes precedence over both environment variables):
    ` + os.Args[0] + ` ai ollama run-proxy --api-key-file /run/secrets/ollama_api_key

    # Other options:
    ` + os.Args[0] + ` ai ollama run-proxy --listen-address ":8080"
    ` + os.Args[0] + ` ai ollama run-proxy --target-url "http://ollama.internal:11434"
    ` + os.Args[0] + ` ai ollama run-proxy --metrics-path "/internal/metrics" --upstream-timeout 10m`,

		Args: cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			listenAddress := mustutils.Must(cmd.Flags().GetString("listen-address"))
			targetUrl := mustutils.Must(cmd.Flags().GetString("target-url"))
			metricsPath := mustutils.Must(cmd.Flags().GetString("metrics-path"))
			apiKeyFile := mustutils.Must(cmd.Flags().GetString("api-key-file"))
			upstreamTimeout := mustutils.Must(cmd.Flags().GetDuration("upstream-timeout"))
			forwardAuthorizationHeader := mustutils.Must(cmd.Flags().GetBool("forward-authorization-header"))

			logging.LogInfoByCtxf(ctx, "Ollama proxy is going to listen on '%s' and forwards to '%s'.", listenAddress, targetUrl)

			mustutils.Must0(ollamaproxy.RunOllamaProxy(ctx, &ollamaproxy.OllamaProxyOptions{
				ListenAddress:              listenAddress,
				TargetUrl:                  targetUrl,
				ApiKeyFile:                 apiKeyFile,
				MetricsPath:                metricsPath,
				UpstreamTimeout:            upstreamTimeout,
				ForwardAuthorizationHeader: forwardAuthorizationHeader,
			}))

			logging.LogGoodByCtxf(ctx, "Ollama proxy stopped.")
		},
	}

	cmd.Flags().String(
		"listen-address",
		ollamaproxy.DefaultOllamaProxyListenAddress,
		"Address the proxy listens on.",
	)

	cmd.Flags().String(
		"target-url",
		ollamaproxy.DefaultOllamaProxyTargetUrl,
		"Base url of the ollama server to forward the requests to.",
	)

	cmd.Flags().String(
		"api-key-file",
		"",
		"Path to a file containing the shared API key. Takes precedence over the environment variables "+
			ollamaproxy.EnvVarNameOllamaApiKey+" and "+ollamaproxy.EnvVarNameOllamaApiKeyFile+". "+
			"Only the path is passed as flag, the key itself never is.",
	)

	cmd.Flags().String(
		"metrics-path",
		ollamaproxy.DefaultOllamaProxyMetricsPath,
		"Path the prometheus metrics are exposed on.",
	)

	cmd.Flags().Duration(
		"upstream-timeout",
		0,
		"Timeout for a single request forwarded to ollama. When 0 no timeout is applied which is recommended for long running prompts.",
	)

	cmd.Flags().Bool(
		"forward-authorization-header",
		false,
		"Forward the incoming 'Authorization' header to the ollama server as well. By default the header is stripped since the shared key is only valid between client and proxy.",
	)

	return cmd
}
