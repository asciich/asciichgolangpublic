package nativefiles_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

// createTestFile creates a fresh file inside t.TempDir() and returns its path.
// Cleanup is handled via defer + t.TempDir (auto-removed by the test framework).
func createTestFile(t *testing.T, ctx context.Context) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "birthdate_testfile.txt")

	err := nativefiles.Create(ctx, path, &filesoptions.CreateOptions{})
	require.NoError(t, err)

	return path
}

func TestGetBirthDate(t *testing.T) {
	ctx := context.Background()

	t.Run("empty path returns error", func(t *testing.T) {
		_, err := nativefiles.GetBirthDate(ctx, "")
		require.Error(t, err)
	})

	t.Run("non existing path returns error", func(t *testing.T) {
		nonExisting := filepath.Join(t.TempDir(), "does_not_exist.txt")

		_, err := nativefiles.GetBirthDate(ctx, nonExisting)
		require.Error(t, err)
	})

	t.Run("freshly created file has a recent birth date", func(t *testing.T) {
		before := time.Now().Add(-5 * time.Second)
		path := createTestFile(t, ctx)
		after := time.Now().Add(5 * time.Second)

		birthDate, err := nativefiles.GetBirthDate(ctx, path)
		if err != nil {
			// Some filesystems / kernels do not expose a birth time. In that
			// case the function correctly returns an error and there is
			// nothing meaningful to assert about the value.
			t.Skipf("Birth date not available on this filesystem: %v", err)
		}

		require.False(t, birthDate.IsZero())
		require.WithinRange(t, birthDate, before, after)
	})

	t.Run("birth date is returned in UTC", func(t *testing.T) {
		path := createTestFile(t, ctx)

		birthDate, err := nativefiles.GetBirthDate(ctx, path)
		if err != nil {
			t.Skipf("Birth date not available on this filesystem: %v", err)
		}

		require.Equal(t, time.UTC, birthDate.Location())
	})

	t.Run("birth date of a directory can be read", func(t *testing.T) {
		dir := t.TempDir()

		birthDate, err := nativefiles.GetBirthDate(ctx, dir)
		if err != nil {
			t.Skipf("Birth date not available on this filesystem: %v", err)
		}

		require.False(t, birthDate.IsZero())
	})
}

func TestGetBirthDateRFC3339(t *testing.T) {
	ctx := context.Background()

	t.Run("empty path returns error", func(t *testing.T) {
		_, err := nativefiles.GetBirthDateRFC3339(ctx, "")
		require.Error(t, err)
	})

	t.Run("non existing path returns error", func(t *testing.T) {
		nonExisting := filepath.Join(t.TempDir(), "does_not_exist.txt")

		_, err := nativefiles.GetBirthDateRFC3339(ctx, nonExisting)
		require.Error(t, err)
	})

	t.Run("returns a parseable RFC3339 string", func(t *testing.T) {
		path := createTestFile(t, ctx)

		rfc3339, err := nativefiles.GetBirthDateRFC3339(ctx, path)
		if err != nil {
			t.Skipf("Birth date not available on this filesystem: %v", err)
		}

		parsed, parseErr := time.Parse(time.RFC3339, rfc3339)
		require.NoError(t, parseErr)
		require.False(t, parsed.IsZero())
	})

	t.Run("RFC3339 string matches GetBirthDate value", func(t *testing.T) {
		path := createTestFile(t, ctx)

		birthDate, err := nativefiles.GetBirthDate(ctx, path)
		if err != nil {
			t.Skipf("Birth date not available on this filesystem: %v", err)
		}

		rfc3339, err := nativefiles.GetBirthDateRFC3339(ctx, path)
		require.NoError(t, err)

		// The formatted string must represent the exact same instant as the
		// time.Time returned by GetBirthDate.
		require.Equal(t, birthDate.Format(time.RFC3339), rfc3339)
	})
}
