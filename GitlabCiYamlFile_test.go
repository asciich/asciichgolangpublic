package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
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
				ctx := getCtx()

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				toWrite := "---\n"
				toWrite += "include:\n"
				toWrite += "  - project: a\n"
				toWrite += "    ref: b\n"
				toWrite += "    file: c.yaml\n"
				err = gitlabCiYamlFile.WriteString(ctx, toWrite, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(ctx)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]

				project, err := include.GetProject()
				require.NoError(t, err)
				require.EqualValues(t, "a", project)

				ref, err := include.GetRef()
				require.NoError(t, err)
				require.EqualValues(t, "b", ref)

				file, err := include.GetFile()
				require.NoError(t, err)
				require.EqualValues(t, "c.yaml", file)
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
				ctx := getCtx()

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				toWrite := "---\n"
				toWrite += "include:\n"
				toWrite += "  - project: a\n"
				toWrite += "    ref: b\n"
				toWrite += "    file:\n"
				toWrite += "     - c.yaml\n"
				err = gitlabCiYamlFile.WriteString(ctx, toWrite, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(ctx)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]

				project, err := include.GetProject()
				require.NoError(t, err)
				require.EqualValues(t, "a", project)

				ref, err := include.GetRef()
				require.NoError(t, err)
				require.EqualValues(t, "b", ref)

				file, err := include.GetFile()
				require.NoError(t, err)
				require.EqualValues(t, "c.yaml", file)
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
				ctx := getCtx()

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
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
				err = gitlabCiYamlFile.WriteString(ctx, toWrite, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				includes, err := gitlabCiYamlFile.GetIncludes(ctx)
				require.NoError(t, err)
				require.Len(t, includes, 1)

				include := includes[0]

				project, err := include.GetProject()
				require.NoError(t, err)
				require.EqualValues(t, "a", project)

				ref, err := include.GetRef()
				require.NoError(t, err)
				require.EqualValues(t, "b", ref)

				file, err := include.GetFile()
				require.NoError(t, err)
				require.EqualValues(t, "c.yaml", file)
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
				ctx := getCtx()

				emptyFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)
				gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(emptyFilePath)
				require.NoError(t, err)

				for i := 0; i < 3; i++ {
					err = gitlabCiYamlFile.AddInclude(
						ctx,
						&GitlabCiYamlInclude{
							Project: "abc",
							File:    "test.yml",
							Ref:     "v1234",
						},
					)
					require.NoError(t, err)

					includes, err := gitlabCiYamlFile.GetIncludes(ctx)
					require.NoError(t, err)
					require.Len(t, includes, 1)
				}

				for i := 0; i < 3; i++ {
					err = gitlabCiYamlFile.AddInclude(
						ctx,
						&GitlabCiYamlInclude{
							Project: "abc_other",
							File:    "test2.yml",
							Ref:     "v12345",
						},
					)
					require.NoError(t, err)

					includes, err := gitlabCiYamlFile.GetIncludes(ctx)
					require.NoError(t, err)
					require.Len(t, includes, 2)
				}
			},
		)
	}
}
