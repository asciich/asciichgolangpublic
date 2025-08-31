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
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)
				
				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				testProject, err := gitlab.CreatePersonalProject(ctx, "testProject")
				require.NoError(t, err)

				const testBranchName string = "mr_test_branch"

				err = testProject.DeleteBranch(ctx, testBranchName, &GitlabDeleteBranchOptions{})
				require.NoError(t, err)

				branch, err := testProject.CreateBranchFromDefaultBranch(ctx, testBranchName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					mergeRequest, err := branch.CreateMergeRequest(
						ctx,
						&GitlabCreateMergeRequestOptions{
							Title: tt.mergeRequestTitle,
						},
					)
					require.NoError(t, err)

					isOpen, err := mergeRequest.IsOpen(ctx)
					require.NoError(t, err)
					require.True(t, isOpen)

					defaultBranchName, err := testProject.GetDefaultBranchName(ctx)
					require.NoError(t, err)
					targetBranchName, err := mergeRequest.GetTargetBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, defaultBranchName, targetBranchName)

					sourceBranchName, err := mergeRequest.GetSourceBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, testBranchName, sourceBranchName)
				}

				mergeRequest, err := testProject.GetOpenMergeRequestByTitle(ctx, tt.mergeRequestTitle)
				require.NoError(t, err)
				for i := 0; i < 2; i++ {
					err := mergeRequest.Close(ctx, "closed for testing")
					require.NoError(t, err)

					isClosed, err := mergeRequest.IsClosed(ctx)
					require.NoError(t, err)
					require.True(t, isClosed)
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
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				testProject, err := gitlab.CreatePersonalProject(ctx, "testProject")
				require.NoError(t, err)

				const testBranchName string = "mr_test_branch"

				err = testProject.DeleteBranch(ctx, testBranchName, &GitlabDeleteBranchOptions{})
				require.NoError(t, err)

				branch, err := testProject.CreateBranchFromDefaultBranch(ctx, testBranchName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					branchName, err := branch.GetName()
					require.NoError(t, err)

					mergeRequest, err := testProject.CreateMergeRequest(
						ctx,
						&GitlabCreateMergeRequestOptions{
							SourceBranchName: branchName,
							Title:            tt.mergeRequestTitle,
							Labels:           tt.labels,
							Description:      tt.description,
						},
					)
					require.NoError(t, err)

					isOpen, err := mergeRequest.IsOpen(ctx)
					require.NoError(t, err)
					require.True(t, isOpen)

					defaultBranchName, err := testProject.GetDefaultBranchName(ctx)
					require.NoError(t, err)
					targetBranchName, err := mergeRequest.GetTargetBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, defaultBranchName, targetBranchName)

					sourceBranch, err := mergeRequest.GetSourceBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, testBranchName, sourceBranch)

					labels, err := mergeRequest.GetLabels(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.labels, labels)

					description, err := mergeRequest.GetDescription(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.description, description)
				}

				branchName, err := branch.GetName()
				require.NoError(t, err)

				defaultBranchName, err := testProject.GetDefaultBranchName(ctx)
				require.NoError(t, err)
				mergeRequest, err := testProject.GetOpenMergeRequestBySourceAndTargetBranch(ctx, branchName, defaultBranchName)
				require.NoError(t, err)
				for i := 0; i < 2; i++ {
					err = mergeRequest.Close(ctx, "closed for testing")
					require.NoError(t, err)

					isClosed, err := mergeRequest.IsClosed(ctx)
					require.NoError(t, err)
					require.True(t, isClosed)
				}
			},
		)
	}
}
