package commandexecutorkubernetes

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// WhoAmI returns the identity used to authenticate against the cluster.
//
// Strategy:
//  1. Ask the API server using `kubectl auth whoami` (SelfSubjectReview API).
//     This is the most reliable way since the server reports the identity.
//  2. If the SelfSubjectReview API is not available (cluster < 1.27/1.28 or
//     kubectl < 1.26), fall back to reading the identity from the kubeconfig:
//     the CN of the client certificate, or the kubeconfig user name.
//
// Errors unrelated to a missing API (e.g. cluster not reachable, unauthorized)
// are returned as they are and do NOT trigger the fallback.
//
// The cluster name is optional (e.g. for the default cluster). If no name and no
// cached context is set, kubectl uses the current-context of the kubeconfig.
func (c *CommandExecutorKubernetes) WhoAmI(ctx context.Context) (*kubernetesgeneric.UserInfo, error) {
	executor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	contextArgs, contextDisplayName, err := c.kubectlContextArgs(ctx)
	if err != nil {
		return nil, err
	}

	clusterName := c.displayName()

	userName, supported, err := c.whoAmIBySelfSubjectReview(ctx, executor, contextArgs)
	if err != nil {
		return nil, err
	}

	if !supported {
		logging.LogInfoByCtxf(
			ctx,
			"Whoami: SelfSubjectReview API not available in cluster '%s'. Fall back to kubeconfig based detection.",
			clusterName,
		)

		userName, err = c.whoAmIByKubeconfig(ctx, executor, contextArgs, contextDisplayName)
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Whoami: Kube context '%s' uses user '%s' to log in to cluster '%s'.", contextDisplayName, userName, clusterName)

	return &kubernetesgeneric.UserInfo{
		Username: userName,
	}, nil
}

// displayName returns the cluster name for log and error messages.
// The default cluster has no explicit name, so a placeholder is returned instead of an error.
func (c *CommandExecutorKubernetes) displayName() string {
	if c.name == "" {
		return "<default cluster>"
	}

	return c.name
}

// kubectlContextArgs returns the kubectl arguments to select the context together
// with a human readable name of the context used for log and error messages.
//
//   - In-cluster authentication: no context is needed (same as in the other methods).
//   - Cluster name or cached context name set: "--context <name>".
//   - Neither set (default cluster): no context argument. kubectl then uses the
//     current-context of the kubeconfig.
func (c *CommandExecutorKubernetes) kubectlContextArgs(ctx context.Context) (args []string, displayName string, err error) {
	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. cluster context is not used.")
		return []string{}, "<in-cluster>", nil
	}

	if c.name == "" && c.cachedContextName == "" {
		logging.LogInfoByCtxf(ctx, "No cluster name set. The current-context of the kubeconfig is used.")
		return []string{}, "<current-context>", nil
	}

	kubeContext, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, "", err
	}

	return []string{"--context", kubeContext}, kubeContext, nil
}

// whoAmIBySelfSubjectReview uses `kubectl auth whoami`.
// supported is false if the API (or the kubectl subcommand) is not available.
// In this case err is nil and the caller is expected to use a fallback.
func (c *CommandExecutorKubernetes) whoAmIBySelfSubjectReview(
	ctx context.Context,
	executor commandexecutorinterfaces.CommandExecutor,
	contextArgs []string,
) (userName string, supported bool, err error) {
	cmd := []string{"kubectl"}
	cmd = append(cmd, contextArgs...)
	cmd = append(cmd, "auth", "whoami", "-ojson")

	output, err := executor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command:           cmd,
		AllowAllExitCodes: true,
	})
	if err != nil {
		return "", false, err
	}

	if !output.IsExitSuccess() {
		stderr, err := output.GetStderrAsString()
		if err != nil {
			return "", false, err
		}

		notSupportedHints := []string{
			"selfsubjectreviews API is not enabled", // kubectl >= 1.27 against a cluster without the API
			"could not find the requested resource", // HTTP 404 for the endpoint
			"unknown command",                       // kubectl < 1.26 has no `auth whoami`
		}

		for _, hint := range notSupportedHints {
			if strings.Contains(stderr, hint) {
				return "", false, nil
			}
		}

		return "", false, tracederrors.TracedErrorf("`kubectl auth whoami` failed: %s", strings.TrimSpace(stderr))
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return "", false, err
	}

	userName, err = jsonutils.RunJqAgainstJsonStringAsString(stdout, ".status.userInfo.username")
	if err != nil {
		return "", false, err
	}

	return userName, true, nil
}

// whoAmIByKubeconfig determines the user from the kubeconfig of the selected context.
//
// If the user authenticates with a client certificate (default for kind and
// kubeadm clusters) the certificate's CN is returned, which is the username the API
// server sees. Otherwise the user name configured in the kubeconfig is returned.
// This is only the identity the client *claims*, not one confirmed by the server.
func (c *CommandExecutorKubernetes) whoAmIByKubeconfig(
	ctx context.Context,
	executor commandexecutorinterfaces.CommandExecutor,
	contextArgs []string,
	contextDisplayName string,
) (string, error) {
	cmd := []string{"kubectl"}
	cmd = append(cmd, contextArgs...)
	cmd = append(cmd, "config", "view", "--raw", "--minify", "-o", "json")

	// --raw is needed to get the certificate data. The silent context prevents
	// the output (which may contain secrets) from being logged.
	stdout, err := executor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return "", err
	}

	var kubeconfig struct {
		Users []struct {
			Name string `json:"name"`
			User struct {
				ClientCertificateData string `json:"client-certificate-data"`
			} `json:"user"`
		} `json:"users"`
	}

	err = json.Unmarshal([]byte(stdout), &kubeconfig)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to parse kubeconfig of context '%s': %w", contextDisplayName, err)
	}

	if len(kubeconfig.Users) == 0 {
		return "", tracederrors.TracedErrorf(
			"No user found in kubeconfig for context '%s'. The identity cannot be determined without the SelfSubjectReview API.",
			contextDisplayName,
		)
	}

	user := kubeconfig.Users[0]

	if user.User.ClientCertificateData != "" {
		cn, err := commonNameFromBase64Certificate(user.User.ClientCertificateData)
		if err != nil {
			return "", tracederrors.TracedErrorf("Unable to get CN from client certificate of kubeconfig user '%s': %w", user.Name, err)
		}

		if cn != "" {
			return cn, nil
		}
	}

	logging.LogInfoByCtxf(
		ctx,
		"Whoami: No client certificate with CN found for kubeconfig user '%s'. Using the kubeconfig user name.",
		user.Name,
	)

	return user.Name, nil
}

// commonNameFromBase64Certificate returns the subject CN of a base64 encoded
// PEM certificate as used in the kubeconfig field `client-certificate-data`.
func commonNameFromBase64Certificate(b64 string) (string, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", fmt.Errorf("no PEM block found in client certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	return cert.Subject.CommonName, nil
}
