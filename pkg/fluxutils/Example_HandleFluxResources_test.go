package fluxutils_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// Maximum time to wait for asynchronous events (watch callbacks, flux status changes).
	// Generous on purpose: CI runners are much slower than local machines.
	eventuallyTimeout = 90 * time.Second
	eventuallyTick    = 500 * time.Millisecond
)

// waitForCounter waits until the given counter fulfills the check function.
func waitForCounter(t *testing.T, name string, counter *atomic.Int32, check func(int32) bool) {
	t.Helper()

	require.Eventuallyf(
		t,
		func() bool { return check(counter.Load()) },
		eventuallyTimeout,
		eventuallyTick,
		"counter '%s' did not reach the expected value. Current value: %d", name, counter.Load(),
	)
}

// waitForCounterEquals waits until the counter equals the expected value.
func waitForCounterEquals(t *testing.T, name string, counter *atomic.Int32, expected int32) {
	t.Helper()

	waitForCounter(t, name, counter, func(v int32) bool { return v == expected })
}

// waitForCounterAtLeast waits until the counter is greater or equal to the expected minimum.
func waitForCounterAtLeast(t *testing.T, name string, counter *atomic.Int32, minimum int32) {
	t.Helper()

	waitForCounter(t, name, counter, func(v int32) bool { return v >= minimum })
}

// waitForStatusContains polls the status message until it contains the expected string.
func waitForStatusContains(t *testing.T, name string, getStatus func() (string, error), expected string) {
	t.Helper()

	var last atomic.Value
	last.Store("")

	ok := assertEventually(func() bool {
		status, err := getStatus()
		if err != nil {
			last.Store("error: " + err.Error())
			return false
		}
		last.Store(status)
		return contains(status, expected)
	})

	require.Truef(
		t,
		ok,
		"status of '%s' never contained '%s'. Last status: '%v'", name, expected, last.Load(),
	)
}

// assertEventually polls the condition until it returns true or the timeout is reached.
func assertEventually(condition func() bool) bool {
	deadline := time.Now().Add(eventuallyTimeout)
	for {
		if condition() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(eventuallyTick)
	}
}

func contains(s string, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

func indexOf(s string, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// This example shows how to handle flux resources.
//
// To start this test use:
//
//	bash -c "cd pkg/fluxutils && go test -v -run Test_HandleFluxResources"
func Test_HandleFluxResources(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const namespaceName = "flux-system"
	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)
	// Install flux using flux-operator
	_, err = fluxutils.InstallFlux(ctx, &fluxparameteroptions.InstalFluxOptions{
		KubernetesCluster: cluster,
		Namespace:         namespaceName,
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get the deployed flux
	fluxDeployment, err := fluxutils.GetFluxDeployment(cluster, namespaceName)
	require.NoError(t, err)

	// Define example resource names:
	const gitRepoName = "example-repo"
	const kustomizationName = "example-kustomization"
	const helmreleaseName = "example-helmrelease"

	// Ensure example resources absent:
	err = fluxDeployment.DeleteGitRepository(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	err = fluxDeployment.DeleteKustomization(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	err = fluxDeployment.DeleteHelmRelease(ctx, helmreleaseName, namespaceName)
	require.NoError(t, err)

	// Check the example GitRepository is absent before registering callback functions:
	exists, err := fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check the example Kustomization is absent before registering callback functions:
	exists, err = fluxDeployment.KustomizationExists(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check the example HelmRelease is absent before registering callback functions:
	exists, err = fluxDeployment.HelmReleaseExists(ctx, helmreleaseName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Define counters and context to watch the flux objects.
	// The callbacks are called from the watcher goroutines so the counters must be atomic.
	var grCreateCounter, grUpdateCounter, grDeleteCounter atomic.Int32 // gr for GitRepository
	var kuCreateCounter, kuUpdateCounter, kuDeleteCounter atomic.Int32 // ku for Kustomization
	var hrCreateCounter, hrUpdateCounter, hrDeleteCounter atomic.Int32 // hr for HelmRelease
	ctxWatch, cancel := context.WithCancel(ctx)                       // ensure we can cancel the watching
	defer cancel()

	// Start watching the GitRepository
	err = fluxDeployment.WatchGitRepository(
		ctxWatch,
		gitRepoName,
		namespaceName,
		func(*unstructured.Unstructured) { grCreateCounter.Add(1) },
		func(*unstructured.Unstructured) { grUpdateCounter.Add(1) },
		func(*unstructured.Unstructured) { grDeleteCounter.Add(1) },
	)
	require.NoError(t, err)

	// Start watching the Kustomization
	err = fluxDeployment.WatchKustomization(
		ctxWatch,
		kustomizationName,
		namespaceName,
		func(*unstructured.Unstructured) { kuCreateCounter.Add(1) },
		func(*unstructured.Unstructured) { kuUpdateCounter.Add(1) },
		func(*unstructured.Unstructured) { kuDeleteCounter.Add(1) },
	)
	require.NoError(t, err)

	// Start watching the HelmRelease
	err = fluxDeployment.WatchHelmRelease(
		ctxWatch,
		helmreleaseName,
		namespaceName,
		func(*unstructured.Unstructured) { hrCreateCounter.Add(1) },
		func(*unstructured.Unstructured) { hrUpdateCounter.Add(1) },
		func(*unstructured.Unstructured) { hrDeleteCounter.Add(1) },
	)
	require.NoError(t, err)

	// check no callback called for GitRepository
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 0, grCreateCounter.Load())
	require.EqualValues(t, 0, grUpdateCounter.Load())
	require.EqualValues(t, 0, grDeleteCounter.Load())

	// check no callback called for Kustomization
	require.EqualValues(t, 0, kuCreateCounter.Load())
	require.EqualValues(t, 0, kuUpdateCounter.Load())
	require.EqualValues(t, 0, kuDeleteCounter.Load())

	// check no callback called for HelmRelease
	require.EqualValues(t, 0, hrCreateCounter.Load())
	require.EqualValues(t, 0, hrUpdateCounter.Load())
	require.EqualValues(t, 0, hrDeleteCounter.Load())

	// Define a fluxcd GitRepository:
	gitRepoYaml := "---\n"
	gitRepoYaml += "apiVersion: source.toolkit.fluxcd.io/v1\n"
	gitRepoYaml += "kind: GitRepository\n"
	gitRepoYaml += "metadata:\n"
	gitRepoYaml += "  name: " + gitRepoName + "\n"
	gitRepoYaml += "  namespace: " + namespaceName + "\n"
	gitRepoYaml += "spec:\n"
	gitRepoYaml += "  interval: 5m0s\n"
	gitRepoYaml += "  url: https://asciich.ch/example/repo\n"
	gitRepoYaml += "  ref:\n"
	gitRepoYaml += "    branch: master\n"
	// Create the GitRepository:
	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: gitRepoYaml})
	require.NoError(t, err)

	// Check the example GitRepository exists:
	exists, err = fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Check Create counter increased for GitRepository:
	waitForCounterEquals(t, "GitRepository create", &grCreateCounter, 1)
	require.GreaterOrEqual(t, grUpdateCounter.Load(), int32(0)) // Every status change calles an update.
	require.EqualValues(t, 0, grDeleteCounter.Load())

	// Define a fluxcd Kustomization:
	kustomizationYaml := "apiVersion: kustomize.toolkit.fluxcd.io/v1\n"
	kustomizationYaml += "kind: Kustomization\n"
	kustomizationYaml += "metadata:\n"
	kustomizationYaml += "  name: " + kustomizationName + "\n"
	kustomizationYaml += "  namespace: " + namespaceName + "\n"
	kustomizationYaml += "spec:\n"
	kustomizationYaml += "  interval: 10m\n"
	kustomizationYaml += "  targetNamespace: " + namespaceName + "\n"
	kustomizationYaml += "  sourceRef:\n"
	kustomizationYaml += "    kind: GitRepository\n"
	kustomizationYaml += "    name: " + gitRepoName + "\n"
	kustomizationYaml += "  path: \"./kustomize\"\n"
	kustomizationYaml += "  prune: true\n"
	kustomizationYaml += "  timeout: 1m\n"

	// Create the Kustomization:
	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: kustomizationYaml})
	require.NoError(t, err)

	// Check the example Kustomization exists:
	exists, err = fluxDeployment.KustomizationExists(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Check Create counter increased for Kuszomization:
	waitForCounterEquals(t, "Kustomization create", &kuCreateCounter, 1)
	require.GreaterOrEqual(t, kuUpdateCounter.Load(), int32(0)) // Every status update generates an updated version
	require.EqualValues(t, 0, kuDeleteCounter.Load())

	// Define a fluxcd HelmRelease:
	helmreleaseYaml := "---\n"
	helmreleaseYaml += "apiVersion: helm.toolkit.fluxcd.io/v2\n"
	helmreleaseYaml += "kind: HelmRelease\n"
	helmreleaseYaml += "metadata:\n"
	helmreleaseYaml += "  name: " + helmreleaseName + "\n"
	helmreleaseYaml += "  namespace: " + namespaceName + "\n"
	helmreleaseYaml += "spec:\n"
	helmreleaseYaml += "  interval: 10m\n"
	helmreleaseYaml += "  timeout: 5m\n"
	helmreleaseYaml += "  chart:\n"
	helmreleaseYaml += "    spec:\n"
	helmreleaseYaml += "      chart: podinfo\n"
	helmreleaseYaml += "      version: '6.5.*'\n"
	helmreleaseYaml += "      sourceRef:\n"
	helmreleaseYaml += "        kind: HelmRepository\n"
	helmreleaseYaml += "        name: podinfo\n"
	helmreleaseYaml += "      interval: 5m\n"
	helmreleaseYaml += "  releaseName: podinfo\n"
	helmreleaseYaml += "  install:\n"
	helmreleaseYaml += "    remediation:\n"
	helmreleaseYaml += "      retries: 3\n"
	helmreleaseYaml += "  upgrade:\n"
	helmreleaseYaml += "    remediation:\n"
	helmreleaseYaml += "      retries: 3\n"
	helmreleaseYaml += "  test:\n"
	helmreleaseYaml += "    enable: true\n"
	helmreleaseYaml += "  driftDetection:\n"
	helmreleaseYaml += "    mode: enabled\n"
	helmreleaseYaml += "    ignore:\n"
	helmreleaseYaml += "    - paths: [\"/spec/replicas\"]\n"
	helmreleaseYaml += "      target:\n"
	helmreleaseYaml += "        kind: Deployment\n"
	helmreleaseYaml += "  values:\n"
	helmreleaseYaml += "    replicaCount: 2\n"

	// Create the HelmRelease:
	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: helmreleaseYaml})
	require.NoError(t, err)

	// Check the example HelmRelease exists:
	exists, err = fluxDeployment.HelmReleaseExists(ctx, helmreleaseName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Check Create counter increased for HelmRelease:
	waitForCounterEquals(t, "HelmRelease create", &hrCreateCounter, 1)
	require.GreaterOrEqual(t, hrUpdateCounter.Load(), int32(0)) // Every status update generates an updated version
	require.EqualValues(t, 0, hrDeleteCounter.Load())

	// Wait until the resources reached the expected (failing) state.
	// The flux controllers work asynchronously. How long this takes depends on the machine,
	// so poll instead of sleeping a fixed time.

	// Status of the GitRepository:
	waitForStatusContains(
		t,
		"GitRepository",
		func() (string, error) {
			return fluxDeployment.GetGitRepositoryStatusMessage(ctx, gitRepoName, namespaceName)
		},
		"failed to checkout and determine revision: unable to clone ", // The repo of this example does not exist.
	)

	// Status of the Kustomization:
	waitForStatusContains(
		t,
		"Kustomization",
		func() (string, error) {
			return fluxDeployment.GetKustomizationStatusMessage(ctx, kustomizationName, namespaceName)
		},
		"Source artifact not found, retrying in 5s", // The repo of this example does not exist.
	)

	// Status of the HelmRelease:
	waitForStatusContains(
		t,
		"HelmRelease",
		func() (string, error) {
			return fluxDeployment.GetHelmReleaseStatusMessage(ctx, helmreleaseName, namespaceName)
		},
		"latest generation of object has not been reconciled", // The repo of this example does not exist.
	)

	// Check update counter increased for GitRepository:
	waitForCounterEquals(t, "GitRepository create", &grCreateCounter, 1)
	waitForCounterAtLeast(t, "GitRepository update", &grUpdateCounter, 3) // Every status update generates an updated version
	require.EqualValues(t, 0, grDeleteCounter.Load())

	// Check update counter increased for Kustomization:
	waitForCounterEquals(t, "Kustomization create", &kuCreateCounter, 1)
	waitForCounterAtLeast(t, "Kustomization update", &kuUpdateCounter, 2) // Every status update generates an updated version
	require.EqualValues(t, 0, kuDeleteCounter.Load())

	// Check update counter increased for HelmRelease:
	waitForCounterEquals(t, "HelmRelease create", &hrCreateCounter, 1)
	waitForCounterAtLeast(t, "HelmRelease update", &hrUpdateCounter, 2) // Every status update generates an updated version
	require.EqualValues(t, 0, hrDeleteCounter.Load())

	// Delete the example GitRepository:
	err = fluxDeployment.DeleteGitRepository(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)

	// Check if example GitRepository is absent:
	exists, err = fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check delete counter increased for GitRepository:
	waitForCounterEquals(t, "GitRepository delete", &grDeleteCounter, 1)
	require.EqualValues(t, 1, grCreateCounter.Load())
	require.GreaterOrEqual(t, grUpdateCounter.Load(), int32(3)) // Every status update generates an updated version

	// Delete the example Kustomization:
	err = fluxDeployment.DeleteKustomization(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)

	// Check if example Kustomization is absent:
	exists, err = fluxDeployment.KustomizationExists(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check delete counter increased for Kustomization:
	waitForCounterEquals(t, "Kustomization delete", &kuDeleteCounter, 1)
	require.EqualValues(t, 1, kuCreateCounter.Load())
	require.GreaterOrEqual(t, kuUpdateCounter.Load(), int32(2)) // Every status update generates an updated version

	// Delete the example HelmRelease:
	err = fluxDeployment.DeleteHelmRelease(ctx, helmreleaseName, namespaceName)
	require.NoError(t, err)

	// Check if example HelmRelease is absent:
	exists, err = fluxDeployment.HelmReleaseExists(ctx, helmreleaseName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check delete counter increased for HelmRelease:
	waitForCounterEquals(t, "HelmRelease delete", &hrDeleteCounter, 1)
	require.EqualValues(t, 1, hrCreateCounter.Load())
	require.GreaterOrEqual(t, hrUpdateCounter.Load(), int32(2)) // Every status update generates an updated version
}
