package ansiblegalaxyutils_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateFileStructure(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		err := ansiblegalaxyutils.CreateFileStructure(ctx, "", nil)
		require.Error(t, err)
	})

	t.Run("temp_dir", func(t *testing.T) {
		const verbose = true

		tempDir, err := tempfiles.CreateEmptyTemporaryDirectory(verbose)
		require.NoError(t, err)

		tempDirPath, err := tempDir.GetPath()
		require.NoError(t, err)

		err = ansiblegalaxyutils.CreateFileStructure(ctx, tempDirPath, &ansiblegalaxyutils.CreateCollectionFileStructureOptions{
			Namespace: "testnamespace",
			Name:      "example_collection",
			Version:   "v0.1.2",
			Authors:   []string{"exampleauthor"},
		})
		require.NoError(t, err)

		readmeExists, err := tempDir.FileInDirectoryExists(verbose, "README.md")
		require.NoError(t, err)
		require.True(t, readmeExists)

		ansibleBin := os.Getenv("ANSIBLE_GALAXY_BIN")
		if ansibleBin != "" {
			_, err = commandexecutorbashoo.Bash().RunCommand(
				commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
				&parameteroptions.RunCommandOptions{
					Command: []string{ansibleBin, "collection", "install", tempDirPath},
				},
			)
			require.NoError(t, err)
		}
	})
}

func Test_CreateFileStructureInDir(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		err := ansiblegalaxyutils.CreateFileStructureInDir(ctx, nil, nil)
		require.Error(t, err)
	})
}
