package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabCiYamlFileGetInclude(t *testing.T) {
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

				gitlabCiYamlFile.MustWriteString("---\n", verbose)
				gitlabCiYamlFile.MustAppendString("include:\n", verbose)
				gitlabCiYamlFile.MustAppendString("  - project: a\n", verbose)
				gitlabCiYamlFile.MustAppendString("    ref: b\n", verbose)
				gitlabCiYamlFile.MustAppendString("    file: c.yaml\n", verbose)

				includes := gitlabCiYamlFile.MustGetIncludes(verbose)
				assert.Len(includes, 1)

				include := includes[0]
				assert.EqualValues("a", include.MustGetProject())
				assert.EqualValues("b", include.MustGetRef())
				assert.EqualValues("c.yaml", include.MustGetFile())
			},
		)
	}
}

func TestGitlabCiYamlFileGetInclude2(t *testing.T) {
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

				gitlabCiYamlFile.MustWriteString("---\n", verbose)
				gitlabCiYamlFile.MustAppendString("include:\n", verbose)
				gitlabCiYamlFile.MustAppendString("  - project: a\n", verbose)
				gitlabCiYamlFile.MustAppendString("    ref: b\n", verbose)
				gitlabCiYamlFile.MustAppendString("    file:\n", verbose)
				gitlabCiYamlFile.MustAppendString("     - c.yaml\n", verbose)

				includes := gitlabCiYamlFile.MustGetIncludes(verbose)
				assert.Len(includes, 1)

				include := includes[0]

				assert.EqualValues("a", include.MustGetProject())
				assert.EqualValues("b", include.MustGetRef())
				assert.EqualValues("c.yaml", include.MustGetFile())
			},
		)
	}
}

// include rules are currently ignored during parsing.
func TestGitlabCiYamlFileGetIncludeIgnoreRules(t *testing.T) {
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

				gitlabCiYamlFile.MustWriteString("---\n", verbose)
				gitlabCiYamlFile.MustAppendString("include:\n", verbose)
				gitlabCiYamlFile.MustAppendString("  - project: a\n", verbose)
				gitlabCiYamlFile.MustAppendString("    ref: b\n", verbose)
				gitlabCiYamlFile.MustAppendString("    file:\n", verbose)
				gitlabCiYamlFile.MustAppendString("     - c.yaml\n", verbose)
				gitlabCiYamlFile.MustAppendString("    rules:\n", verbose)
				gitlabCiYamlFile.MustAppendString("     - if: $CI_PIPELINE_SOURCE != \"pipeline\"\n", verbose)

				includes := gitlabCiYamlFile.MustGetIncludes(verbose)
				assert.Len(includes, 1)

				include := includes[0]

				assert.EqualValues("a", include.MustGetProject())
				assert.EqualValues("b", include.MustGetRef())
				assert.EqualValues("c.yaml", include.MustGetFile())
			},
		)
	}
}

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
