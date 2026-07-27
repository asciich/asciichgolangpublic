package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitlabGetDefaultGitlabCiYamlFileName(t *testing.T) {
	require.EqualValues(t, Gitlab().GetDefaultGitlabCiYamlFileName(), ".gitlab-ci.yml")
}
