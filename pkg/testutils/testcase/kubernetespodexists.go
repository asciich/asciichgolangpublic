package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
)

type TestCaseExecutorKubernetesPodExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesPodExists) GetName() (string, error) {
	return "kubernetes_pod_exists", nil
}

func (t *TestCaseExecutorKubernetesPodExists) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	podName, err := t.GetResourceName()
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

	// Check if pod exists:
	err = ns.CheckPodByNameExists(ctx, podName)
	if err != nil {
		// Pod does not exist
		tEnd := time.Now()

		err = result.SetFailedMessage(
			fmt.Sprintf("The Kubernetes pod '%s' in namespace '%s' cluster '%s' does not exist.", podName, namespace, cluster),
		)
		if err != nil {
			return nil, err
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

	// Pod exists
	tEnd := time.Now()

	err = result.SetSuccessMessage(
		fmt.Sprintf("The Kubernetes pod '%s' in namespace '%s' cluster '%s' exists.", podName, namespace, cluster),
	)
	if err != nil {
		return nil, err
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
