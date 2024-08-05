package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func TestGitlabProjectsProjectDoesNotExist(t *testing.T) {
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

func TestGitlabProjectsGetFileContentAsString(t *testing.T) {
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
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				assert.True(gitlabProject.MustExists(verbose))
			},
		)
	}
}

func TestGitlabProjectsGetProjectIdAndPath(t *testing.T) {
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
			},
		)
	}
}
*/