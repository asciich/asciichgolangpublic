package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabCiYamlFileGetInclude(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				toWrite := "---\n"
				toWrite += "include:\n"
				toWrite += "  - project: a\n"
				toWrite += "    ref: b\n"
				toWrite += "    file: c.yaml\n"
				err = gitlabCiYamlFile.WriteString(toWrite, verbose)
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(verbose)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]
				require.EqualValues(t, "a", include.MustGetProject())
				require.EqualValues(t, "b", include.MustGetRef())
				require.EqualValues(t, "c.yaml", include.MustGetFile())
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				toWrite := "---\n"
				toWrite += "include:\n"
				toWrite += "  - project: a\n"
				toWrite += "    ref: b\n"
				toWrite += "    file:\n"
				toWrite += "     - c.yaml\n"
				err = gitlabCiYamlFile.WriteString(toWrite, verbose)
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(verbose)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]

				require.EqualValues(t, "a", include.MustGetProject())
				require.EqualValues(t, "b", include.MustGetRef())
				require.EqualValues(t, "c.yaml", include.MustGetFile())
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				toWrite := "---\n"
				toWrite += "include:\n"
				toWrite += "  - project: a\n"
				toWrite += "    ref: b\n"
				toWrite += "    file:\n"
				toWrite += "     - c.yaml\n"
				toWrite += "    rules:\n"
				toWrite += "     - if: $CI_PIPELINE_SOURCE != \"pipeline\"\n"
				err = gitlabCiYamlFile.WriteString(toWrite, verbose)
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(verbose)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]

				require.EqualValues(t, "a", include.MustGetProject())
				require.EqualValues(t, "b", include.MustGetRef())
				require.EqualValues(t, "c.yaml", include.MustGetFile())
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				for i := 0; i < 3; i++ {
					err = gitlabCiYamlFile.AddInclude(
						&GitlabCiYamlInclude{
							Project: "abc",
							File:    "test.yml",
							Ref:     "v1234",
						},
						verbose,
					)
					require.NoError(t, err)

					includes, err := gitlabCiYamlFile.GetIncludes(verbose)
					require.NoError(t, err)
					require.Len(t, includes, 1)
				}

				for i := 0; i < 3; i++ {
					err = gitlabCiYamlFile.AddInclude(
						&GitlabCiYamlInclude{
							Project: "abc_other",
							File:    "test2.yml",
							Ref:     "v12345",
						},
						verbose,
					)
					require.NoError(t, err)

					includes, err := gitlabCiYamlFile.GetIncludes(verbose)
					require.NoError(t, err)
					require.Len(t, includes, 2)
				}
			},
		)
	}
}
