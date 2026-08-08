package nativekubernetesoo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// Test_WaitUntilPodReady_EmptyPodName tests error handling for empty pod name
func Test_WaitUntilPodReady_EmptyPodName(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Create a mock namespace - we can't test actual cluster operations without kind
	// This test verifies the input validation logic
	namespace := &nativekubernetesoo.NativeNamespace{}

	// Wait with empty pod name
	err := namespace.WaitUntilPodReady(ctx, "", 5*time.Second)
	require.Error(t, err, "WaitUntilPodReady should fail for empty pod name")
	require.Contains(t, err.Error(), "empty string", "Error should mention empty string")
}

// Test_WaitUntilPodReady_ContextCancellation tests that context cancellation is respected
func Test_WaitUntilPodReady_ContextCancellation(t *testing.T) {
	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Create a mock namespace
	namespace := &nativekubernetesoo.NativeNamespace{}

	// Wait with cancelled context - should fail quickly due to context cancellation
	err := namespace.WaitUntilPodReady(ctx, "any-pod", 30*time.Second)
	require.Error(t, err, "WaitUntilPodReady should fail with cancelled context")
}
