package kubernetesutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/filesutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
)

func IsInClusterAuthenticationAvailable(ctx context.Context) bool {
	for _, f := range []string{"/var/run/secrets/kubernetes.io/serviceaccount/token", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"} {
		if !filesutils.IsFile(contextutils.WithVerbosityContextByBool(ctx, false), f) {
			logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is not available.")
			return false
		}
	}

	logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is available.")
	return true
}
