package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
)

type TestCaseExecutorKubernetesSecretExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesSecretExists) GetName() (string, error) {
	return "kubernetes_secret_exists", nil
}

func (t *TestCaseExecutorKubernetesSecretExists) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	secretName, err := t.GetResourceName()
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

	// Check if secret exists:
	err = ns.CheckSecretByNameExists(ctx, secretName)
	if err != nil {
		// Secret does not exist
		tEnd := time.Now()

		err = result.SetFailedMessage(
			fmt.Sprintf("The Kubernetes secret '%s' in namespace '%s' cluster '%s' does not exist.", secretName, namespace, cluster),
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

	// Secret exists
	tEnd := time.Now()

	err = result.SetSuccessMessage(
		fmt.Sprintf("The Kubernetes secret '%s' in namespace '%s' cluster '%s' exists.", secretName, namespace, cluster),
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
