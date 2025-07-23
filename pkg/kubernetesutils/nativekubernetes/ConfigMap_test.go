package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/nativekubernetes"
)

func TestIsConfigMapContentEqual(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapContentEqual(nil, nil))
	})

	t.Run("empty and nil", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapContentEqual(map[string]string{}, nil))
	})

	t.Run("nil and empty", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapContentEqual(nil, map[string]string{}))
	})

	t.Run("empty and empty", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapContentEqual(map[string]string{}, map[string]string{}))
	})

	t.Run("empty and oneEntry", func(t *testing.T) {
		require.False(t, nativekubernetes.IsConfigMapContentEqual(map[string]string{}, map[string]string{"a": "bc"}))
	})

	t.Run("oneEntry and oneEntry", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapContentEqual(map[string]string{"a": "bc"}, map[string]string{"a": "bc"}))
	})
}

func Test_IsConfigMapLabelsEqual(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapLabelsEqual(nil, nil))
	})

	t.Run("empty and nil", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapLabelsEqual(map[string]string{}, nil))
	})

	t.Run("nil and empty", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapLabelsEqual(nil, map[string]string{}))
	})

	t.Run("empty and empty", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapLabelsEqual(map[string]string{}, map[string]string{}))
	})

	t.Run("empty and oneEntry", func(t *testing.T) {
		require.False(t, nativekubernetes.IsConfigMapLabelsEqual(map[string]string{}, map[string]string{"a": "bc"}))
	})

	t.Run("oneEntry and oneEntry", func(t *testing.T) {
		require.True(t, nativekubernetes.IsConfigMapLabelsEqual(map[string]string{"a": "bc"}, map[string]string{"a": "bc"}))
	})
}
