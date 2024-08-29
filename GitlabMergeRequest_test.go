package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRequestCreateAndClose(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in github")
		return
	}

	tests := []struct {
		mergeRequestTitle string
	}{
		{"merge request title 1"},
		{"merge request title 2"},
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

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				const testBranchName string = "mr_test_branch"

				testProject.MustDeleteBranch(testBranchName, &GitlabDeleteBranchOptions{
					Verbose: verbose,
				})
				branch := testProject.MustCreateBranchFromDefaultBranch(testBranchName, verbose)

				for i := 0; i < 2; i++ {
					mergeRequest := branch.MustCreateMergeRequest(
						&GitlabCreateMergeRequestOptions{
							Title:   tt.mergeRequestTitle,
							Verbose: verbose,
						},
					)
					assert.True(mergeRequest.MustIsOpen())
					assert.EqualValues(
						testProject.MustGetDefaultBranchName(),
						mergeRequest.MustGetTargetBranchName(),
					)
					assert.EqualValues(
						testBranchName,
						mergeRequest.MustGetSourceBranchName(),
					)
				}

				mergeRequest := testProject.MustGetOpenMergeRequestByTitle(tt.mergeRequestTitle, verbose)
				for i := 0; i < 2; i++ {
					mergeRequest.MustClose("closed for testing", verbose)
					assert.True(mergeRequest.MustIsClosed())
				}
			},
		)
	}
}

func TestMergeRequestCreateAndClose_withLabels(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in github")
		return
	}

	tests := []struct {
		mergeRequestTitle string
		labels            []string
		description       string
	}{
		{
			"merge request title 1",
			[]string{"label_a"},
			"MR description",
		},
		{
			"merge request title 2",
			[]string{"label_a", "label_b"},
			"MR description\nmultiline",
		},
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

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				const testBranchName string = "mr_test_branch"

				testProject.MustDeleteBranch(testBranchName, &GitlabDeleteBranchOptions{
					Verbose: verbose,
				})
				branch := testProject.MustCreateBranchFromDefaultBranch(testBranchName, verbose)

				for i := 0; i < 2; i++ {
					mergeRequest := testProject.MustCreateMergeRequest(
						&GitlabCreateMergeRequestOptions{
							SourceBranchName: branch.MustGetName(),
							Title:            tt.mergeRequestTitle,
							Labels:           tt.labels,
							Description:      tt.description,
							Verbose:          verbose,
						},
					)
					assert.True(mergeRequest.MustIsOpen())
					assert.EqualValues(
						testProject.MustGetDefaultBranchName(),
						mergeRequest.MustGetTargetBranchName(),
					)
					assert.EqualValues(
						testBranchName,
						mergeRequest.MustGetSourceBranchName(),
					)

					assert.EqualValues(
						tt.labels,
						mergeRequest.MustGetLabels(),
					)

					assert.EqualValues(
						tt.description,
						mergeRequest.MustGetDescription(),
					)
				}

				mergeRequest := testProject.MustGetOpenMergeRequestByTitle(tt.mergeRequestTitle, verbose)
				for i := 0; i < 2; i++ {
					mergeRequest.MustClose("closed for testing", verbose)
					assert.True(mergeRequest.MustIsClosed())
				}
			},
		)
	}
}
