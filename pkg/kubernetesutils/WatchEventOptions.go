package kubernetesutils

import "fmt"

type WatchEventOptions struct {
	// Namespace to watch events. If not set all namespaces are watched:
	Namespace string

	// InvolvedObject API version, e.g apiextensions.k8s.io/v1 :
	APIVersion string

	// InvolvedObject Kind:
	Kind string

	// InvolvedObject Name:
	Name string
}

func (w *WatchEventOptions) String() string {
	return fmt.Sprintf("Name='%s' Kind='%s' APIVersion='%s' Namespace='%s'", w.Name, w.Kind, w.APIVersion, w.Namespace)
}