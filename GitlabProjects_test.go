package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabProjectsProjectDoesNotExist(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		testcase string
	}{
		{"thisProjectDoesNotExist"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.UseUnauthenticatedClient(ctx)
				require.NoError(t, err)
				doesExist, err := gitlab.ProjectByProjectPathExists(ctx, "this/project_does_not_exist")
				require.NoError(t, err)
				require.False(t, doesExist)
			},
		)
	}
}

func TestGitlabProjectsGetProjectIdAndPath(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)
				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				projectId, err := gitlabProject.GetId(ctx)
				require.NoError(t, err)

				gitlabProject2, err := gitlab.GetGitlabProjectById(ctx, projectId)
				require.NoError(t, err)

				exists, err = gitlabProject2.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				require.EqualValues(t, projectPath, mustutils.Must(gitlabProject.GetCachedPath(ctx)))
				require.EqualValues(t, projectPath, mustutils.Must(gitlabProject2.GetCachedPath(ctx)))
				require.EqualValues(t, projectPath, mustutils.Must(gitlabProject.GetPath(ctx)))
				require.EqualValues(t, projectPath, mustutils.Must(gitlabProject2.GetPath(ctx)))
				require.EqualValues(t, "https://"+gitlabFQDN+"/"+projectPath, mustutils.Must(gitlabProject.GetProjectUrl(ctx)))
				require.EqualValues(t, "https://"+gitlabFQDN+"/"+projectPath, mustutils.Must(gitlabProject2.GetProjectUrl(ctx)))
			},
		)
	}
}

func TestGitlabProjectsGetFileContentAsString(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)
				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, "test_group/testproject")
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				fileName := "test.txt"

				for _, content := range []string{"a", "hello", "world"} {
					_, err = gitlabProject.WriteFileContent(
						ctx,
						&GitlabWriteFileOptions{
							Path:          fileName,
							Content:       []byte(content),
							CommitMessage: "commit during automated testing",
						},
					)
					require.NoError(t, err)

					readBack, err := gitlabProject.ReadFileContentAsString(
						ctx,
						&GitlabReadFileOptions{
							Path: fileName,
						},
					)
					require.NoError(t, err)

					require.EqualValues(t, content, readBack)
				}
			},
		)
	}
}
