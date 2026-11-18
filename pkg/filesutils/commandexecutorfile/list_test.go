package commandexecutorfile

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// fakeCommandExecutor only implements the methods used by ListSubDirectories.
// The embedded (nil) interface satisfies all other methods of the interface,
// calling any of them would panic and show an unexpected dependency.
type fakeCommandExecutor struct {
	commandexecutorinterfaces.CommandExecutor

	hostDescription    string
	hostDescriptionErr error

	stdoutLines []string
	runErr      error

	runCalls    int
	lastOptions *parameteroptions.RunCommandOptions
}

func (f *fakeCommandExecutor) GetHostDescription() (string, error) {
	return f.hostDescription, f.hostDescriptionErr
}

func (f *fakeCommandExecutor) RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]string, error) {
	f.runCalls++
	f.lastOptions = options

	if f.runErr != nil {
		return nil, f.runErr
	}

	// Return a copy so the function under test cannot modify our fixture.
	return slices.Clone(f.stdoutLines), nil
}

func newFakeExecutor(stdout ...string) *fakeCommandExecutor {
	return &fakeCommandExecutor{
		hostDescription: "fake-host",
		stdoutLines:     stdout,
	}
}

func TestListSubDirectories_NilCommandExecutor(t *testing.T) {
	result, err := ListSubDirectories(context.Background(), nil, "/some/path", nil)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestListSubDirectories_EmptyPath(t *testing.T) {
	executor := newFakeExecutor()

	result, err := ListSubDirectories(context.Background(), executor, "", nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, 0, executor.runCalls, "no command must be run for an invalid path")
}

func TestListSubDirectories_GetHostDescriptionFails(t *testing.T) {
	executor := newFakeExecutor()
	executor.hostDescriptionErr = errors.New("no host description")

	result, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.ErrorIs(t, err, executor.hostDescriptionErr)
	require.Nil(t, result)
	require.Equal(t, 0, executor.runCalls)
}

func TestListSubDirectories_RunCommandFails(t *testing.T) {
	executor := newFakeExecutor()
	executor.runErr = errors.New("find failed")

	result, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.ErrorIs(t, err, executor.runErr)
	require.Nil(t, result)
	require.Equal(t, 1, executor.runCalls)
}

func TestListSubDirectories_NilOptionsIsNotRecursive(t *testing.T) {
	executor := newFakeExecutor("/some/path")

	_, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.NoError(t, err)
	require.Equal(t, 1, executor.runCalls)
	require.NotNil(t, executor.lastOptions)
	require.Equal(t,
		[]string{"find", "/some/path", "-type", "d", "-maxdepth", "1"},
		executor.lastOptions.Command,
	)
}

func TestListSubDirectories_NonRecursiveCommand(t *testing.T) {
	executor := newFakeExecutor("/some/path")

	_, err := ListSubDirectories(
		context.Background(),
		executor,
		"/some/path",
		&parameteroptions.ListDirectoryOptions{Recursive: false},
	)

	require.NoError(t, err)
	require.Equal(t,
		[]string{"find", "/some/path", "-type", "d", "-maxdepth", "1"},
		executor.lastOptions.Command,
	)
}

func TestListSubDirectories_RecursiveCommand(t *testing.T) {
	executor := newFakeExecutor("/some/path")

	_, err := ListSubDirectories(
		context.Background(),
		executor,
		"/some/path",
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.NoError(t, err)
	require.Equal(t,
		[]string{"find", "/some/path", "-type", "d"},
		executor.lastOptions.Command,
	)
}

func TestListSubDirectories_PathWithSpacesIsPassedAsSingleArgument(t *testing.T) {
	executor := newFakeExecutor()

	_, err := ListSubDirectories(context.Background(), executor, "/path with spaces", nil)

	require.NoError(t, err)
	require.Contains(t, executor.lastOptions.Command, "/path with spaces")
}

func TestListSubDirectories_NoDirectoriesFound(t *testing.T) {
	executor := newFakeExecutor() // find prints nothing

	result, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.NoError(t, err)
	require.Empty(t, result)
}

// This test describes the INTENDED behavior and currently FAILS because
// ListSubDirectories returns the never-filled `ret` slice instead of `foundPaths`.
func TestListSubDirectories_ReturnsFoundDirectories(t *testing.T) {
	executor := newFakeExecutor(
		"/some/path/c",
		"/some/path/a",
		"/some/path/b",
	)

	result, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.NoError(t, err)
	require.ElementsMatch(t,
		[]string{"/some/path/a", "/some/path/b", "/some/path/c"},
		result,
	)
}

// Also FAILS until the bug is fixed.
func TestListSubDirectories_ResultIsSorted(t *testing.T) {
	executor := newFakeExecutor(
		"/some/path/c",
		"/some/path/a",
		"/some/path/b",
	)

	result, err := ListSubDirectories(context.Background(), executor, "/some/path", nil)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.True(t, slices.IsSorted(result), "result is not sorted: %v", result)
}

// Also FAILS until the bug is fixed.
func TestListSubDirectories_RecursiveReturnsNestedDirectories(t *testing.T) {
	executor := newFakeExecutor(
		"/some/path/a/nested",
		"/some/path/a",
	)

	result, err := ListSubDirectories(
		context.Background(),
		executor,
		"/some/path",
		&parameteroptions.ListDirectoryOptions{Recursive: true},
	)

	require.NoError(t, err)
	require.Contains(t, result, "/some/path/a")
	require.Contains(t, result, "/some/path/a/nested")
}
