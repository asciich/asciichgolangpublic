package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestCommitGetHash(t *testing.T) {
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

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})
				require.NoError(t, err)

				testProject, err := gitlab.CreatePersonalProject(ctx, "testProject")
				require.NoError(t, err)

				branchNames := []string{"testbranch1", "testbranch2"}
				for _, branchName := range branchNames {
					testProject.CreateBranchFromDefaultBranch(ctx, branchName)
				}

				for _, branchName := range branchNames {
					content, err := randomgenerator.GetRandomString(16)
					require.NoError(t, err)

					_, err = testProject.WriteFileContent(
						ctx,
						&GitlabWriteFileOptions{
							Path:          "testfile",
							Content:       []byte(content),
							BranchName:    branchName,
							CommitMessage: "test commit",
						},
					)
					require.NoError(t, err)
				}

				hashes := []string{}
				for _, branchName := range branchNames {
					toAdd, err := testProject.GetLatestCommitHashAsString(ctx, branchName)
					require.NoError(t, err)
					hashes = append(hashes, toAdd)
				}

				require.True(t, slicesutils.ContainsOnlyUniqeStrings(hashes))
				require.True(t, slicesutils.ContainsNoEmptyStrings(hashes))
			},
		)
	}
}

func TestGitlabCommitGetParentCommit(t *testing.T) {
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

				testProject, err := gitlab.CreatePersonalProject(ctx, "testProject")
				require.NoError(t, err)

				latestCommit, err := testProject.GetLatestCommitOfDefaultBranch(ctx)
				require.NoError(t, err)

				_, err = testProject.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          "getParentCommit.txt",
						Content:       []byte("only test content."),
						CommitMessage: "For test case get parent commit",
					},
				)
				require.NoError(t, err)

				latestCommitAfterWrite, err := testProject.GetLatestCommitOfDefaultBranch(ctx)
				require.NoError(t, err)

				hash, err := latestCommit.GetCommitHash()
				require.NoError(t, err)
				hashAfterWrite, err := latestCommitAfterWrite.GetCommitHash()
				require.NoError(t, err)
				require.NotEqualValues(t, hash, hashAfterWrite)

				parents, err := latestCommitAfterWrite.GetParentCommits(ctx)
				require.NoError(t, err)
				require.Len(t, parents, 1)

				latestHash, err := latestCommit.GetCommitHash()
				require.NoError(t, err)
				parentHash, err := parents[0].GetCommitHash()
				require.NoError(t, err)
				require.EqualValues(t, latestHash, parentHash)
			},
		)
	}
}

func TestGitlabCommitGetIsMergeCommit(t *testing.T) {
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

				testProject, err := gitlab.CreatePersonalProject(ctx, "testProject")
				require.NoError(t, err)

				const testFileName string = "isMergeCommit.txt"

				err = testProject.DeleteFileInDefaultBranch(
					ctx,
					testFileName,
					fmt.Sprintf(
						"Ensure %s is absent for TestGitlabCommitGetIsMergeCommit",
						testFileName,
					),
				)
				require.NoError(t, err)

				latestCommit, err := testProject.GetLatestCommitOfDefaultBranch(ctx)
				require.NoError(t, err)

				branch, err := testProject.GetBranchByName("test_merge_commit")
				require.NoError(t, err)
				err = branch.Delete(ctx, &GitlabDeleteBranchOptions{})
				require.NoError(t, err)

				err = branch.CreateFromDefaultBranch(ctx)
				require.NoError(t, err)

				_, err = branch.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          testFileName,
						Content:       []byte("only test content."),
						CommitMessage: "For TestGitlabCommitGetIsMergeCommit",
					},
				)
				require.NoError(t, err)

				latestCommitAfterWrite, err := branch.GetLatestCommit(ctx)
				latestHash, err := latestCommit.GetCommitHash()
				require.NoError(t, err)
				afterWriteHash, err := latestCommitAfterWrite.GetCommitHash()
				require.NoError(t, err)
				require.NotEqualValues(t, latestHash, afterWriteHash)
				require.NoError(t, err)

				mergeRequest, err := branch.CreateMergeRequest(
					ctx,
					&GitlabCreateMergeRequestOptions{
						Title:       "Merge for isMergeCommit test",
						Description: "Merge for isMergeCommit test",
					},
				)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					_, err = mergeRequest.Merge(ctx)
					require.NoError(t, err)

					isMerged, err := mergeRequest.IsMerged(ctx)
					require.NoError(t, err)
					require.True(t, isMerged)
				}

				commitAfterMerge, err := testProject.GetLatestCommitOfDefaultBranch(ctx)
				require.NoError(t, err)

				isMergeCommit, err := commitAfterMerge.IsMergeCommit(ctx)
				require.NoError(t, err)
				require.True(t, isMergeCommit)

				isMergeCommit, err = latestCommitAfterWrite.IsMergeCommit(ctx)
				require.NoError(t, err)
				require.False(t, isMergeCommit)

				isParentOf, err := latestCommit.IsParentCommitOf(ctx, commitAfterMerge)
				require.NoError(t, err)
				require.True(t, isParentOf)

				isParentOf, err = latestCommit.IsParentCommitOf(ctx, latestCommitAfterWrite)
				require.NoError(t, err)
				require.True(t, isParentOf)
			},
		)
	}
}
