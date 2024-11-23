package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabReleaseCreateAndDelete(t *testing.T) {
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

				const releaseName string = "test_release"
				const releaseDescription string = "Release description."

				project := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				release := project.MustGetReleaseByName(releaseName)

				for i := 0; i < 2; i++ {
					release.MustDelete(
						&GitlabDeleteReleaseOptions{
							Verbose:                verbose,
							DeleteCorrespondingTag: true,
						},
					)

					assert.False(release.MustExists(verbose))
				}

				release = project.MustCreateReleaseFromLatestCommitInDefaultBranch(
					&GitlabCreateReleaseOptions{
						Name:        releaseName,
						Verbose:     verbose,
						Description: releaseDescription,
					},
				)
				tag := release.MustGetTag()

				assert.True(release.MustExists(verbose))
				assert.True(tag.MustExists(verbose))

				for i := 0; i < 2; i++ {
					release.MustDelete(
						&GitlabDeleteReleaseOptions{
							Verbose:                verbose,
							DeleteCorrespondingTag: true,
						},
					)

					assert.False(release.MustExists(verbose))
					assert.False(tag.MustExists(verbose))
				}
			},
		)
	}
}
