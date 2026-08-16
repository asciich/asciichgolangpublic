package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorKubernetesReplicaSetExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesReplicaSetExists) GetName() (string, error) {
	return "kubernetes_replicaset_exists", nil
}

func (t *TestCaseExecutorKubernetesReplicaSetExists) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
	tStart := time.Now()

	name, err := t.GetTestCaseName()
	if err != nil {
		return nil, err
	}

	result := &testresults.TestCaseResult{
		Name: name,
	}

	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	replicaSetName, err := t.GetResourceName()
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

	// Check if running on localhost or remote
	isLocalhost, err := commandExecutor.IsRunningOnLocalhost()
	if err != nil {
		return nil, err
	}

	var exists bool
	var kubernetesCluster kubernetesinterfaces.KubernetesCluster
	if isLocalhost {
		kubernetesCluster, err = nativekubernetesoo.GetClusterByName(ctx, cluster)
		if err != nil {
			return nil, err
		}
	} else {
		kubernetesCluster, err = commandexecutorkubernetes.GetCommandExecutorKubernetsByName(commandExecutor, cluster)
		if err != nil {
			return nil, err
		}
	}

	ns, err := kubernetesCluster.GetNamespaceByName(namespace)
	if err != nil {
		return nil, err
	}

	exists, err = ns.ReplicaSetByNameExists(ctx, replicaSetName)
	if err != nil {
		return nil, err
	}

	tEnd := time.Now()

	if exists {
		err = result.SetSuccessMessage(
			fmt.Sprintf("The Kubernetes replicaset '%s' in namespace '%s' cluster '%s' exists.", replicaSetName, namespace, cluster),
		)
		if err != nil {
			return nil, err
		}
	} else {
		baseMessage := fmt.Sprintf("The Kubernetes replicaset '%s' in namespace '%s' cluster '%s' does not exist.", replicaSetName, namespace, cluster)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}
		err = result.SetFailedMessage(failedMessage)
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
