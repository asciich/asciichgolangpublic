package filesgeneric_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
)

// fileNotFoundMessages are error messages which must be detected as file not found.
var fileNotFoundMessages = []string{
	"Command failed: 'stat -c %s -- /does/not/exist', exit status 1\nstat: cannot statx '/does/not/exist': No such file or directory\n\nstat: cannot statx '/does/not/exist': No such file or directory\n",
	"stat: cannot statx '/does/not/exist': No such file or directory\n",
	"stat: cannot statx '/does/not/exist': No such file or directory",
	// BusyBox (e.g. alpine) variant:
	"stat: can't stat '/does/not/exist': No such file or directory",
}

// otherErrorMessages are error messages which must NOT match any known error.
var otherErrorMessages = []string{
	"another error",
	"permission denied",
	"stat: cannot statx '/root/secret': Permission denied",
	"",
}

func Test_IsErrFileNotFound(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, filesgeneric.IsErrFileNotFound(nil))
	})

	t.Run("another error", func(t *testing.T) {
		require.False(t, filesgeneric.IsErrFileNotFound(fmt.Errorf("another error")))
	})

	t.Run("direct", func(t *testing.T) {
		require.True(t, filesgeneric.IsErrFileNotFound(filesgeneric.ErrFileNotFound))
	})

	t.Run("wrapped", func(t *testing.T) {
		require.True(t, filesgeneric.IsErrFileNotFound(fmt.Errorf("This is wrapping: %w", filesgeneric.ErrFileNotFound)))
	})

	t.Run("The default os.ErrNotExist returned by os.Open is ErrFileNotFound as well", func(t *testing.T) {
		_, err := os.Open("/this/file/does/not/exist")
		require.True(t, filesgeneric.IsErrFileNotFound(err))
	})
}

func Test_GetAsFileNotFoundErrorIfMessageMatches(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(nil))
	})

	t.Run("other errors are returned untouched", func(t *testing.T) {
		for _, msg := range otherErrorMessages {
			t.Run(fmt.Sprintf("%q", msg), func(t *testing.T) {
				err := errors.New(msg)

				got := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(err)
				require.Same(t, err, got)
				require.False(t, filesgeneric.IsErrFileNotFound(got))
			})
		}
	})

	t.Run("ErrFileNotFound is returned untouched", func(t *testing.T) {
		got := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(filesgeneric.ErrFileNotFound)
		require.Same(t, filesgeneric.ErrFileNotFound, got)
		require.True(t, filesgeneric.IsErrFileNotFound(got))
	})

	t.Run("os.ErrNotExist is returned untouched", func(t *testing.T) {
		_, err := os.Open("/this/file/does/not/exist")
		require.Error(t, err)

		got := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(err)
		require.Same(t, err, got)
		require.True(t, filesgeneric.IsErrFileNotFound(got))
	})

	t.Run("file not found messages are wrapped as ErrFileNotFound", func(t *testing.T) {
		for _, msg := range fileNotFoundMessages {
			t.Run(fmt.Sprintf("%q", msg), func(t *testing.T) {
				err := errors.New(msg)
				require.False(t, filesgeneric.IsErrFileNotFound(err))

				got := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(err)
				require.Error(t, got)
				require.True(t, filesgeneric.IsErrFileNotFound(got))

				// Original error must be preserved:
				require.True(t, errors.Is(got, err))
				require.Contains(t, got.Error(), msg)
			})
		}
	})

	t.Run("wrapping is idempotent", func(t *testing.T) {
		err := errors.New("stat: cannot statx '/does/not/exist': No such file or directory")

		first := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(err)
		second := filesgeneric.GetAsFileNotFoundErrorIfMessageMatches(first)
		require.Same(t, first, second)
	})
}

func Test_GetAsError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, filesgeneric.GetAsError(nil))
	})

	t.Run("other errors are returned untouched", func(t *testing.T) {
		for _, msg := range otherErrorMessages {
			t.Run(fmt.Sprintf("%q", msg), func(t *testing.T) {
				err := errors.New(msg)

				got := filesgeneric.GetAsError(err)
				require.Same(t, err, got)
				require.False(t, filesgeneric.IsErrFileNotFound(got))
			})
		}
	})

	t.Run("ErrFileNotFound is returned untouched", func(t *testing.T) {
		got := filesgeneric.GetAsError(filesgeneric.ErrFileNotFound)
		require.Same(t, filesgeneric.ErrFileNotFound, got)
		require.True(t, filesgeneric.IsErrFileNotFound(got))
	})

	t.Run("os.ErrNotExist is returned untouched", func(t *testing.T) {
		_, err := os.Open("/this/file/does/not/exist")
		require.Error(t, err)

		got := filesgeneric.GetAsError(err)
		require.Same(t, err, got)
		require.True(t, filesgeneric.IsErrFileNotFound(got))
	})

	t.Run("file not found messages are wrapped as ErrFileNotFound", func(t *testing.T) {
		for _, msg := range fileNotFoundMessages {
			t.Run(fmt.Sprintf("%q", msg), func(t *testing.T) {
				err := errors.New(msg)
				require.False(t, filesgeneric.IsErrFileNotFound(err))

				got := filesgeneric.GetAsError(err)
				require.Error(t, got)
				require.True(t, filesgeneric.IsErrFileNotFound(got))

				// Original error must be preserved:
				require.True(t, errors.Is(got, err))
				require.Contains(t, got.Error(), msg)
			})
		}
	})

	t.Run("wrapped original error is detected", func(t *testing.T) {
		inner := errors.New("stat: cannot statx '/does/not/exist': No such file or directory")
		err := fmt.Errorf("outer context: %w", inner)

		got := filesgeneric.GetAsError(err)
		require.True(t, filesgeneric.IsErrFileNotFound(got))
		require.True(t, errors.Is(got, err))
		require.True(t, errors.Is(got, inner))
	})

	t.Run("is idempotent", func(t *testing.T) {
		err := errors.New("stat: cannot statx '/does/not/exist': No such file or directory")

		first := filesgeneric.GetAsError(err)
		second := filesgeneric.GetAsError(first)
		require.Same(t, first, second)
	})
}
