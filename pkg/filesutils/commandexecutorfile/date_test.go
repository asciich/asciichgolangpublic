package commandexecutorfile_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

// getCommandExecutorsToTest returns all command executors the functions must
// behave identically on (local execution variants).
func getCommandExecutorsToTest() []struct {
	Name     string
	Executor commandexecutorinterfaces.CommandExecutor
} {
	return []struct {
		Name     string
		Executor commandexecutorinterfaces.CommandExecutor
	}{
		{Name: "Exec", Executor: commandexecutorexecoo.Exec()},
		{Name: "Bash", Executor: commandexecutorbashoo.Bash()},
	}
}

// createTestFile creates a fresh file on the host targeted by the executor and
// returns its path. It uses t.TempDir() so cleanup is automatic.
func createTestFile(t *testing.T, ctx context.Context, executor commandexecutorinterfaces.CommandExecutor) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "birthdate_testfile.txt")

	err := commandexecutorfile.AppendString(ctx, executor, path, "birthdate test content\n")
	require.NoError(t, err)

	return path
}

func TestGetBirthDate(t *testing.T) {
	ctx := context.Background()

	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		_, err := commandexecutorfile.GetBirthDate(ctx, nil, "/tmp/whatever")
		require.Error(t, err)
	})

	for _, impl := range getCommandExecutorsToTest() {
		impl := impl

		t.Run(impl.Name+"_empty path returns error", func(t *testing.T) {
			_, err := commandexecutorfile.GetBirthDate(ctx, impl.Executor, "")
			require.Error(t, err)
		})

		t.Run(impl.Name+"_non existing path returns error", func(t *testing.T) {
			nonExisting := filepath.Join(t.TempDir(), "does_not_exist.txt")

			_, err := commandexecutorfile.GetBirthDate(ctx, impl.Executor, nonExisting)
			require.Error(t, err)
		})

		t.Run(impl.Name+"_freshly created file has a recent birth date", func(t *testing.T) {
			before := time.Now().Add(-5 * time.Second)
			path := createTestFile(t, ctx, impl.Executor)
			after := time.Now().Add(5 * time.Second)

			birthDate, err := commandexecutorfile.GetBirthDate(ctx, impl.Executor, path)
			if err != nil {
				// Some filesystems / kernels do not expose a birth time
				// ('stat -c %W' returns 0). In that case the function correctly
				// returns an error and there is nothing to assert about a value.
				t.Skipf("Birth date not available on this filesystem: %v", err)
			}

			require.False(t, birthDate.IsZero())
			require.WithinRange(t, birthDate, before, after)
		})

		t.Run(impl.Name+"_birth date is returned in UTC", func(t *testing.T) {
			path := createTestFile(t, ctx, impl.Executor)

			birthDate, err := commandexecutorfile.GetBirthDate(ctx, impl.Executor, path)
			if err != nil {
				t.Skipf("Birth date not available on this filesystem: %v", err)
			}

			require.Equal(t, time.UTC, birthDate.Location())
		})
	}
}

func TestGetBirthDateRFC3339(t *testing.T) {
	ctx := context.Background()

	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		_, err := commandexecutorfile.GetBirthDateRFC3339(ctx, nil, "/tmp/whatever")
		require.Error(t, err)
	})

	for _, impl := range getCommandExecutorsToTest() {
		impl := impl

		t.Run(impl.Name+"_empty path returns error", func(t *testing.T) {
			_, err := commandexecutorfile.GetBirthDateRFC3339(ctx, impl.Executor, "")
			require.Error(t, err)
		})

		t.Run(impl.Name+"_non existing path returns error", func(t *testing.T) {
			nonExisting := filepath.Join(t.TempDir(), "does_not_exist.txt")

			_, err := commandexecutorfile.GetBirthDateRFC3339(ctx, impl.Executor, nonExisting)
			require.Error(t, err)
		})

		t.Run(impl.Name+"_returns a parseable RFC3339 string", func(t *testing.T) {
			path := createTestFile(t, ctx, impl.Executor)

			rfc3339, err := commandexecutorfile.GetBirthDateRFC3339(ctx, impl.Executor, path)
			if err != nil {
				t.Skipf("Birth date not available on this filesystem: %v", err)
			}

			parsed, parseErr := time.Parse(time.RFC3339, rfc3339)
			require.NoError(t, parseErr)
			require.False(t, parsed.IsZero())
		})

		t.Run(impl.Name+"_RFC3339 string matches GetBirthDate value", func(t *testing.T) {
			path := createTestFile(t, ctx, impl.Executor)

			birthDate, err := commandexecutorfile.GetBirthDate(ctx, impl.Executor, path)
			if err != nil {
				t.Skipf("Birth date not available on this filesystem: %v", err)
			}

			rfc3339, err := commandexecutorfile.GetBirthDateRFC3339(ctx, impl.Executor, path)
			require.NoError(t, err)

			require.Equal(t, birthDate.Format(time.RFC3339), rfc3339)
		})
	}
}

// TestGetBirthDateRoundedToSecondsAcrossImplementations verifies the two
// command executors report the same birth time for the same file.
func TestGetBirthDateSameAcrossImplementations(t *testing.T) {
	ctx := context.Background()

	// A single file on the local host, shared by all executors.
	path := filepath.Join(t.TempDir(), "shared_birthdate_testfile.txt")

	execExecutor := commandexecutorexecoo.Exec()

	err := commandexecutorfile.AppendString(ctx, execExecutor, path, "shared content\n")
	require.NoError(t, err)

	execBirthDate, err := commandexecutorfile.GetBirthDate(ctx, execExecutor, path)
	if err != nil {
		t.Skipf("Birth date not available on this filesystem: %v", err)
	}

	bashBirthDate, err := commandexecutorfile.GetBirthDate(ctx, commandexecutorbashoo.Bash(), path)
	require.NoError(t, err)

	// 'stat -c %W' has second precision, so both must be identical.
	require.Equal(t, execBirthDate, bashBirthDate)
}
