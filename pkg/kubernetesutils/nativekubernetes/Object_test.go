package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_NativeResurce_GetApiVersion(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		r := &nativekubernetes.NativeObject{}
		apiVersion, err := r.GetApiVersion(getCtx())
		require.NoError(t, err)
		require.EqualValues(t, "v1", apiVersion)
	})

	t.Run("default", func(t *testing.T) {
		r := &nativekubernetes.NativeObject{}

		err := r.SetKind("FluxInstance")
		require.NoError(t, err)

		apiVersion, err := r.GetApiVersion(getCtx())
		require.NoError(t, err)
		require.EqualValues(t, "fluxcd.controlplane.io/v1", apiVersion)
	})
}
