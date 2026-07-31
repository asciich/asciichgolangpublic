package commandexecutorkubernetes

import (
	"context"
	"encoding/json"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorKubernetes struct {
	name              string
	commandExecutor   commandexecutorinterfaces.CommandExecutor
	cachedContextName string
}

func GetCommandExecutorKubernetsByName(commandExecutor commandexecutorinterfaces.CommandExecutor, clusterName string) (kubernetes kubernetesinterfaces.KubernetesCluster, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	toReturn := NewCommandExecutorKubernetes()

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(clusterName)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetClusterByName(clusterName string) (kubernetes kubernetesinterfaces.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	return GetCommandExecutorKubernetsByName(commandexecutorbashoo.Bash(), clusterName)
}

func NewCommandExecutorKubernetes() (c *CommandExecutorKubernetes) {
	return new(CommandExecutorKubernetes)
}

// Returns the kubernetes cluster name
func (c *CommandExecutorKubernetes) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorKubernetes) CreateNamespaceByName(ctx context.Context, name string) (createdNamespace kubernetesinterfaces.Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(ctx, name)
	if err != nil {
		return nil, err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already exists in cluster '%s'.", name, clusterName)
	} else {
		cmd := []string{"kubectl"}

		if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
			logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. cluster context is not used.")
		} else {
			kubectlContext, err := c.GetCachedKubectlContext(ctx)
			if err != nil {
				return nil, err
			}

			cmd = append(cmd, "--context", kubectlContext)
		}

		cmd = append(cmd, "create", "namespace", name)

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: cmd,
			},
		)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' in cluster '%s' created.", name, clusterName)
	}

	err = c.WaitForDefaultServiceAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	return c.GetNamespaceByName(name)
}

func (c *CommandExecutorKubernetes) WaitForDefaultServiceAccount(ctx context.Context, namespaceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	logging.LogInfoByCtxf(ctx, "Wait for default service account in namespace '%s' started.", namespaceName)

	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	for {
		cmd := []string{"kubectl"}

		if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
			// No context needed for in-cluster auth
		} else {
			kubectlContext, err := c.GetCachedKubectlContext(ctx)
			if err != nil {
				return err
			}

			cmd = append(cmd, "--context", kubectlContext)
		}

		cmd = append(cmd, "get", "serviceaccount", "default", "--namespace", namespaceName, "-o", "name")

		output, err := c.RunCommand(
			timeoutCtx,
			&parameteroptions.RunCommandOptions{
				Command:           cmd,
				AllowAllExitCodes: true,
			},
		)
		if err != nil {
			return err
		}

		if output.IsExitSuccess() {
			logging.LogInfoByCtxf(ctx, "Default service account in namespace '%s' is available.", namespaceName)
			return nil
		}

		stderr, err := output.GetStderrAsString()
		if err != nil {
			return err
		}

		if strings.Contains(stderr, "not found") {
			logging.LogInfoByCtxf(ctx, "Default service account in namespace '%s' not yet available, waiting...", namespaceName)
		} else {
			return tracederrors.TracedErrorf("Failed to check default service account in namespace '%s': %s", namespaceName, stderr)
		}

		select {
		case <-timeoutCtx.Done():
			return tracederrors.TracedErrorf("Timed out waiting for default service account in namespace '%s': %w", namespaceName, timeoutCtx.Err())
		case <-time.After(1 * time.Second):
			// retry
		}
	}
}

func (c *CommandExecutorKubernetes) DeleteNamespaceByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(ctx, name)
	if err != nil {
		return err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return err
	}

	if exists {

		context, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return err
		}

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"kubectl",
					"--context",
					context,
					"delete",
					"namespace",
					name,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' in cluster '%s' deleted.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already absent in cluster '%s'.", name, clusterName)
	}

	return nil
}

func (c *CommandExecutorKubernetes) GetCachedContextName() (cachedContextName string, err error) {
	if c.cachedContextName == "" {
		return "", tracederrors.TracedErrorf("cachedContextName not set")
	}

	return c.cachedContextName, nil
}

func (c *CommandExecutorKubernetes) GetCachedKubectlContext(ctx context.Context) (context string, err error) {
	if c.cachedContextName == "" {
		return c.GetKubectlContext(ctx)
	}

	context = c.cachedContextName

	clusterName, err := c.GetName()
	if err != nil {
		return "", err
	}

	logging.LogInfof(
		"Kubectl context for cluster '%s' is '%s'.",
		clusterName,
		context,
	)

	return
}

func (c *CommandExecutorKubernetes) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("CommandExecutor not set")
	}

	return c.commandExecutor, nil
}

func (c *CommandExecutorKubernetes) GetKubectlContext(ctx context.Context) (context string, err error) {
	contexts, err := c.GetKubectlContexts()
	if err != nil {
		return "", err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return "", err
	}

	for _, con := range contexts {
		clusterNameToCeck, err := con.GetCluster()
		if err != nil {
			return "", err
		}

		if clusterNameToCeck == clusterName {
			context, err = con.GetName()
			if err != nil {
				return "", err
			}

			logging.LogInfoByCtxf(ctx, "Kubectl context for cluster '%s' is '%s'.", clusterName, context)

			return context, nil
		}
	}

	return "", tracederrors.TracedErrorf(
		"No kubectl context for cluster '%s' found.",
		clusterName,
	)
}

func (c *CommandExecutorKubernetes) GetKubectlContexts() (contexts []kubernetesutils.KubectlContext, err error) {
	lines, err := c.RunCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"kubectl", "config", "get-contexts", "--no-headers"},
		},
	)
	if err != nil {
		return nil, err
	}

	contexts = []kubernetesutils.KubectlContext{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", " ")
		line = stringsutils.RepeatReplaceAll(line, "  ", " ")
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		splitted := strings.Split(line, " ")
		if len(splitted) <= 2 {
			return nil, tracederrors.TracedErrorf(
				"Unable to get context from line: '%s'",
				line,
			)
		}

		toAdd := kubernetesutils.NewKubectlContext()
		err = toAdd.SetName(splitted[0])
		if err != nil {
			return nil, err
		}

		err = toAdd.SetCluster(splitted[1])
		if err != nil {
			return nil, err
		}

		contexts = append(contexts, *toAdd)
	}

	return contexts, nil
}

func (c *CommandExecutorKubernetes) GetPodByNames(namespaceName string, podName string) (kubernetesinterfaces.Pod, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetPodByName(podName)
}

func (c *CommandExecutorKubernetes) GetReplicaSetByNames(namespaceName string, replicaSetName string) (kubernetesinterfaces.ReplicaSet, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetReplicaSetByName(replicaSetName)
}

func (c *CommandExecutorKubernetes) GetDeploymentByNames(namespaceName string, deploymentName string) (kubernetesinterfaces.Deployment, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetDeploymentByName(deploymentName)
}

func (c *CommandExecutorKubernetes) GetNamespaceByName(name string) (namespace kubernetesinterfaces.Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	toReturn := NewCommandExecutorNamespace()

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetKubernetesCluster(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorKubernetes) GetObjectByNames(objectName string, objectType string, namespaceName string) (object kubernetesinterfaces.Object, err error) {
	if objectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectName")
	}

	if objectType == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectType")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetObjectByNames(objectName, objectType)
}

func (c *CommandExecutorKubernetes) ListNamespaceNames(ctx context.Context) (namespaceNames []string, err error) {
	namespaces, err := c.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespaceNames = []string{}
	for _, namespace := range namespaces {
		toAdd, err := namespace.GetName()
		if err != nil {
			return nil, err
		}

		namespaceNames = append(namespaceNames, toAdd)
	}

	return namespaceNames, nil
}

func (c *CommandExecutorKubernetes) ListNamespaces(ctx context.Context) (namespaces []kubernetesinterfaces.Namespace, err error) {

	cmd := []string{"kubectl"}

	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext for ListNamespaces.")
	} else {
		context, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, "--context", context)
	}

	cmd = append(cmd, "get", "namespaces", "-o", "name")

	lines, err := c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return nil, err
	}

	namespaces = []kubernetesinterfaces.Namespace{}
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		line = strings.TrimPrefix(line, "namespace/")

		toAdd, err := c.GetNamespaceByName(line)
		if err != nil {
			return nil, err
		}

		namespaces = append(namespaces, toAdd)
	}

	return namespaces, nil
}

func (c *CommandExecutorKubernetes) ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objectNames []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	namespaceName, err := options.GetNamespace()
	if err != nil {
		return nil, err
	}

	context, err := c.GetKubectlContext(contextutils.GetVerbosityContextByBool(options.Verbose))
	if err != nil {
		return nil, err
	}

	objectType, err := options.GetObjectType()
	if err != nil {
		return nil, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		contextutils.GetVerbosityContextByBool(options.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"get",
				"--context",
				context,
				"--namespace",
				namespaceName,
				"-o",
				"name",
				objectType,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	objectNames = []string{}
	for _, name := range output {
		objectNames = append(objectNames, strings.TrimPrefix(name, objectType+"/"))
	}

	sort.Strings(objectNames)

	return objectNames, nil
}

func (c *CommandExecutorKubernetes) ListObjects(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) (objects []kubernetesinterfaces.Object, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	objectNames, err := c.ListObjectNames(options)
	if err != nil {
		return nil, err
	}

	namespaceName, err := options.GetNamespace()
	if err != nil {
		return nil, err
	}

	objectType, err := options.GetObjectType()
	if err != nil {
		return nil, err
	}

	objects = []kubernetesinterfaces.Object{}
	for _, name := range objectNames {
		toAdd, err := c.GetObjectByNames(name, objectType, namespaceName)
		if err != nil {
			return nil, err
		}

		objects = append(objects, toAdd)
	}

	return objects, nil
}

func (c *CommandExecutorKubernetes) NamespaceByNameExists(ctx context.Context, name string) (exists bool, err error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	namespaces, err := c.ListNamespaceNames(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(namespaces, name)

	clusterName, err := c.GetName()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' exists in kubernetes cluster '%s'.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' does not exist in kubernetes cluster '%s'.", name, clusterName)
	}

	return exists, nil
}

func (c *CommandExecutorKubernetes) RunCommand(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runCommandOptions)
}

func (c *CommandExecutorKubernetes) RunCommandAndGetStdoutAsLines(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (lines []string, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	output, err := c.RunCommand(ctx, runCommandOptions)
	if err != nil {
		return nil, err
	}

	return output.GetStdoutAsLines(false)
}

func (c *CommandExecutorKubernetes) SetCachedContextName(cachedContextName string) (err error) {
	if cachedContextName == "" {
		return tracederrors.TracedErrorf("cachedContextName is empty string")
	}

	c.cachedContextName = cachedContextName

	return nil
}

func (c *CommandExecutorKubernetes) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorKubernetes) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorKubernetes) CreateSecret(ctx context.Context, namespaceName string, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret kubernetesinterfaces.Secret, err error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateSecret(ctx, secretName, options)
}

func (c *CommandExecutorKubernetes) SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.SecretByNameExists(ctx, secretName)
}

// CheckSecretByNameExists checks if a secret exists by name.
// Returns nil if it exists, error if it does not exist.
func (c *CommandExecutorKubernetes) CheckSecretByNameExists(ctx context.Context, namespaceName string, secretName string) error {
	exists, err := c.SecretByNameExists(ctx, namespaceName, secretName)
	if err != nil {
		return err
	}
	if !exists {
		tracederrors.TracedErrorf("Secret '%s' does not exist in namespace '%s'", secretName, namespaceName)
	}
	return nil
}

func (c *CommandExecutorKubernetes) DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteSecretByName(ctx, secretName)
}

func (c *CommandExecutorKubernetes) CreateConfigMap(ctx context.Context, namespaceName string, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap kubernetesinterfaces.ConfigMap, err error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateConfigMap(ctx, configMapName, options)
}

func (c *CommandExecutorKubernetes) ConfigMapByNameExists(ctx context.Context, namespaceName string, configMapName string) (exists bool, err error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.ConfigMapByNameExists(ctx, configMapName)
}

func (c *CommandExecutorKubernetes) DeleteConfigMapByName(ctx context.Context, namespaceName string, configMapName string) (err error) {
	namspace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namspace.DeleteConfigMapByName(ctx, configMapName)
}

func (c *CommandExecutorKubernetes) CheckAccessible(ctx context.Context) error {
	clusterName, err := c.GetName()
	if err != nil {
		return err
	}

	_, err = c.WhoAmI(ctx)
	if err != nil {
		return tracederrors.TracedErrorf("Cluster '%s' is not reachable.", clusterName)
	}

	logging.LogInfoByCtxf(ctx, "Cluster '%s' is reachable.", clusterName)

	return err
}

func (c *CommandExecutorKubernetes) WhoAmI(ctx context.Context) (*kubernetesimplementationindependend.UserInfo, error) {
	executor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	kubeContext, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	stdout, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"kubectl", "--context", kubeContext, "auth", "whoami", "-ojson"},
	})
	if err != nil {
		return nil, err
	}

	userName, err := jsonutils.RunJqAgainstJsonStringAsString(stdout, ".status.userInfo.username")
	if err != nil {
		return nil, err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Whoami: Kube context '%s' uses user '%s' to log in to cluster '%s'.", kubeContext, userName, clusterName)

	return &kubernetesimplementationindependend.UserInfo{
		Username: userName,
	}, nil
}

func (c *CommandExecutorKubernetes) WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.WaitForPodsOptions) error {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKubernetes) CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (kubernetesinterfaces.Object, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

// RunCommandInTemporaryPod runs a command in a temporary Kubernetes pod and returns the output.
//
// Implementation note — why we use a run → wait → logs → delete approach instead of `kubectl run --rm -i`:
//
// Using `kubectl run` with `-i` or `--attach` causes output to be printed twice. This is a
// known bug tracked in https://github.com/kubernetes/kubernetes/issues/27264. The root cause
// is a race condition: kubectl first attaches to the container's stdout/stderr while it is
// running, and then fetches the pod logs again as a fallback to ensure no output was missed
// during the attach phase. Since kubectl cannot determine whether the attach captured all
// output, it defensively reads the logs a second time — resulting in every line appearing twice.
//
// To avoid this, we decouple the lifecycle into four explicit steps:
//  1. `kubectl run`    — start the pod without attaching
//  2. `kubectl wait`   — block until the pod has completed
//  3. `kubectl logs`   — fetch the output exactly once
//  4. `kubectl delete` — clean up the pod
func (c *CommandExecutorKubernetes) RunCommandInTemporaryPod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	podName, err := options.GetPodName()
	if err != nil {
		return nil, err
	}

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	commandToExecute, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' started.", podName, namespaceName, imageName)

	kubeContext, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	// Step 1: Start the pod without attaching.
	runCommand := []string{
		"kubectl", "--context", kubeContext, "run", podName,
		"--namespace", namespaceName,
		"--image", imageName,
		"--restart=Never",
	}

	// Build overrides for secrets (environment variables and mounts)
	overrides := buildPodOverridesForSecrets(podName, options)
	if overrides != "" {
		runCommand = append(runCommand,
			"--override-type", "strategic",
			"--overrides", overrides,
		)
	}

	runCommand = append(runCommand, "--")
	runCommand = append(runCommand, commandToExecute...)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: runCommand,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: Wait for the pod to complete.
	waitCommand := []string{
		"kubectl", "--context", kubeContext, "wait", "pod", podName,
		"--namespace", namespaceName,
		"--for=condition=Ready",
		"--timeout=60s",
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: waitCommand,
	})
	if err != nil {
		return nil, err
	}

	// Step 3: Fetch the logs exactly once.
	logsCommand := []string{
		"kubectl", "--context", kubeContext, "logs", podName,
		"--namespace", namespaceName,
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: logsCommand,
	})
	if err != nil {
		return nil, err
	}

	// Step 4: Delete the pod to clean up.
	deleteCommand := []string{
		"kubectl", "--context", kubeContext, "delete", "pod", podName,
		"--namespace", namespaceName,
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: deleteCommand,
	})
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary pod '%s' in namespace '%s' using container image '%s' finished.", podName, namespaceName, imageName)

	return output, nil
}

func (c *CommandExecutorKubernetes) ReadSecret(ctx context.Context, namespaceName string, secretName string) (map[string][]byte, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKubernetes) ListNodeNames(ctx context.Context) ([]string, error) {
	commandexecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	context, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	output, err := commandexecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{Command: []string{"kubectl", "--context", context, "get", "nodes", "-o", "name"}})
	if err != nil {
		return nil, err
	}

	nodeNames := []string{}
	for _, line := range stringsutils.SplitLines(output, true) {
		toAdd := strings.TrimSpace(strings.TrimPrefix(line, "node/"))
		if toAdd == "" {
			continue
		}

		nodeNames = append(nodeNames, toAdd)
	}

	logging.LogInfoByCtxf(ctx, "The kubernetes cluster has '%d' nodes.", len(nodeNames))

	return nodeNames, nil
}

func (c *CommandExecutorKubernetes) DeletePodByNames(ctx context.Context, namespaceName string, podName string) error {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeletePodByName(ctx, podName)
}

func (c *CommandExecutorKubernetes) DeleteReplicaSetByNames(ctx context.Context, namespaceName string, replicaSetName string) error {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteReplicaSetByName(ctx, replicaSetName)
}

func (c *CommandExecutorKubernetes) DeleteDeploymentByNames(ctx context.Context, namespaceName string, deploymentName string) error {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return err
	}

	return namespace.DeleteDeploymentByName(ctx, deploymentName)
}

func (c *CommandExecutorKubernetes) CreatePod(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.Pod, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreatePod(ctx, options)
}

func (c *CommandExecutorKubernetes) CreateReplicaSet(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.ReplicaSet, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateReplicaSet(ctx, options)
}

func (c *CommandExecutorKubernetes) CreateDeployment(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.Deployment, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.CreateDeployment(ctx, options)
}

func (c *CommandExecutorKubernetes) PodByNameExists(ctx context.Context, namespaceName string, podName string) (bool, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.PodByNameExists(ctx, podName)
}

func (c *CommandExecutorKubernetes) ReplicaSetByNameExists(ctx context.Context, namespaceName string, replicaSetName string) (bool, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.ReplicaSetByNameExists(ctx, replicaSetName)
}

func (c *CommandExecutorKubernetes) DeploymentByNameExists(ctx context.Context, namespaceName string, deploymentName string) (bool, error) {
	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return false, err
	}

	return namespace.DeploymentByNameExists(ctx, deploymentName)
}

// toJsonEnvVarFromSecretsForKubectl builds JSON for environment variables sourced from secrets
// This function is specifically designed for kubectl --overrides format
func toJsonEnvVarFromSecretsForKubectl(secretEnvVars map[string]kubernetesparameteroptions.SecretEnvVarSource) string {
	if len(secretEnvVars) == 0 {
		return "[]"
	}

	type envVar struct {
		Name      string `json:"name"`
		ValueFrom struct {
			SecretKeyRef struct {
				Name string `json:"name"`
				Key  string `json:"key"`
			} `json:"secretKeyRef"`
		} `json:"valueFrom"`
	}

	envVars := make([]envVar, 0, len(secretEnvVars))
	for envVarName, secretSource := range secretEnvVars {
		ev := envVar{
			Name: envVarName,
		}
		ev.ValueFrom.SecretKeyRef.Name = secretSource.SecretName
		ev.ValueFrom.SecretKeyRef.Key = secretSource.SecretKey
		envVars = append(envVars, ev)
	}

	jsonBytes, err := json.Marshal(envVars)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

// buildPodOverridesForSecrets builds JSON overrides for pod spec including
// both environment variables from secrets and secret volume mounts
func buildPodOverridesForSecrets(podName string, options *kubernetesparameteroptions.RunCommandOptions) string {
	containerOverrides := map[string]interface{}{
		"name": podName,
	}

	hasEnvVars := len(options.SecretEnvVars) > 0
	hasMounts := len(options.SecretMounts) > 0

	if !hasEnvVars && !hasMounts {
		return ""
	}

	// Build environment variables
	if hasEnvVars {
		envVarsJSON := toJsonEnvVarFromSecretsForKubectl(options.SecretEnvVars)
		var envVars []interface{}
		json.Unmarshal([]byte(envVarsJSON), &envVars)
		containerOverrides["env"] = envVars
	}

	// Build volume mounts and volumes
	if hasMounts {
		volumeMounts := []map[string]interface{}{}
		volumes := []map[string]interface{}{}

		for mountPath, secretSource := range options.SecretMounts {
			volumeName := "secret-" + secretSource.SecretName

			// Volume mount for the container
			volumeMounts = append(volumeMounts, map[string]interface{}{
				"name":      volumeName,
				"mountPath": mountPath,
				"readOnly":  true,
			})

			// Volume definition for the pod
			volumes = append(volumes, map[string]interface{}{
				"name": volumeName,
				"secret": map[string]interface{}{
					"secretName": secretSource.SecretName,
				},
			})
		}

		containerOverrides["volumeMounts"] = volumeMounts

		// Return full pod spec override with both containers and volumes
		overrides := map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{containerOverrides},
				"volumes":    volumes,
			},
		}

		jsonBytes, err := json.Marshal(overrides)
		if err != nil {
			return ""
		}
		return string(jsonBytes)
	}

	// Only environment variables (legacy behavior)
	overrides := map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []map[string]interface{}{containerOverrides},
		},
	}

	jsonBytes, err := json.Marshal(overrides)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// ListKindNames retrieves a sorted list of all available resource kind names
// from the Kubernetes API server. It uses `kubectl api-resources` to query the
// server's available resources across all API groups and versions, and returns
// their kind names in alphabetical order.
//
// Returns:
//   - []string: A sorted slice of unique resource kind names (e.g., "Pod", "Service", "Deployment").
//   - error: An error if the kubectl command fails or the API server cannot be queried.
func (c *CommandExecutorKubernetes) ListKindNames(ctx context.Context) ([]string, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	cmd := []string{"kubectl"}

	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext for ListKindNames.")
	} else {
		kubeContext, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, "--context", kubeContext)
	}

	cmd = append(cmd, "api-resources", "--no-headers", "-o", "name")

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return nil, err
	}

	// Use a map to deduplicate kind names
	kindSet := map[string]bool{}

	// Parse the output using a second call that includes the KIND column
	cmd2 := []string{"kubectl"}

	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		// No context needed
	} else {
		kubeContext, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return nil, err
		}

		cmd2 = append(cmd2, "--context", kubeContext)
	}

	cmd2 = append(cmd2, "api-resources", "--no-headers")

	output, err = commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd2,
	})
	if err != nil {
		return nil, err
	}

	for _, line := range stringsutils.SplitLines(output, true) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		// The KIND column is always the last field in `kubectl api-resources --no-headers`
		kind := fields[len(fields)-1]
		kindSet[kind] = true
	}

	apiKinds := []string{}
	for kind := range kindSet {
		apiKinds = append(apiKinds, kind)
	}

	sort.Strings(apiKinds)

	return apiKinds, nil
}

// CronJob stub implementations - not yet implemented for command executor
func (c *CommandExecutorKubernetes) CronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKubernetes) CreateCronJob(ctx context.Context, namespaceName string, cronJobName string, schedule string, image string, command []string, labels map[string]string) (kubernetesinterfaces.CronJob, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKubernetes) DeleteCronJobByName(ctx context.Context, namespaceName string, cronJobName string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}
