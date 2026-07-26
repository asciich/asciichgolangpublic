package nativekubernetesoo_test

import (
	"os"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

var testClusterName string

func TestMain(m *testing.M) {
	ctx := contextutils.ContextVerbose()

	// Create a dedicated cluster for this package's tests
	testClusterName = testutils.GetKindClusterNameForTest(&testing.T{})
	_, err := kindutils.CreateCluster(ctx, testClusterName)
	if err != nil {
		panic("Failed to create test cluster: " + err.Error())
	}

	// Run all tests
	exitCode := m.Run()

	// Teardown: Delete the cluster after all tests
	err = kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, testClusterName)
	if err != nil {
		panic("Failed to delete test cluster: " + err.Error())
	}

	os.Exit(exitCode)
}
