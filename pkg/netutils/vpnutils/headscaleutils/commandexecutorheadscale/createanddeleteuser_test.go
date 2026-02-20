package commandexecutorheadscale_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_headscale_CreateAndDeleteUsers(t *testing.T) {
	// Currently not available in Github CI:
	continuousintegration.SkipInGithubCi(t, "Expose docker port does not work yet.")

	ctx := getCtx()

	const containerName = "test-headscale"
	nativedocker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})

	// Use a minimal config:
	configPath, err := headscalegeneric.WriteMinimalConfigAsTemporaryFile(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, nativefiles.Delete(ctx, configPath, &filesoptions.DeleteOptions{}))
	}()

	headscaleContainer, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName:        "headscale/headscale:latest",
		Name:             containerName,
		Command:          []string{"serve"},
		Mounts:           []string{configPath + ":/etc/headscale/config.yaml"},
		WaitForPortsOpen: true,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, headscaleContainer.Remove(ctx, &dockeroptions.RemoveOptions{Force: true}))
	}()

	// Create the testuser
	const username = "testuser"
	ctxCreate := contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.CreateUser(ctxCreate, headscaleContainer, username)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate))

	// Create the testuser again to check idempotence:
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.CreateUser(ctxCreate, headscaleContainer, username)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxCreate))

	// We only expect one user:
	users, err := commandexecutorheadscale.ListUserNames(ctx, headscaleContainer)
	require.NoError(t, err)
	require.EqualValues(t, []string{"testuser"}, users)

	// Create another testuser
	const username2 = "another_testuser"
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.CreateUser(ctxCreate, headscaleContainer, username2)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxCreate))

	// Create the testuser again to check idempotence:
	ctxCreate = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.CreateUser(ctxCreate, headscaleContainer, username2)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxCreate))

	// We only expect two users:
	users, err = commandexecutorheadscale.ListUserNames(ctx, headscaleContainer)
	require.NoError(t, err)
	require.EqualValues(t, []string{"another_testuser", "testuser"}, users)

	// Delete the first testuser:
	ctxDelete := contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.DeleteUser(ctxDelete, headscaleContainer, username)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	// Delete the first testuser again to check idempotence:
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.DeleteUser(ctxDelete, headscaleContainer, username)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDelete))

	// List remaining user:
	users, err = commandexecutorheadscale.ListUserNames(ctx, headscaleContainer)
	require.NoError(t, err)
	require.EqualValues(t, []string{"another_testuser"}, users)

	exists, err := commandexecutorheadscale.UserExists(ctx, headscaleContainer, username2)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the second testuser:
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.DeleteUser(ctxDelete, headscaleContainer, username2)
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	// Delete the second testuser again to check idempotence:
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = commandexecutorheadscale.DeleteUser(ctxDelete, headscaleContainer, username2)
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDelete))

	// No users to list left:
	users, err = commandexecutorheadscale.ListUserNames(ctx, headscaleContainer)
	require.NoError(t, err)
	require.EqualValues(t, []string{}, users)

	exists, err = commandexecutorheadscale.UserExists(ctx, headscaleContainer, username)
	require.NoError(t, err)
	require.False(t, exists)

	exists, err = commandexecutorheadscale.UserExists(ctx, headscaleContainer, username2)
	require.NoError(t, err)
	require.False(t, exists)
}
