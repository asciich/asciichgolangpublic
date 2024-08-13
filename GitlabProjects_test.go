package asciichgolangpublic


import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabProjectsProjectDoesNotExist(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
		return
	}
	
	tests := []struct {
		testcase string
	}{
		{"thisProjectDoesNotExist"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustUseUnauthenticatedClient(verbose)
				doesExist := gitlab.MustProjectByProjectPathExists("this/project_does_not_exist", verbose)
				assert.False(doesExist)
			},
		)
	}
}

func TestGitlabProjectsGetProjectIdAndPath(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
		return
	}
	
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

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				assert.True(gitlabProject.MustExists(verbose))

				projectId := gitlabProject.MustGetId()
				gitlabProject2 := gitlab.MustGetGitlabProjectById(projectId, verbose)
				assert.True(gitlabProject2.MustExists(verbose))

				assert.EqualValues(
					projectPath,
					gitlabProject.MustGetCachedPath(),
				)
				assert.EqualValues(
					projectPath,
					gitlabProject2.MustGetCachedPath(),
				)
				assert.EqualValues(
					projectPath,
					gitlabProject.MustGetPath(),
				)
				assert.EqualValues(
					projectPath,
					gitlabProject2.MustGetPath(),
				)

				assert.EqualValues(
					"https://"+gitlabFQDN+"/"+projectPath,
					gitlabProject.MustGetProjectUrl(),
				)
				assert.EqualValues(
					"https://"+gitlabFQDN+"/"+projectPath,
					gitlabProject2.MustGetProjectUrl(),
				)
			},
		)
	}
}

func TestGitlabProjectsGetFileContentAsString(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
		return
	}
	
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

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				assert.True(gitlabProject.MustExists(verbose))

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

					assert.EqualValues(
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
