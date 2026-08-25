package gitutils_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegitoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_GetRepositoryRootPathByPath(t *testing.T) {
	absPath, err := filepath.Abs(".")
	require.NoError(t, err)

	tests := []struct {
		Path string
	}{
		{"."},
		{absPath},
	}

	for _, tt := range tests {
		ctx := getCtx()

		repoRootPath, err := gitutils.GetRepositoryRootPathByPath(ctx, tt.Path)
		require.NoError(t, err)
		require.EqualValues(t, "asciichgolangpublic", filepath.Base(repoRootPath))
	}
}

func getGitRepoByImplementationNameAndPath(t *testing.T, implementationName string, path string) gitinterfaces.GitRepository {
	if implementationName == "nativegitoo" {
		ret, err := nativegitoo.NewGitRepositoryFromPath(path)
		require.NoError(t, err)
		return ret
	}

	if implementationName == "commandexecutorgitoo" {
		ret, err := commandexecutorgitoo.NewGitRepositoryFromPath(path)
		require.NoError(t, err)
		return ret
	}

	t.Fatalf("Unknown implementationName='%s'", implementationName)

	return nil
}

func Test_GetRootDirectoryPath(t *testing.T) {
	var tests = []struct {
		implementationName string
	}{
		{"nativegitoo"},
		{"commandexecutorgitoo"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				_, err = nativegit.InitializeEmptyGitRepository(ctx, tempDir)
				require.NoError(t, err)

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				repoRoot, err := gitRepo.GetRootDirectoryPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tempDir, repoRoot)
			},
		)
	}
}
