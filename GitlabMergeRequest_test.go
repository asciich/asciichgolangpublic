package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestMergeRequestCreateAndClose(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		mergeRequestTitle string
	}{
		{"merge request title 1"},
		{"merge request title 2"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
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
					mergeRequest, err := branch.CreateMergeRequest(
						&GitlabCreateMergeRequestOptions{
							Title:   tt.mergeRequestTitle,
							Verbose: verbose,
						},
					)
					require.NoError(t, err)

					require.True(t, mergeRequest.MustIsOpen())
					require.EqualValues(t, testProject.MustGetDefaultBranchName(), mergeRequest.MustGetTargetBranchName())
					require.EqualValues(t, testBranchName, mergeRequest.MustGetSourceBranchName())
				}

				mergeRequest := testProject.MustGetOpenMergeRequestByTitle(tt.mergeRequestTitle, verbose)
				for i := 0; i < 2; i++ {
					mergeRequest.MustClose("closed for testing", verbose)
					require.True(t, mergeRequest.MustIsClosed())
				}
			},
		)
	}
}

func TestMergeRequestCreateAndClose_withLabels(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
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
					branchName, err := branch.GetName()
					require.NoError(t, err)

					mergeRequest := testProject.MustCreateMergeRequest(
						&GitlabCreateMergeRequestOptions{
							SourceBranchName: branchName,
							Title:            tt.mergeRequestTitle,
							Labels:           tt.labels,
							Description:      tt.description,
							Verbose:          verbose,
						},
					)
					require.True(t, mergeRequest.MustIsOpen())
					require.EqualValues(t, testProject.MustGetDefaultBranchName(), mergeRequest.MustGetTargetBranchName())
					require.EqualValues(t, testBranchName, mergeRequest.MustGetSourceBranchName())

					require.EqualValues(t, tt.labels, mergeRequest.MustGetLabels())

					require.EqualValues(t, tt.description, mergeRequest.MustGetDescription())
				}

				branchName, err := branch.GetName()
				require.NoError(t, err)

				mergeRequest := testProject.MustGetOpenMergeRequestBySourceAndTargetBranch(
					branchName,
					testProject.MustGetDefaultBranchName(),
					verbose,
				)
				for i := 0; i < 2; i++ {
					mergeRequest.MustClose("closed for testing", verbose)
					require.True(t, mergeRequest.MustIsClosed())
				}
			},
		)
	}
}
