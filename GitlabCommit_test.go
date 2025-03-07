package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/randomgenerator"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				branchNames := []string{"testbranch1", "testbranch2"}
				for _, branchName := range branchNames {
					testProject.CreateBranchFromDefaultBranch(branchName, verbose)
				}

				for _, branchName := range branchNames {
					content := randomgenerator.MustGetRandomString(16)
					testProject.MustWriteFileContent(
						&GitlabWriteFileOptions{
							Path:          "testfile",
							Content:       []byte(content),
							Verbose:       verbose,
							BranchName:    branchName,
							CommitMessage: "test commit",
						},
					)
				}

				hashes := []string{}
				for _, branchName := range branchNames {
					hashes = append(hashes, testProject.MustGetLatestCommitHashAsString(branchName, verbose))
				}

				require.True(slicesutils.ContainsOnlyUniqeStrings(hashes))
				require.True(slicesutils.ContainsNoEmptyStrings(hashes))
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
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				latestCommit := testProject.MustGetLatestCommitOfDefaultBranch(verbose)
				testProject.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          "getParentCommit.txt",
						Content:       []byte("only test content."),
						CommitMessage: "For test case get parent commit",
						Verbose:       verbose,
					},
				)
				latestCommitAfterWrite := testProject.MustGetLatestCommitOfDefaultBranch(verbose)
				require.NotEqualValues(
					latestCommit.MustGetCommitHash(),
					latestCommitAfterWrite.MustGetCommitHash(),
				)

				parents := latestCommitAfterWrite.MustGetParentCommits(verbose)
				require.Len(parents, 1)

				require.EqualValues(
					latestCommit.MustGetCommitHash(),
					parents[0].MustGetCommitHash(),
				)
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
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)
				const testFileName string = "isMergeCommit.txt"

				testProject.MustDeleteFileInDefaultBranch(
					testFileName,
					fmt.Sprintf(
						"Ensure %s is absent for TestGitlabCommitGetIsMergeCommit",
						testFileName,
					),
					verbose,
				)

				latestCommit := testProject.MustGetLatestCommitOfDefaultBranch(verbose)

				branch := testProject.MustGetBranchByName("test_merge_commit")
				branch.MustDelete(&GitlabDeleteBranchOptions{Verbose: verbose})
				branch.MustCreateFromDefaultBranch(verbose)

				branch.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          testFileName,
						Content:       []byte("only test content."),
						CommitMessage: "For TestGitlabCommitGetIsMergeCommit",
						Verbose:       verbose,
					},
				)
				latestCommitAfterWrite := branch.MustGetLatestCommit(verbose)
				require.NotEqualValues(
					latestCommit.MustGetCommitHash(),
					latestCommitAfterWrite.MustGetCommitHash(),
				)

				mergeRequest := branch.MustCreateMergeRequest(
					&GitlabCreateMergeRequestOptions{
						Title:       "Merge for isMergeCommit test",
						Description: "Merge for isMergeCommit test",
						Verbose:     verbose,
					},
				)

				for i := 0; i < 2; i++ {
					mergeRequest.MustMerge(
						&GitlabMergeOptions{
							Verbose: verbose,
						},
					)
					require.True(mergeRequest.MustIsMerged())
				}

				commitAfterMerge := testProject.MustGetLatestCommitOfDefaultBranch(verbose)

				require.True(commitAfterMerge.MustIsMergeCommit(verbose))
				require.False(latestCommitAfterWrite.MustIsMergeCommit(verbose))

				require.True(latestCommit.MustIsParentCommitOf(commitAfterMerge, verbose))
				require.True(latestCommit.MustIsParentCommitOf(latestCommitAfterWrite, verbose))
			},
		)
	}
}
