package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitGetHash(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in github")
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

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				branchNames := []string{"testbranch1", "testbranch2"}
				for _, branchName := range branchNames {
					testProject.CreateBranchFromDefaultBranch(branchName, verbose)
				}

				for _, branchName := range branchNames {
					content := RandomGenerator().MustGetRandomString(16)
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

				assert.True(Slices().ContainsOnlyUniqeStrings(hashes))
				assert.True(Slices().ContainsNoEmptyStrings(hashes))
			},
		)
	}
}

func TestGitlabCommitGetParentCommit(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in github")
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

				testProject := gitlab.MustCreatePersonalProject("testProject", verbose)

				latestCommit := testProject.MustGetLatestCommitOfDefaultBranch(verbose)
				testProject.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:    "getParentCommit.txt",
						Content: []byte("only test content."),
						CommitMessage: "For test case get parent commit",
						Verbose: verbose,
					},
				)
				latestCommitAfterWrite := testProject.MustGetLatestCommitOfDefaultBranch(verbose)
				assert.NotEqualValues(
					latestCommit.MustGetCommitHash(),
					latestCommitAfterWrite.MustGetCommitHash(),
				)

				parents := latestCommitAfterWrite.MustGetParentCommits(verbose)
				assert.Len(parents, 1)

				assert.EqualValues(
					latestCommit.MustGetCommitHash(),
					parents[0].MustGetCommitHash(),
				)
			},
		)
	}
}
