package uuidutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGenerate(t *testing.T) {
	generated := uuidutils.Generate(getCtx())
	require.Len(t, generated, len("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"))
}
