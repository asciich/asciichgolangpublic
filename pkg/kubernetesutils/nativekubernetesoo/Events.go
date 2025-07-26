package nativekubernetesoo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func EventMatchesOptions(event *corev1.Event, options *kubernetesutils.WatchEventOptions) bool {
	if event == nil {
		return false
	}

	if options == nil {
		return true
	}

	if options.Namespace != "" {
		if !strings.EqualFold(event.InvolvedObject.Namespace, options.Namespace) {
			return false
		}
	}

	if options.InvolvedObjectAPIVersion != "" {
		if !strings.EqualFold(event.InvolvedObject.APIVersion, options.InvolvedObjectAPIVersion) {
			return false
		}
	}

	if options.InvolvedObjectKind != "" {
		if !strings.EqualFold(event.InvolvedObject.Kind, options.InvolvedObjectKind) {
			return false
		}
	}

	if options.InvolvedObjectName != "" {
		if !strings.EqualFold(event.InvolvedObject.Name, options.InvolvedObjectName) {
			return false
		}
	}

	return true
}

func notifyCallbackWithEvent(event *corev1.Event, options *kubernetesutils.WatchEventOptions, callback func(*corev1.Event)) {
	if event == nil {
		return
	}

	if EventMatchesOptions(event, options) {
		callback(event)
	}
}

func WatchEvents(ctx context.Context, options *kubernetesutils.WatchEventOptions, onCreate func(*corev1.Event), onUpdate func(*corev1.Event), onDelete func(*corev1.Event)) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	tStart := time.Now()

	logging.LogInfoByCtxf(ctx, "Watch kubernetes events with options='%s' setup started.", options)

	clientset, err := GetClientSet(ctx, options.ClusterName)
	if err != nil {
		return err
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	eventInformer := factory.Core().V1().Events().Informer()

	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if onCreate == nil {
				return
			}

			event, ok := obj.(*corev1.Event)
			if !ok {
				return
			}

			if event.LastTimestamp.Time.Before(tStart) {
				return
			}

			notifyCallbackWithEvent(event, options, onCreate)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if onUpdate == nil {
				return
			}

			event, ok := newObj.(*corev1.Event)
			if !ok {
				return
			}

			if event.LastTimestamp.Time.Before(tStart) {
				return
			}

			notifyCallbackWithEvent(event, options, onUpdate)
		},
		// DeleteFunc is called when an event is deleted.
		DeleteFunc: func(obj interface{}) {
			if onDelete == nil {
				return
			}

			event, ok := obj.(*corev1.Event)
			if !ok {
				return
			}

			if event.LastTimestamp.Time.Before(tStart) {
				return
			}

			notifyCallbackWithEvent(event, options, onDelete)
		},
	})

	factory.Start(ctx.Done())
	factory.WaitForCacheSync(ctx.Done())

	logging.LogInfoByCtxf(ctx, "Watch kubernetes events with options='%s' setup finished. Matching events are now watched", options)

	return nil
}

func EventToString(event *corev1.Event) string {
	if event == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"K8s event: Timestamp: %s Namespace: %s Name: %s Kind: %s Reason: %s Message: %s",
		event.LastTimestamp.Format(time.RFC3339),
		event.InvolvedObject.Namespace,
		event.InvolvedObject.Name,
		event.InvolvedObject.Kind,
		event.Reason,
		event.Message,
	)
}
