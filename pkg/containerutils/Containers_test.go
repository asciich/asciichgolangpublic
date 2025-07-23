package containerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/containerutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func TestContainersIsRunningInsideContainer(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				require.False(mustutils.Must(containerutils.IsRunningInsideContainer(verbose)))
			},
		)
	}
}
