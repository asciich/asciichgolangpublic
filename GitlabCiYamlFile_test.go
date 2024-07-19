package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabCiYamlFileAddIncludes(t *testing.T) {
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

				const verbose bool = true

				emptyFilePath := TemporaryFiles().MustCreateEmptyTemporaryFileAndGetPath(verbose)
				gitlabCiYamlFile := MustGetGitlabCiYamlFileByPath(emptyFilePath)

				for i := 0; i < 3; i++ {
					gitlabCiYamlFile.MustAddInclude(
						&GitlabCiYamlInclude{
							Project: "abc",
							File:    "test.yml",
							Ref:     "v1234",
						},
						verbose,
					)

					includes := gitlabCiYamlFile.MustGetIncludes(verbose)
					assert.Len(includes, 1)
				}

				for i := 0; i < 3; i++ {
					gitlabCiYamlFile.MustAddInclude(
						&GitlabCiYamlInclude{
							Project: "abc_other",
							File:    "test2.yml",
							Ref:     "v12345",
						},
						verbose,
					)

					includes := gitlabCiYamlFile.MustGetIncludes(verbose)
					assert.Len(includes, 2)
				}
			},
		)
	}
}
