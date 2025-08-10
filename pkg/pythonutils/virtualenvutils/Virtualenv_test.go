package virtualenvutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pythonutils/virtualenvutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsVirtualenv(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()
		isVirtualenv, err := virtualenvutils.IsVirtualEnv(ctx, "")
		require.Error(t, err)
		require.False(t, isVirtualenv)
	})

	t.Run("empty dir", func(t *testing.T) {
		ctx := getCtx()
		emptyDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, emptyDir, &filesoptions.DeleteOptions{})

		isVirtualenv, err := virtualenvutils.IsVirtualEnv(ctx, emptyDir)
		require.NoError(t, err)
		require.False(t, isVirtualenv)
	})

	t.Run("empty virtualenv", func(t *testing.T) {
		ctx := getCtx()
		vePath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, vePath, &filesoptions.DeleteOptions{})

		_, err = commandexecutorexec.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"virtualenv", vePath},
			},
		)
		require.NoError(t, err)

		isVirtualenv, err := virtualenvutils.IsVirtualEnv(ctx, vePath)
		require.NoError(t, err)
		require.True(t, isVirtualenv)
	})
}

func Test_CreateVirtualEnv(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		ctx := getCtx()
		ve, err := virtualenvutils.CreateVirtualEnv(ctx, nil)
		require.Error(t, err)
		require.Nil(t, ve)
	})

	t.Run("empty options", func(t *testing.T) {
		// This test is expected to fail as there is no path for the ve specified.
		ctx := getCtx()
		ve, err := virtualenvutils.CreateVirtualEnv(ctx, nil)
		require.Error(t, err)
		require.Nil(t, ve)
	})

	t.Run("happy path no packages", func(t *testing.T) {
		ctx := getCtx()

		vePath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		ve, err := virtualenvutils.CreateVirtualEnv(ctx, &virtualenvutils.CreateVirtualenvOptions{
			Path: vePath,
		})
		require.NoError(t, err)

		// No package was defined so expect pyyaml to be absent:
		pyyamlInstalled, err := ve.IsPackageInstalled(ctx, "pyyaml")
		require.NoError(t, err)
		require.False(t, pyyamlInstalled)
	})

	t.Run("happy path pyyaml package", func(t *testing.T) {
		ctx := getCtx()

		vePath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		// Run twice to validate idempotence
		for range 2 {
			ve, err := virtualenvutils.CreateVirtualEnv(ctx, &virtualenvutils.CreateVirtualenvOptions{
				Path:     vePath,
				Packages: []string{"pyyaml"},
			})
			require.NoError(t, err)

			pyyamlInstalled, err := ve.IsPackageInstalled(ctx, "pyyaml")
			require.NoError(t, err)
			require.True(t, pyyamlInstalled)
		}
	})
}
