package nativekubernetesoo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func TestIsConfigMapContentEqual(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapContentEqual(nil, nil))
	})

	t.Run("empty and nil", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapContentEqual(map[string]string{}, nil))
	})

	t.Run("nil and empty", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapContentEqual(nil, map[string]string{}))
	})

	t.Run("empty and empty", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapContentEqual(map[string]string{}, map[string]string{}))
	})

	t.Run("empty and oneEntry", func(t *testing.T) {
		require.False(t, nativekubernetesoo.IsConfigMapContentEqual(map[string]string{}, map[string]string{"a": "bc"}))
	})

	t.Run("oneEntry and oneEntry", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapContentEqual(map[string]string{"a": "bc"}, map[string]string{"a": "bc"}))
	})
}

func Test_IsConfigMapLabelsEqual(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapLabelsEqual(nil, nil))
	})

	t.Run("empty and nil", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapLabelsEqual(map[string]string{}, nil))
	})

	t.Run("nil and empty", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapLabelsEqual(nil, map[string]string{}))
	})

	t.Run("empty and empty", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapLabelsEqual(map[string]string{}, map[string]string{}))
	})

	t.Run("empty and oneEntry", func(t *testing.T) {
		require.False(t, nativekubernetesoo.IsConfigMapLabelsEqual(map[string]string{}, map[string]string{"a": "bc"}))
	})

	t.Run("oneEntry and oneEntry", func(t *testing.T) {
		require.True(t, nativekubernetesoo.IsConfigMapLabelsEqual(map[string]string{"a": "bc"}, map[string]string{"a": "bc"}))
	})
}
