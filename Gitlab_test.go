package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabGetDefaultGitlabCiYamlFileName(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					Gitlab().GetDefaultGitlabCiYamlFileName(),
					".gitlab-ci.yml",
				)
			},
		)
	}
}
