package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabProjectsProjectDoesNotExist(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"thisProjectDoesNotExist"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustUseUnauthenticatedClient(verbose)
				doesExist := gitlab.MustProjectByProjectPathExists("this/project_does_not_exist", verbose)
				require.False(doesExist)
			},
		)
	}
}

func TestGitlabProjectsGetProjectIdAndPath(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(gitlabProject.MustExists(verbose))

				projectId := gitlabProject.MustGetId()
				gitlabProject2 := gitlab.MustGetGitlabProjectById(projectId, verbose)
				require.True(gitlabProject2.MustExists(verbose))

				require.EqualValues(
					projectPath,
					gitlabProject.MustGetCachedPath(),
				)
				require.EqualValues(
					projectPath,
					gitlabProject2.MustGetCachedPath(),
				)
				require.EqualValues(
					projectPath,
					gitlabProject.MustGetPath(),
				)
				require.EqualValues(
					projectPath,
					gitlabProject2.MustGetPath(),
				)

				require.EqualValues(
					"https://"+gitlabFQDN+"/"+projectPath,
					gitlabProject.MustGetProjectUrl(),
				)
				require.EqualValues(
					"https://"+gitlabFQDN+"/"+projectPath,
					gitlabProject2.MustGetProjectUrl(),
				)
			},
		)
	}
}

func TestGitlabProjectsGetFileContentAsString(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				require.True(gitlabProject.MustExists(verbose))

				fileName := "test.txt"

				for _, content := range []string{"a", "hello", "world"} {
					gitlabProject.MustWriteFileContent(
						&GitlabWriteFileOptions{
							Path:          fileName,
							Content:       []byte(content),
							CommitMessage: "commit during automated testing",
							Verbose:       verbose,
						},
					)

					require.EqualValues(
						content,
						gitlabProject.MustReadFileContentAsString(
							&GitlabReadFileOptions{
								Path:    fileName,
								Verbose: verbose,
							},
						),
					)
				}
			},
		)
	}
}
