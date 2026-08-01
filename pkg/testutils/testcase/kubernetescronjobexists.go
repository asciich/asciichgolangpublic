package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
)

type TestCaseExecutorKubernetesCronJobExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesCronJobExists) GetName() (string, error) {
	return "kubernetes_cronjob_exists", nil
}

func (t *TestCaseExecutorKubernetesCronJobExists) Run(ctx context.Context) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	cronJobName, err := t.GetResourceName()
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

	// Check if cronjob exists:
	err = ns.CheckCronJobByNameExists(ctx, cronJobName)
	if err != nil {
		// CronJob does not exist
		tEnd := time.Now()

		err = result.SetFailedMessage(
			fmt.Sprintf("The Kubernetes cronjob '%s' in namespace '%s' cluster '%s' does not exist.", cronJobName, namespace, cluster),
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

	// CronJob exists
	tEnd := time.Now()

	err = result.SetSuccessMessage(
		fmt.Sprintf("The Kubernetes cronjob '%s' in namespace '%s' cluster '%s' exists.", cronJobName, namespace, cluster),
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
