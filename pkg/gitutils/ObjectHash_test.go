package gitutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/gitutils"
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
