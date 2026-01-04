package osutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CommandAvailable(t *testing.T) {
	tests := []struct {
		command   string
		imageName string
		expected  bool
	}{
		{"which", "ubuntu", true},
		{"vim", "ubuntu", false},
		{"whereis", "ubuntu", true},
		{"which", "archlinux", false},
		{"vim", "archlinux", false},
		{"whereis", "archlinux", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				container, err := nativedocker.NewDocker().RunContainer(
					ctx,
					&dockeroptions.DockerRunContainerOptions{
						Name:      "test-command-available",
						ImageName: tt.imageName,
						Command:   []string{"sleep", "1m"},
					},
				)
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				isCommandAvailable, err := osutils.IsCommandAvailable(ctx, container, tt.command)
				require.NoError(t, err)
				require.EqualValues(t, tt.expected, isCommandAvailable)
			},
		)
	}
}
