package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorKindClusterExists struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKindClusterExists) GetName() (string, error) {
	return "kind_cluster_exists", nil
}

func (t *TestCaseExecutorKindClusterExists) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
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
	var kind kindutils.Kind
	if isLocalhost {
		kind, err = kindutils.GetLocalCommandExecutorKind()
		if err != nil {
			return nil, err
		}
	} else {
		kind, err = kindutils.GetCommandExecutorKind(commandExecutor)
		if err != nil {
			return nil, err
		}
	}

	exists, err = kind.ClusterByNameExists(ctx, cluster)
	if err != nil {
		return nil, err
	}

	tEnd := time.Now()

	if exists {
		err = result.SetSuccessMessage(
			fmt.Sprintf("The kind cluster '%s' exists.", cluster),
		)
		if err != nil {
			return nil, err
		}
	} else {
		baseMessage := fmt.Sprintf("The kind cluster '%s' does not exist.", cluster)
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
