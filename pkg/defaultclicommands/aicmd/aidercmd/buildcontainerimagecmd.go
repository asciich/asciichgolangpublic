package aidercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiderutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewBuildContainerImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "build-container-image",
		Short: "Build a local docker container image to run aider.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(aiderutils.BuildAiderDockerContainer(ctx))
		},
	}

	return cmd
}