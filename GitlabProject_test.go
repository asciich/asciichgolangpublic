package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabProjectSyncFilesToBranch(t *testing.T) {
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
				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(t, gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				_, err := defaultBranch.WriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("hello"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				content, err := defaultBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "hello", content)

				content, err = syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				err = syncBranch.SyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})
				require.NoError(t, err)

				content, err = defaultBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)
			},
		)
	}
}

func TestGitlabProjectSyncFilesToBranch_notExistingTargetFile(t *testing.T) {
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
				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(t, gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				err := defaultBranch.DeleteRepositoryFile(
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
					verbose,
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				exists, err := defaultBranch.RepositoryFileExists(filePath, verbose)
				require.NoError(t, err)
				require.False(t, exists)

				content, err := syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				err = syncBranch.SyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})
				require.NoError(t, err)

				content, err = defaultBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)
			},
		)
	}
}

func TestGitlabProjectSyncFilesToBranch_notExistingTargetFile_usingMR(t *testing.T) {
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
				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(t, gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				err := defaultBranch.DeleteRepositoryFile(
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
					verbose,
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				exits, err := defaultBranch.RepositoryFileExists(filePath, verbose)
				require.NoError(t, err)
				require.False(t, exits)

				content, err := syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				mergeRequest, err := syncBranch.SyncFilesToBranchUsingMergeRequest(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})
				require.NoError(t, err)

				require.True(t, mergeRequest.MustIsOpen())

				exists, err := defaultBranch.RepositoryFileExists(filePath, verbose)
				require.NoError(t, err)
				require.False(t, exists)

				content, err = syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.EqualValues(t, "world", content)

				mergeRequest.MustMerge(
					&GitlabMergeOptions{
						Verbose: verbose,
					},
				)

				require.True(t, mergeRequest.MustIsMerged())

				content, err = defaultBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)
			},
		)
	}
}
