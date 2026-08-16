package testcase

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorKubernetesValidateSshKeyInSecret struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesValidateSshKeyInSecret) GetName() (string, error) {
	return "kubernetes_validate_ssh_key_in_secret", nil
}
func (t *TestCaseExecutorKubernetesValidateSshKeyInSecret) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
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

	secretKey, err := t.GetSecretKey()
	if err != nil {
		return nil, err
	}

	targetHost, err := t.GetTargetHost()
	if err != nil {
		return nil, err
	}

	targetUser, err := t.GetTargetUser()
	if err != nil {
		return nil, err
	}

	targetPort, err := t.GetTargetPort()
	if err != nil {
		return nil, err
	}

	// Check if running on localhost or remote
	isLocalhost, err := commandExecutor.IsRunningOnLocalhost()
	if err != nil {
		return nil, err
	}

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

	options := &kubernetesparameteroptions.ValidateSshKeyInSecretOptions{
		Namespace:             namespace,
		SecretName:            secretName,
		SecretKey:             secretKey,
		TargetHost:            targetHost,
		TargetUser:            targetUser,
		TargetPort:            targetPort,
		SkipHostKeyValidation: true,
		ConnectionTimeout:     "10",
		ConnectionAttempts:    1,
	}

	success, err := kubernetesCluster.ValidateSSHKeyInSecret(ctx, options)

	tEnd := time.Now()

	// If there's an error (e.g. secret not found), treat it as a failed test case, not an infrastructure error
	if err != nil {
		baseMessage := fmt.Sprintf("SSH key validation failed for secret '%s' in namespace '%s' cluster '%s': %v",
			secretName, namespace, cluster, err)
		failedMessage, formatErr := t.FormatFailedMessage(baseMessage)
		if formatErr != nil {
			return nil, formatErr
		}
		err = result.SetFailedMessage(failedMessage)
		if err != nil {
			return nil, err
		}
	} else if success {
		err = result.SetSuccessMessage(
			fmt.Sprintf("The SSH key in secret '%s' (key '%s') in namespace '%s' cluster '%s' successfully authenticated to '%s@%s:%d'.",
				secretName, secretKey, namespace, cluster, targetUser, targetHost, targetPort),
		)
		if err != nil {
			return nil, err
		}
	} else {
		baseMessage := fmt.Sprintf("The SSH key in secret '%s' (key '%s') in namespace '%s' cluster '%s' failed to authenticate to '%s@%s:%d'.",
			secretName, secretKey, namespace, cluster, targetUser, targetHost, targetPort)
		failedMessage, formatErr := t.FormatFailedMessage(baseMessage)
		if formatErr != nil {
			return nil, formatErr
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
