package hostreconciler_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/hostreconciler"
	"github.com/asciich/asciichgolangpublic/pkg/hostreconciler/hostreconcileroptions"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

// kvmTestSource is a minimal Source used to test the KVM reconciliation.
type kvmTestSource struct {
	beforeScript []string
	kvmConfig    *hostreconciler.KVMConfig
}

func (s *kvmTestSource) GetHostName(ctx context.Context) (string, error) {
	return "kvmtesthost", nil
}

func (s *kvmTestSource) GetNetworkInterfaces(ctx context.Context) ([]hostreconciler.NetworkInterface, error) {
	return []hostreconciler.NetworkInterface{}, nil
}

func (s *kvmTestSource) GetBeforeScript(ctx context.Context) ([]string, error) {
	return s.beforeScript, nil
}

func (s *kvmTestSource) GetDockerComposeProjects(ctx context.Context) ([]hostreconciler.DockerComposeProject, error) {
	return []hostreconciler.DockerComposeProject{}, nil
}

func (s *kvmTestSource) GetKVMConfig(ctx context.Context) (*hostreconciler.KVMConfig, error) {
	return s.kvmConfig, nil
}

func Test_ReconcileOnce_KVM(t *testing.T) {
	ctx := context.Background()

	t.Run("nil KVM config is rejected", func(t *testing.T) {
		err := hostreconciler.ReconcileOnce(ctx, &kvmTestSource{kvmConfig: nil}, &hostreconcileroptions.RunOptions{})
		require.Error(t, err)
	})

	t.Run("host is left untouched without root if KVM actions are configured", func(t *testing.T) {
		isRoot, err := userutils.IsRunningAsRoot(ctx)
		require.NoError(t, err)

		if isRoot {
			t.Skip("Test validates behavior without root privileges but is running as root.")
		}

		markerPath := filepath.Join(t.TempDir(), "before_script_executed")

		source := &kvmTestSource{
			beforeScript: []string{"touch '" + markerPath + "'"},
			kvmConfig:    &hostreconciler.KVMConfig{Start: true, Enable: true},
		}

		err = hostreconciler.ReconcileOnce(ctx, source, &hostreconcileroptions.RunOptions{})
		require.Error(t, err)

		exists, err := nativefiles.Exists(ctx, markerPath)
		require.NoError(t, err)
		require.False(t, exists, "before_script must not run if root privileges are missing.")
	})
}
