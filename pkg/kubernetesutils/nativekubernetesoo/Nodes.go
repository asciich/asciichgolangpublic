package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
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

	sort.Strings(nodeNames)

	return nodeNames, nil
}
