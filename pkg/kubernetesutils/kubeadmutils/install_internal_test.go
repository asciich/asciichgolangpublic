package kubeadmutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_prepareInstallOptions(t *testing.T) {
	t.Run("nil options should use defaults", func(t *testing.T) {
		opts, err := prepareInstallOptions(nil)
		require.NoError(t, err)
		require.NotNil(t, opts)

		// Verify default install path is for kubeadm, not kubelet or kubectl
		require.Equal(t, "/bin/kubeadm", opts.InstallPath)

		// Verify it doesn't accidentally use kubelet or kubectl paths
		require.NotEqual(t, "/bin/kubelet", opts.InstallPath)
		require.NotEqual(t, "/bin/kubectl", opts.InstallPath)

		// Verify UseSudo defaults to true
		require.True(t, opts.UseSudo)
	})

	t.Run("empty options should use defaults", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{})
		require.NoError(t, err)

		// Verify default install path is for kubeadm
		require.Equal(t, "/bin/kubeadm", opts.InstallPath)
	})

	t.Run("custom install path should be respected", func(t *testing.T) {
		customPath := "/usr/local/bin/kubeadm"
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			InstallPath: customPath,
		})
		require.NoError(t, err)

		require.Equal(t, customPath, opts.InstallPath)
	})

	t.Run("UseSudo false should be respected", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			UseSudo: false,
		})
		require.NoError(t, err)

		require.False(t, opts.UseSudo)
	})

	t.Run("source URL should point to kubeadm binary", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			Version: "v1.36.2",
		})
		require.NoError(t, err)

		// Verify the URL points to kubeadm, not kubelet or kubectl
		require.True(t, strings.HasSuffix(opts.SrcUrl, "/kubeadm"))
		require.False(t, strings.HasSuffix(opts.SrcUrl, "/kubelet"))
		require.False(t, strings.HasSuffix(opts.SrcUrl, "/kubectl"))

		require.Equal(t, "https://dl.k8s.io/release/v1.36.2/bin/linux/amd64/kubeadm", opts.SrcUrl)
	})

	t.Run("unsupported version should return error", func(t *testing.T) {
		_, err := prepareInstallOptions(&InstallKubeadmOptions{
			Version: "v99.99.99",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported kubeadm version")
	})

	t.Run("install options should have correct mode", func(t *testing.T) {
		opts, err := prepareInstallOptions(nil)
		require.NoError(t, err)

		require.Equal(t, "u=rwx,g=rx,o=rx", opts.Mode)
	})

	t.Run("install options should have ReplaceExisting true", func(t *testing.T) {
		opts, err := prepareInstallOptions(nil)
		require.NoError(t, err)

		require.True(t, opts.ReplaceExisting)
	})

	t.Run("SHA256 checksum should be set for supported version", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			Version: "v1.36.2",
		})
		require.NoError(t, err)

		require.NotEmpty(t, opts.Sha256Sum)

		// Verify it's a valid SHA256 hash (64 hex characters)
		require.Len(t, opts.Sha256Sum, 64)
	})
}

// Test_prepareInstallOptions_PackageIsolation ensures kubeadmutils doesn't
// accidentally use paths or URLs from kubeletutils or kubectlutils.
func Test_prepareInstallOptions_PackageIsolation(t *testing.T) {
	t.Run("should not use kubelet install path", func(t *testing.T) {
		opts, err := prepareInstallOptions(nil)
		require.NoError(t, err)

		require.NotEqual(t, "/bin/kubelet", opts.InstallPath,
			"CRITICAL: kubeadmutils is using kubelet install path - packages are mixed up!")
	})

	t.Run("should not use kubectl install path", func(t *testing.T) {
		opts, err := prepareInstallOptions(nil)
		require.NoError(t, err)

		require.NotEqual(t, "/bin/kubectl", opts.InstallPath,
			"CRITICAL: kubeadmutils is using kubectl install path - packages are mixed up!")
	})

	t.Run("should not use kubelet download URL", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			Version: "v1.36.2",
		})
		require.NoError(t, err)

		require.NotEqual(t, "https://dl.k8s.io/release/v1.36.2/bin/linux/amd64/kubelet", opts.SrcUrl,
			"CRITICAL: kubeadmutils is using kubelet download URL - packages are mixed up!")
	})

	t.Run("should not use kubectl download URL", func(t *testing.T) {
		opts, err := prepareInstallOptions(&InstallKubeadmOptions{
			Version: "v1.36.2",
		})
		require.NoError(t, err)

		require.NotEqual(t, "https://dl.k8s.io/release/v1.36.2/bin/linux/amd64/kubectl", opts.SrcUrl,
			"CRITICAL: kubeadmutils is using kubectl download URL - packages are mixed up!")
	})
}
