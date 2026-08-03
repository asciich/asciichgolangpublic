package osutils_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func TestIsExecutableNotFoundError(t *testing.T) {
	t.Run("returns_false_for_nil", func(t *testing.T) {
		result := osutils.IsExecutableNotFoundError(nil)
		require.False(t, result, "IsExecutableNotFoundError should return false for nil")
	})

	t.Run("returns_true_for_ErrExecutableNotFound", func(t *testing.T) {
		result := osutils.IsExecutableNotFoundError(osutils.ErrExecutableNotFound)
		require.True(t, result, "IsExecutableNotFoundError should return true for ErrExecutableNotFound")
	})

	t.Run("returns_false_for_error_with_same_message", func(t *testing.T) {
		wrappedErr := errors.New("executable not found")
		result := osutils.IsExecutableNotFoundError(wrappedErr)
		require.False(t, result, "IsExecutableNotFoundError should return false for error with same message but not the same error")
	})

	t.Run("returns_true_for_error_wrapped_with_join", func(t *testing.T) {
		wrappedErr := errors.Join(osutils.ErrExecutableNotFound, errors.New("additional context"))
		result := osutils.IsExecutableNotFoundError(wrappedErr)
		require.True(t, result, "IsExecutableNotFoundError should return true for error wrapped with Join")
	})

	t.Run("returns_false_for_other_error", func(t *testing.T) {
		otherErr := errors.New("some other error")
		result := osutils.IsExecutableNotFoundError(otherErr)
		require.False(t, result, "IsExecutableNotFoundError should return false for other errors")
	})
}
