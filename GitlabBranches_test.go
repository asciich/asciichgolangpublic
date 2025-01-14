package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
)

func TestGitlabProjectBranches_pagination(t *testing.T) {
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

				branchName := gitlabProject.MustGetDefaultBranchName()

				gitlabProject.MustDeleteAllBranchesExceptDefaultBranch(verbose)

				branchList := gitlabProject.MustGetBranchNames(verbose)
				assert.EqualValues(
					[]string{branchName},
					branchList,
				)

				expectedBranchesList := []string{}
				for i := 0; i < 21; i++ {
					expectedBranchesList = append(
						expectedBranchesList,
						fmt.Sprintf("testbranch_%03d", i),
					)
				}

				for _, toCreate := range expectedBranchesList {
					gitlabProject.MustCreateBranchFromDefaultBranch(toCreate, verbose)
				}

				branchList = gitlabProject.MustGetBranchNames(verbose)
				assert.EqualValues(branchList, branchList)
			},
		)
	}
}
