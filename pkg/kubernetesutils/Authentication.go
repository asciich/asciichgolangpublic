package kubernetesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
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
