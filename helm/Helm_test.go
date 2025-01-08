package helm

import (
	"testing"

	"github.com/asciich/asciichgolangpublic"
)

func getHelmImplementationByName(implementationName string) (helm Helm) {
	if implementationName == "commandExecutorHelm" {
		return MustGetLocalCommandExecutorHelm()
	}

	asciichgolangpublic.LogFatalf("Unknwon implementation name '%s'", implementationName)

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
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				// assert := assert.New(t)

				const verbose bool = true

				kubernetes := getHelmImplementationByName(tt.implementationName)
				kubernetes.MustAddRepositoryByName("argo", "https://argoproj.github.io/argo-helm", verbose)
			},
		)
	}
}
