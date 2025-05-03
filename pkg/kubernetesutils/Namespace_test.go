package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getKubernetesByImplementationName(ctx context.Context, implementationName string) kubernetesutils.KubernetesCluster {
	if implementationName == "commandExecutorKubernetes" {
		// Directly call kind binary to avoid cyclic import...
		commandexecutor.Bash().RunOneLiner(ctx, "kind create cluster -n kind || true")

		return kubernetesutils.MustGetLocalCommandExecutorKubernetesByName("kind-kind")
	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
		implementationName,
	)

	return nil
}

func TestNamespace_CreateAndDeleteNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				err := kubernetes.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(kubernetes.NamespaceByNameExists(ctx, namespaceName)))

				for i := 0; i < 2; i++ {
					_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)
					require.True(t, mustutils.Must(kubernetes.NamespaceByNameExists(ctx, namespaceName)))
				}

				for i := 0; i < 2; i++ {
					err := kubernetes.DeleteNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(kubernetes.NamespaceByNameExists(ctx, namespaceName)))
				}
			},
		)
	}
}
