package installkvm_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/systemdutils"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/installkvm"
)

func Test_RequiresRoot(t *testing.T) {
	ctx := context.Background()

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	require.NoError(t, err)

	if isRoot {
		t.Skip("Test validates behavior without root privileges but is running as root.")
	}

	t.Run("InstallPackages fails without root", func(t *testing.T) {
		err := installkvm.InstallPackages(ctx)
		require.Error(t, err)
	})

	t.Run("Start fails without root", func(t *testing.T) {
		err := installkvm.Start(ctx)
		require.Error(t, err)
	})

	t.Run("Enable fails without root", func(t *testing.T) {
		err := installkvm.Enable(ctx)
		require.Error(t, err)
	})

}

func Test_Start(t *testing.T) {
	ctx := context.Background()

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	require.NoError(t, err)

	if !isRoot {
		t.Skip("Test requires root privileges.")
	}

	t.Run("libvirtd is running after start", func(t *testing.T) {
		err := installkvm.Start(ctx)
		require.NoError(t, err)

		isActive, err := systemdutils.IsServiceActive(ctx, commandexecutorexecoo.Exec(), "libvirtd")
		require.NoError(t, err)
		require.True(t, isActive)
	})

	t.Run("start is idempotent", func(t *testing.T) {
		for i := range 2 {
			err := installkvm.Start(ctx)
			require.NoError(t, err, "Start run '%d' failed.", i)
		}

		isActive, err := systemdutils.IsServiceActive(ctx, commandexecutorexecoo.Exec(), "libvirtd")
		require.NoError(t, err)
		require.True(t, isActive)
	})
}

func Test_Enable(t *testing.T) {
	ctx := context.Background()

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	require.NoError(t, err)

	if !isRoot {
		t.Skip("Test requires root privileges.")
	}

	t.Run("enable is idempotent", func(t *testing.T) {
		for i := range 2 {
			err := installkvm.Enable(ctx)
			require.NoError(t, err, "Enable run '%d' failed.", i)
		}
	})
}
