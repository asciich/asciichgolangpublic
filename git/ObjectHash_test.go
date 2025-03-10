package git

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
)

func Test_GetBlobObjectHashFromString(t *testing.T) {
	expected := commandexecutor.Bash().MustRunOneLinerAndGetStdoutAsString("echo -en 'hello world' | git hash-object --stdin", true)
	expected = strings.TrimSpace(expected)

	require.EqualValues(t, expected, GetBlobOjectHashFromString("hello world"))
}
