package helmutils_test

import (
	"testing"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helminterfaces"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getHelmImplementationByName(implementationName string) (helm helminterfaces.Helm) {
	if implementationName == "commandExecutorHelm" {
		return helmutils.MustGetLocalCommandExecutorHelm()
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
				const verbose bool = true

				kubernetes := getHelmImplementationByName(tt.implementationName)
				kubernetes.MustAddRepositoryByName("argo", "https://argoproj.github.io/argo-helm", verbose)
			},
		)
	}
}
