package testcase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testresults"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TestCaseExecutorKubernetesConfigMapYamlValueEqual struct {
	TestCaseExecutorBase
}

func (t *TestCaseExecutorKubernetesConfigMapYamlValueEqual) GetName() (string, error) {
	return "kubernetes_configmap_yamlvalue_equal", nil
}

// formatExpectedValueAsString renders the expected value by its string representation
// to allow scalar comparisons like 8080 == "8080".
func formatExpectedValueAsString(expectedValue any) string {
	if stringValue, ok := expectedValue.(string); ok {
		return stringValue
	}

	return fmt.Sprintf("%v", expectedValue)
}

// toYqPath converts a dot notation yaml path like 'server.listeners[0].port'
// to a yq query. Already qualified paths (starting with '.') are used as is.
func toYqPath(yamlPath string) string {
	if strings.HasPrefix(yamlPath, ".") {
		return yamlPath
	}

	return "." + yamlPath
}

func (t *TestCaseExecutorKubernetesConfigMapYamlValueEqual) Run(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (testutilsinterfaces.TestResult, error) {
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

	configMapKey, err := t.GetConfigMapKey()
	if err != nil {
		return nil, err
	}

	yamlPath, err := t.GetYamlPath()
	if err != nil {
		return nil, err
	}

	expectedValue, err := t.GetExpectedValue()
	if err != nil {
		return nil, err
	}

	expectedValueString := formatExpectedValueAsString(expectedValue)

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

	ns, err := kubernetesCluster.GetNamespaceByName(namespace)
	if err != nil {
		return nil, err
	}

	configMap, err := ns.GetConfigMapByName(configMapName)
	if err != nil {
		return nil, err
	}

	exists, err := configMap.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		tEnd := time.Now()

		baseMessage := fmt.Sprintf(
			"The Kubernetes configmap '%s' in namespace '%s' cluster '%s' does not exist.",
			configMapName,
			namespace,
			cluster,
		)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}

		err = result.SetFailedMessage(failedMessage)
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

	data, err := configMap.GetAllData(ctx)
	if err != nil {
		return nil, err
	}

	yamlContent, keyExists := data[configMapKey]
	if !keyExists {
		tEnd := time.Now()

		baseMessage := fmt.Sprintf(
			"The Kubernetes configmap '%s' in namespace '%s' cluster '%s' has no key '%s' in its data.",
			configMapName,
			namespace,
			cluster,
			configMapKey,
		)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}

		err = result.SetFailedMessage(failedMessage)
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

	if !yamlutils.IsYaml(yamlContent, nil) {
		tEnd := time.Now()

		baseMessage := fmt.Sprintf(
			"The content of key '%s' in Kubernetes configmap '%s' in namespace '%s' cluster '%s' is no valid yaml.",
			configMapKey,
			configMapName,
			namespace,
			cluster,
		)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}

		err = result.SetFailedMessage(failedMessage)
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

	// yaml_path is relative to the yaml document stored in configmap_key:
	value, err := yamlutils.RunYqQueryAginstYamlStringAsString(yamlContent, toYqPath(yamlPath))
	if err != nil {
		return nil, err
	}

	if value == "" || value == "null" {
		tEnd := time.Now()

		baseMessage := fmt.Sprintf(
			"The yaml path '%s' in key '%s' of Kubernetes configmap '%s' in namespace '%s' cluster '%s' resolves to no value.",
			yamlPath,
			configMapKey,
			configMapName,
			namespace,
			cluster,
		)
		failedMessage, err := t.FormatFailedMessage(baseMessage)
		if err != nil {
			return nil, err
		}

		err = result.SetFailedMessage(failedMessage)
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

	tEnd := time.Now()

	if value == expectedValueString {
		successMessage := fmt.Sprintf(
			"The yaml path '%s' in key '%s' of Kubernetes configmap '%s' in namespace '%s' cluster '%s' equals the expected value.",
			yamlPath,
			configMapKey,
			configMapName,
			namespace,
			cluster,
		)
		err = result.SetSuccessMessage(successMessage)
		if err != nil {
			return nil, err
		}
	} else {
		baseMessage := fmt.Sprintf(
			"The yaml path '%s' in key '%s' of Kubernetes configmap '%s' in namespace '%s' cluster '%s' does not equal the expected value. Expected: '%s', actual: '%s'.",
			yamlPath,
			configMapKey,
			configMapName,
			namespace,
			cluster,
			expectedValueString,
			value,
		)
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
