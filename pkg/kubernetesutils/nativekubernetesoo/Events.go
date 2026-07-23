package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	corev1 "k8s.io/api/core/v1"
)

func EventMatchesOptions(event *corev1.Event, options *kubernetesutils.WatchEventOptions) bool {
	return nativekubernetes.EventMatchesOptions(event, options)
}

func WatchEvents(ctx context.Context, options *kubernetesutils.WatchEventOptions, onCreate func(*corev1.Event), onUpdate func(*corev1.Event), onDelete func(*corev1.Event)) error {
	return nativekubernetes.WatchEvents(ctx, options, onCreate, onUpdate, onDelete)
}

func EventToString(event *corev1.Event) string {
	return nativekubernetes.EventToString(event)
}
