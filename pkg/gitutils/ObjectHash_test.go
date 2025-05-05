package gitutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetBlobObjectHashFromString(t *testing.T) {
	expected, err := commandexecutor.Bash().RunOneLinerAndGetStdoutAsString(getCtx(), "echo -en 'hello world' | git hash-object --stdin")
	require.NoError(t, err)
	expected = strings.TrimSpace(expected)

	require.EqualValues(t, expected, gitutils.GetBlobOjectHashFromString("hello world"))
}
