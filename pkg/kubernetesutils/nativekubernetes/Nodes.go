package nativekubernetes

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListNodeNames(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeNames := []string{}
	for _, n := range nodes.Items {
		nodeNames = append(nodeNames, n.Name)
	}

	logging.LogInfoByCtxf(ctx, "The kubernetes cluster has '%d' nodes.", len(nodeNames))

	return nodeNames, nil
}
