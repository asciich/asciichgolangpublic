package aidercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiderutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"os"
)

func NewBuildContainerImageCmd() *cobra.Command {
	const short = "Build a local docker container image to run aider."

	cmd := &cobra.Command{
		Use:   "build-container-image",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai aider build-container-image`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(aiderutils.BuildAiderDockerContainer(ctx))
		},
	}

	return cmd
}
