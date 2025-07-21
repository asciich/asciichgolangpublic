package helmutils_test

import (
	"context"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helminterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getHelmImplementationByName(implementationName string) (helm helminterfaces.Helm) {
	if implementationName == "commandExecutorHelm" {
		return mustutils.Must(helmutils.GetLocalCommandExecutorHelm())
	}

	logging.LogFatalf("Unknown implementation name '%s'", implementationName)

	return nil
}

func TestRole_AddHelmRepo(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorHelm"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				kubernetes := getHelmImplementationByName(tt.implementationName)
				mustutils.Must0(kubernetes.AddRepositoryByName(ctx, "argo", "https://argoproj.github.io/argo-helm"))
			},
		)
	}
}
