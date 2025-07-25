package kubernetesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func IsInClusterAuthenticationAvailable(ctx context.Context) bool {
	for _, f := range []string{"/var/run/secrets/kubernetes.io/serviceaccount/token", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"} {
		if !nativefiles.IsFile(contextutils.WithSilent(ctx), f) {
			logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is not available.")
			return false
		}
	}

	logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is available.")
	return true
}
