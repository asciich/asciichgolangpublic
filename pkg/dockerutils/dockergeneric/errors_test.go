package dockergeneric_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Test_IsContainerNotFoundError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, dockergeneric.IsErrorContainerNotFound(nil))
	})

	t.Run("another error", func(t *testing.T) {
		require.False(t, dockergeneric.IsErrorContainerNotFound(errors.New("another error")))
	})

	t.Run("container not found error", func(t *testing.T) {
		require.True(t, dockergeneric.IsErrorContainerNotFound(dockergeneric.ErrDockerContainerNotFound))
	})

	t.Run("TracedError wrapping the ErrDockerContainerNotFound", func(t *testing.T) {
		err := tracederrors.TracedErrorf("wrapping: %w", dockergeneric.ErrDockerContainerNotFound)
		require.True(t, dockergeneric.IsErrorContainerNotFound(err))
	})
}
