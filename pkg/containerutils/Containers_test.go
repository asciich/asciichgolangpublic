package containerutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func getCtx() context.Context {
	return contextutils.ContextSilent()
}

func TestContainersIsRunningInsideContainer(t *testing.T) {
	ctx := getCtx()
	require.False(t, mustutils.Must(containerutils.IsRunningInsideContainer(ctx)))
}
