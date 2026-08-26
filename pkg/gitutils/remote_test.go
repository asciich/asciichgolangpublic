package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_SetRemoteUrl(t *testing.T) {
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

				t.Run("overwrite existing remote url", func(t *testing.T) {
					originDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

					_, err = nativegit.InitializeEmptyGitRepository(ctx, originDir)
					require.NoError(t, err)

					cloneDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

					err = gitRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
					require.NoError(t, err)

					err = gitRepo.SetRemoteUrl(ctx, "https://example.com/new-repo.git")
					require.NoError(t, err)

					remoteConfigs, err := gitRepo.GetRemoteConfigs(ctx)
					require.NoError(t, err)
					require.Len(t, remoteConfigs, 1)
				})
			},
		)
	}
}
