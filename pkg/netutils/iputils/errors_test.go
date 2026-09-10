package iputils_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
)

func TestIsInvalidIPError(t *testing.T) {
	t.Run("nil error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPError(nil))
	})

	t.Run("unrelated error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPError(errors.New("some other error")))
	})

	t.Run("ErrInvalidIP must return true", func(t *testing.T) {
		require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIP))
	})

	t.Run("ErrInvalidIPv4 must return true", func(t *testing.T) {
		require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIPv4))
	})

	t.Run("ErrInvalidIPv6 must return true", func(t *testing.T) {
		require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIPv6))
	})

	t.Run("wrapped ErrInvalidIP must return true", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", iputils.ErrInvalidIP)
		require.True(t, iputils.IsInvalidIPError(wrapped))
	})

	t.Run("wrapped ErrInvalidIPv4 must return true", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", iputils.ErrInvalidIPv4)
		require.True(t, iputils.IsInvalidIPError(wrapped))
	})

	t.Run("wrapped ErrInvalidIPv6 must return true", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", iputils.ErrInvalidIPv6)
		require.True(t, iputils.IsInvalidIPError(wrapped))
	})
}

func TestIsInvalidIPv4Error(t *testing.T) {
	t.Run("nil error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv4Error(nil))
	})

	t.Run("unrelated error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv4Error(errors.New("some other error")))
	})

	t.Run("ErrInvalidIP must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIP))
	})

	t.Run("ErrInvalidIPv6 must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIPv6))
	})

	t.Run("ErrInvalidIPv4 must return true", func(t *testing.T) {
		require.True(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIPv4))
	})

	t.Run("wrapped ErrInvalidIPv4 must return true", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", iputils.ErrInvalidIPv4)
		require.True(t, iputils.IsInvalidIPv4Error(wrapped))
	})
}

func TestIsInvalidIPv6Error(t *testing.T) {
	t.Run("nil error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv6Error(nil))
	})

	t.Run("unrelated error must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv6Error(errors.New("some other error")))
	})

	t.Run("ErrInvalidIP must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIP))
	})

	t.Run("ErrInvalidIPv4 must return false", func(t *testing.T) {
		require.False(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIPv4))
	})

	t.Run("ErrInvalidIPv6 must return true", func(t *testing.T) {
		require.True(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIPv6))
	})

	t.Run("wrapped ErrInvalidIPv6 must return true", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", iputils.ErrInvalidIPv6)
		require.True(t, iputils.IsInvalidIPv6Error(wrapped))
	})
}
