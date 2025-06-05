package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	v1 "k8s.io/api/core/v1"
)

func Test_EventMatchesOptions(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, nativekubernetes.EventMatchesOptions(nil, nil))
	})

	t.Run("event nil", func(t *testing.T) {
		require.False(t, nativekubernetes.EventMatchesOptions(nil, &kubernetesutils.WatchEventOptions{}))
	})

	t.Run("options nil", func(t *testing.T) {
		// Is like all values unset -> matching all events.
		require.True(t, nativekubernetes.EventMatchesOptions(&v1.Event{}, nil))
	})

	t.Run("event namespace mismatches", func(t *testing.T) {
		event := &v1.Event{}
		event.Namespace = "abc"

		require.False(t, nativekubernetes.EventMatchesOptions(
			event,
			&kubernetesutils.WatchEventOptions{Namespace: "abc"},
		))
	})

	t.Run("event involvedOpject namespace", func(t *testing.T) {
		event := &v1.Event{}
		event.InvolvedObject.Namespace = "abc"

		require.True(t, nativekubernetes.EventMatchesOptions(
			event,
			&kubernetesutils.WatchEventOptions{Namespace: "abc"},
		))
	})

	t.Run("kind matches", func(t *testing.T) {
		event := &v1.Event{}
		event.InvolvedObject.Kind = "abc"

		require.True(t, nativekubernetes.EventMatchesOptions(
			event,
			&kubernetesutils.WatchEventOptions{InvolvedObjectKind: "abc"},
		))
	})

	t.Run("kind matches ignorecase", func(t *testing.T) {
		event := &v1.Event{}
		event.InvolvedObject.Kind = "ABC"

		require.True(t, nativekubernetes.EventMatchesOptions(
			event,
			&kubernetesutils.WatchEventOptions{InvolvedObjectKind: "abc"},
		))
	})
}
