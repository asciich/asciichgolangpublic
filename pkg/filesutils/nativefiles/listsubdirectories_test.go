package nativefiles_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// createTestTree creates the given directories and files below root.
// Entries ending with "/" are created as directories, all others as empty files.
func createTestTree(t *testing.T, root string, entries ...string) {
	t.Helper()

	for _, entry := range entries {
		full := filepath.Join(root, filepath.FromSlash(entry))

		if len(entry) > 0 && entry[len(entry)-1] == '/' {
			require.NoError(t, os.MkdirAll(full, 0o755))
			continue
		}

		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		require.NoError(t, os.WriteFile(full, []byte("content"), 0o644))
	}
}

// joinAll joins every relative path with root.
func joinAll(root string, relPaths ...string) []string {
	result := make([]string, 0, len(relPaths))
	for _, rel := range relPaths {
		result = append(result, filepath.Join(root, filepath.FromSlash(rel)))
	}
	return result
}

func TestListSubDirectories_EmptyPath(t *testing.T) {
	for _, recursive := range []bool{false, true} {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			"",
			&parameteroptions.ListDirectoryOptions{Recursive: recursive},
		)

		require.Error(t, err)
		require.ErrorContains(t, err, "path is empty")
		require.Nil(t, got)
	}
}

func TestListSubDirectories_NonExistingPath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	for _, recursive := range []bool{false, true} {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			missing,
			&parameteroptions.ListDirectoryOptions{Recursive: recursive},
		)

		require.Error(t, err)
		require.Nil(t, got)
	}
}

func TestListSubDirectories_PathIsFile_NonRecursive(t *testing.T) {
	root := t.TempDir()
	createTestTree(t, root, "file.txt")

	got, err := nativefiles.ListSubDirectories(
		context.Background(),
		filepath.Join(root, "file.txt"),
		&parameteroptions.ListDirectoryOptions{},
	)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestListSubDirectories_EmptyDirectory(t *testing.T) {
	for _, recursive := range []bool{false, true} {
		root := t.TempDir()

		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{Recursive: recursive},
		)

		require.NoError(t, err)
		require.NotNil(t, got, "an empty result must be an empty slice, not nil")
		require.Empty(t, got)
	}
}

func TestListSubDirectories(t *testing.T) {
	// Layout used by all cases:
	//
	//   root/
	//   ├── .hidden/
	//   ├── alpha/
	//   │   ├── alpha_file.txt
	//   │   └── nested/
	//   │       └── deep/
	//   ├── beta/
	//   ├── file.txt
	//   └── gamma/
	//       └── inner/
	tree := []string{
		".hidden/",
		"alpha/nested/deep/",
		"alpha/alpha_file.txt",
		"beta/",
		"file.txt",
		"gamma/inner/",
	}

	tests := []struct {
		name     string
		options  *parameteroptions.ListDirectoryOptions
		expected []string // relative to root
	}{
		{
			name:    "nil options behave like non recursive",
			options: nil,
			expected: []string{
				".hidden",
				"alpha",
				"beta",
				"gamma",
			},
		},
		{
			name:    "default options are non recursive",
			options: &parameteroptions.ListDirectoryOptions{},
			expected: []string{
				".hidden",
				"alpha",
				"beta",
				"gamma",
			},
		},
		{
			name:    "explicit non recursive",
			options: &parameteroptions.ListDirectoryOptions{Recursive: false},
			expected: []string{
				".hidden",
				"alpha",
				"beta",
				"gamma",
			},
		},
		{
			name:    "recursive",
			options: &parameteroptions.ListDirectoryOptions{Recursive: true},
			expected: []string{
				".hidden",
				"alpha",
				"alpha/nested",
				"alpha/nested/deep",
				"beta",
				"gamma",
				"gamma/inner",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			createTestTree(t, root, tree...)

			got, err := nativefiles.ListSubDirectories(context.Background(), root, tt.options)

			require.NoError(t, err)
			// Equal (not ElementsMatch) on purpose: the output is expected to be
			// deterministic and lexically ordered.
			require.Equal(t, joinAll(root, tt.expected...), got)
		})
	}
}

func TestListSubDirectories_DoesNotContainRootOrFiles(t *testing.T) {
	for _, recursive := range []bool{false, true} {
		root := t.TempDir()
		createTestTree(t, root, "sub/", "sub/file.txt", "root_file.txt")

		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{Recursive: recursive},
		)

		require.NoError(t, err)
		require.NotContains(t, got, root)
		require.NotContains(t, got, filepath.Join(root, "root_file.txt"))
		require.NotContains(t, got, filepath.Join(root, "sub", "file.txt"))
		require.Equal(t, joinAll(root, "sub"), got)
	}
}

func TestListSubDirectories_ReturnedPathsAreBelowGivenPath(t *testing.T) {
	root := t.TempDir()
	createTestTree(t, root, "a/b/c/")

	got, err := nativefiles.ListSubDirectories(
		context.Background(),
		root,
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.NoError(t, err)
	require.Len(t, got, 3)

	for _, p := range got {
		rel, err := filepath.Rel(root, p)
		require.NoError(t, err)
		require.NotEqual(t, ".", rel)
		require.False(t, filepath.IsAbs(rel))
		require.NotContains(t, rel, "..")

		info, err := os.Stat(p)
		require.NoError(t, err)
		require.True(t, info.IsDir())
	}
}

func TestListSubDirectories_RelativeInputPath(t *testing.T) {
	root := t.TempDir()
	createTestTree(t, root, "one/two/", "three/")

	t.Chdir(root)

	got, err := nativefiles.ListSubDirectories(
		context.Background(),
		".",
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.NoError(t, err)
	require.Equal(t, []string{"one", filepath.Join("one", "two"), "three"}, got)
}

func TestListSubDirectories_SymlinkToDirectoryIsNotListed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	createTestTree(t, root, "real/", "real/inner/")

	require.NoError(t, os.Symlink(
		filepath.Join(root, "real"),
		filepath.Join(root, "link_to_real"),
	))

	t.Run("non recursive", func(t *testing.T) {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{},
		)

		require.NoError(t, err)
		require.Equal(t, joinAll(root, "real"), got)
	})

	t.Run("recursive", func(t *testing.T) {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{Recursive: true},
		)

		require.NoError(t, err)
		// The symlink is neither listed nor followed.
		require.Equal(t, joinAll(root, "real", "real/inner"), got)
	})
}

func TestListSubDirectories_Recursive_CancelledContext(t *testing.T) {
	root := t.TempDir()
	createTestTree(t, root, "a/b/", "c/")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := nativefiles.ListSubDirectories(
		ctx,
		root,
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.Error(t, err)
	require.ErrorContains(t, err, "context canceled")
	require.Nil(t, got)
}

func TestListSubDirectories_Recursive_ContextCancelledWhileWalking(t *testing.T) {
	// Not cancelled: the same tree must be listed successfully.
	root := t.TempDir()
	createTestTree(t, root, "a/b/", "c/")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got, err := nativefiles.ListSubDirectories(
		ctx,
		root,
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.NoError(t, err)
	require.Equal(t, joinAll(root, "a", "a/b", "c"), got)
}

func TestListSubDirectories_UnreadableSubDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permissions are not supported on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}

	root := t.TempDir()
	createTestTree(t, root, "readable/", "locked/inner/")

	locked := filepath.Join(root, "locked")
	require.NoError(t, os.Chmod(locked, 0o000))
	// Restore permissions so t.TempDir() can clean up.
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	t.Run("non recursive still lists the locked directory itself", func(t *testing.T) {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{},
		)

		require.NoError(t, err)
		require.Equal(t, joinAll(root, "locked", "readable"), got)
	})

	t.Run("recursive fails because the locked directory cannot be read", func(t *testing.T) {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{Recursive: true},
		)

		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestListSubDirectories_UnreadableRootDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permissions are not supported on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}

	root := t.TempDir()
	createTestTree(t, root, "sub/")

	require.NoError(t, os.Chmod(root, 0o000))
	t.Cleanup(func() { _ = os.Chmod(root, 0o755) })

	for _, recursive := range []bool{false, true} {
		got, err := nativefiles.ListSubDirectories(
			context.Background(),
			root,
			&parameteroptions.ListDirectoryOptions{Recursive: recursive},
		)

		require.Error(t, err)
		require.Nil(t, got)
	}
}
