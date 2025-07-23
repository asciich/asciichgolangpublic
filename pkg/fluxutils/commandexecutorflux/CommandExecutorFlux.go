package commandexecutorflux

import (
	"context"
	"os"
	"time"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/fluxinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/fluxparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils/helminterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils/helmparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type CommandExecutorFlux struct {
	commandExecutor commandexecutor.CommandExecutor
}

func NewcommandExecutorFlux(executor commandexecutor.CommandExecutor) fluxinterfaces.Flux {
	return &CommandExecutorFlux{
		commandExecutor: executor,
	}
}

func (c *CommandExecutorFlux) GetCommandExecutor() (commandexecutor.CommandExecutor, error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return c.commandExecutor, nil
}

func (c *CommandExecutorFlux) GetHelm() (helminterfaces.Helm, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return helmutils.GetCommandExecutorHelm(commandExecutor)
}

func (c *CommandExecutorFlux) InstallFluxOperatorUsingHelm(ctx context.Context, cluster kubernetesinterfaces.KubernetesCluster, namespace string) error {
	if cluster == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	kubeContext, err := cluster.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install flux-operator using helm in namespace '%s' of cluster with context '%s' started.", namespace, kubeContext)

	helm, err := c.GetHelm()
	if err != nil {
		return err
	}

	err = helm.InstallHelmChart(ctx, &helmparameteroptions.InstallHelmChartOptions{
		// Source: https://fluxcd.io/flux/installation/#install-the-flux-operator
		ChartReference:    "flux-operator",
		Namespace:         namespace,
		ChartUri:          "oci://ghcr.io/controlplaneio-fluxcd/charts/flux-operator",
		KubernetesCluster: cluster,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install flux-operator using helm in namespace '%s' of cluster with context '%s' finished.", namespace, kubeContext)

	return nil
}

func (c *CommandExecutorFlux) ConfigureFluxInstance(ctx context.Context, cluster kubernetesinterfaces.KubernetesCluster, namespaceName string) error {
	if cluster == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	kubeContext, err := cluster.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Configuring flux in namespace '%s' of cluster with context '%s' started.", namespaceName, kubeContext)

	resourceYaml := "---\n"
	resourceYaml += "apiVersion: fluxcd.controlplane.io/v1\n"
	resourceYaml += "kind: FluxInstance\n"
	resourceYaml += "metadata:\n"
	resourceYaml += "  name: flux\n"
	resourceYaml += "  namespace: " + namespaceName + "\n"
	resourceYaml += "  annotations:\n"
	resourceYaml += "    fluxcd.controlplane.io/reconcileEvery: \"1h\"\n"
	resourceYaml += "    fluxcd.controlplane.io/reconcileTimeout: \"5m\"\n"
	resourceYaml += "spec:\n"
	resourceYaml += "  distribution:\n"
	resourceYaml += "    version: \"2.x\"\n"
	resourceYaml += "    registry: \"ghcr.io/fluxcd\"\n"
	resourceYaml += "    artifact: \"oci://ghcr.io/controlplaneio-fluxcd/flux-operator-manifests\"\n"
	resourceYaml += "  components:\n"
	resourceYaml += "    - source-controller\n"
	resourceYaml += "    - kustomize-controller\n"
	resourceYaml += "    - helm-controller\n"
	resourceYaml += "    - notification-controller\n"
	resourceYaml += "    - image-reflector-controller\n"
	resourceYaml += "    - image-automation-controller\n"
	resourceYaml += "  cluster:\n"
	resourceYaml += "    type: kubernetes\n"
	resourceYaml += "    multitenant: false\n"
	resourceYaml += "    networkPolicy: true\n"
	resourceYaml += "    domain: \"cluster.local\"\n"
	resourceYaml += "  kustomize:\n"
	resourceYaml += "    patches:\n"
	resourceYaml += "      - target:\n"
	resourceYaml += "          kind: Deployment\n"
	resourceYaml += "          name: \"(kustomize-controller|helm-controller)\"\n"
	resourceYaml += "        patch: |\n"
	resourceYaml += "          - op: add\n"
	resourceYaml += "            path: /spec/template/spec/containers/0/args/-\n"
	resourceYaml += "            value: --concurrent=10\n"
	resourceYaml += "          - op: add\n"
	resourceYaml += "            path: /spec/template/spec/containers/0/args/-\n"
	resourceYaml += "            value: --requeue-dependency=5s\n"

	os.WriteFile("debug.yaml", []byte(resourceYaml), 0644)

	resource, err := cluster.GetObjectByNames("flux", "FluxInstance", namespaceName)
	if err != nil {
		return err
	}

	err = resource.CreateByYamlString(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: resourceYaml,
	})
	if err != nil {
		return err
	}

	timeoutCtx, _ := context.WithTimeout(ctx, time.Minute*2)
	err = cluster.WaitUntilAllPodsInNamespaceAreRunning(timeoutCtx, namespaceName, &kubernetesparameteroptions.WaitForPodsOptions{MinNumberOfPods: 4})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Configuring flux in namespace '%s' of cluster with context '%s' finished.", namespaceName, kubeContext)

	return nil
}

func (c *CommandExecutorFlux) InstallFlux(ctx context.Context, options *fluxparameteroptions.InstalFluxOptions) (fluxinterfaces.FluxDeployment, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	cluster, err := options.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	kubeContext, err := cluster.GetKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := options.GetNamespace()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Install flux in namespace '%s' of kubernetes cluster with context '%s' started.", namespace, kubeContext)

	err = c.InstallFluxOperatorUsingHelm(ctx, cluster, namespace)
	if err != nil {
		return nil, err
	}

	err = c.ConfigureFluxInstance(ctx, cluster, namespace)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Install flux in namespace '%s' of kubernetes cluster with context '%s' finished.", namespace, kubeContext)

	return c.GetDeployedFlux(cluster)
}

func (c *CommandExecutorFlux) GetDeployedFlux(cluster kubernetesinterfaces.KubernetesCluster) (fluxinterfaces.FluxDeployment, error) {
	if cluster == nil {
		return nil, tracederrors.TracedErrorNil("cluster")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	toReturn := new(CommandExecutorDeployedFlux)
	toReturn.commandExecutor = commandExecutor
	toReturn.cluster = cluster

	return toReturn, nil
}
