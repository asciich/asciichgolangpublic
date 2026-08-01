package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
)

type TestCaseExecutorKubernetesConfigMapExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesConfigMapExists) GetName() (string, error) {
	return "kubernetes_configmap_exists", nil
}

func (t *TestCaseExecutorKubernetesConfigMapExists) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	configMapName, err := t.GetResourceName()
	if err != nil {
		return nil, err
	}

	namespace, err := t.GetNamespace()
	if err != nil {
		return nil, err
	}

	cluster, err := t.GetCluster()
	if err != nil {
		return nil, err
	}

	// Get Kubernetes cluster:
	kubernetesCluster, err := nativekubernetesoo.GetClusterByName(ctx, cluster)
	if err != nil {
		return nil, err
	}

	// Get namespace:
	ns, err := kubernetesCluster.GetNamespaceByName(namespace)
	if err != nil {
		return nil, err
	}

	// Check if configMap exists:
	exists, err := ns.ConfigMapByNameExists(ctx, configMapName)
	if err != nil {
		return nil, err
	}

	tEnd := time.Now()

	if !exists {
		// ConfigMap does not exist
		err = result.SetFailedMessage(
			fmt.Sprintf("The Kubernetes configMap '%s' in namespace '%s' cluster '%s' does not exist.", configMapName, namespace, cluster),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// ConfigMap exists
		err = result.SetSuccessMessage(
			fmt.Sprintf("The Kubernetes configMap '%s' in namespace '%s' cluster '%s' exists.", configMapName, namespace, cluster),
		)
		if err != nil {
			return nil, err
		}
	}

	err = result.SetTimeStart(&tStart)
	if err != nil {
		return nil, err
	}

	err = result.SetTimeEnd(&tEnd)
	if err != nil {
		return nil, err
	}

	return result, nil
}
