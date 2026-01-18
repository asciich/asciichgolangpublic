package nativedocker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
)

func Test_IsRemovalAlreadyInProgressError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, nativedocker.IsRemovalAlreadyInProgressError(nil))
	})

	t.Run("hello world error", func(t *testing.T) {
		require.False(t, nativedocker.IsRemovalAlreadyInProgressError(fmt.Errorf("hello world")))
	})

	t.Run("hello world error", func(t *testing.T) {
		require.True(t, nativedocker.IsRemovalAlreadyInProgressError(fmt.Errorf("Error response from daemon: removal of container test-yay-is-package-update-available is already in progress")))
	})
}
