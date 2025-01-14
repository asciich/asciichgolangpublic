package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
)

func TestGitlabProjectBranchCreateAndDelete(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		branchName string
	}{
		{"testbranch"},
		{"testbranch2"},
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

				branch := gitlabProject.MustGetBranchByName(tt.branchName)

				branch.MustDelete(&GitlabDeleteBranchOptions{
					SkipWaitForDeletion: false,
					Verbose:             verbose,
				})
				assert.False(branch.MustExists())

				for i := 0; i < 2; i++ {
					branch.CreateFromDefaultBranch(verbose)
					assert.True(branch.MustExists())
				}

				for i := 0; i < 2; i++ {
					branch.MustDelete(&GitlabDeleteBranchOptions{
						SkipWaitForDeletion: false,
						Verbose:             verbose,
					})
					assert.False(branch.MustExists())
				}
			},
		)
	}
}
