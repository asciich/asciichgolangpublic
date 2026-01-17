package commandexecutorgitoo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateTemporaryRepository(t *testing.T) {
	tests := []struct {
		initialize bool
	}{
		{true},
		{false},
	}
	ctx := getCtx()

	for _, tt := range tests {
		repo, err := commandexecutorgitoo.CreateLocalTemporaryRepository(ctx, &parameteroptions.CreateRepositoryOptions{
			InitializeWithEmptyCommit:   tt.initialize,
			InitializeWithDefaultAuthor: tt.initialize,
		})
		require.NoError(t, err)
		defer func() {
			err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
			require.NoError(t, err)
		}()

		exists, err := repo.Exists(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		isIntialized, err := repo.IsInitialized(ctx)
		require.NoError(t, err)
		require.EqualValues(t, tt.initialize, isIntialized)
	}
}

func Test_CloneToTemporaryDirectory(t *testing.T) {
	ctx := getCtx()

	repo, err := commandexecutorgitoo.CreateLocalTemporaryRepository(ctx, &parameteroptions.CreateRepositoryOptions{
		InitializeWithEmptyCommit:   true,
		InitializeWithDefaultAuthor: true,
	})
	require.NoError(t, err)
	defer func() {
		err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
		require.NoError(t, err)
	}()

	isInitalized, err := repo.IsInitialized(ctx)
	require.NoError(t, err)
	require.True(t, isInitalized)

	cloned, err := repo.CloneToTemporaryRepository(ctx)
	require.NoError(t, err)
	defer cloned.Delete(ctx, &filesoptions.DeleteOptions{})
	require.NotNil(t, cloned)

	remoteConfigs, err := cloned.GetRemoteConfigs(ctx)
	require.NoError(t, err)

	require.Len(t, remoteConfigs, 1)
	remoteName, err := remoteConfigs[0].GetRemoteName()
	require.NoError(t, err)
	require.EqualValues(t, "origin", remoteName)

	repoPath, err := repo.GetPath()
	require.NoError(t, err)
	remoteUrl, err := remoteConfigs[0].GetUrlFetch()
	require.NoError(t, err)
	require.EqualValues(t, repoPath, remoteUrl)
}
