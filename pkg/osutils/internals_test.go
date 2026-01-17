package osutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_removeManpagePaths(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		paths := removeManpagePaths(nil)
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("empty", func(t *testing.T) {
		paths := removeManpagePaths([]string{})
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("yay in bin", func(t *testing.T) {
		paths := removeManpagePaths([]string{"/bin/yay"})
		require.EqualValues(t, []string{"/bin/yay"}, paths)
	})

	t.Run("yay in bin and manpage", func(t *testing.T) {
		paths := removeManpagePaths([]string{"/bin/yay", "/usr/share/man/man8/yay.8.gz"})
		require.EqualValues(t, []string{"/bin/yay"}, paths)
	})

	t.Run("yay  manpage", func(t *testing.T) {
		paths := removeManpagePaths([]string{"/usr/share/man/man8/yay.8.gz"})
		require.EqualValues(t, []string{}, paths)
	})
}

func Test_extractBinaryPathsFromStdout(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("")
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("No paths", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("yay:")
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("No paths with line break", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("yay:\n")
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("Only man paths", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("yay: /usr/share/man/man8/yay.8.gz\n")
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("Only man paths no prefix", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("/usr/share/man/man8/yay.8.gz\n")
		require.EqualValues(t, []string{}, paths)
	})

	t.Run("multiple pacman", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("pacman: /usr/bin/pacman /usr/share/pacman /usr/share/man/man8/pacman.8.gz /usr/share/man/man6/pacman.6.gz")
		require.EqualValues(t, []string{"/usr/bin/pacman", "/usr/share/pacman"}, paths)
	})

	t.Run("multiple pacman no man", func(t *testing.T) {
		paths := extractBinaryPathsFromStdout("pacman: /usr/bin/pacman /usr/share/pacman")
		require.EqualValues(t, []string{"/usr/bin/pacman", "/usr/share/pacman"}, paths)
	})
}
