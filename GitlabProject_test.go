package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				defaultBranch.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("hello"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)

				syncBranch.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)

				require.EqualValues(
					"hello",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				syncBranch.MustSyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				require.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
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
				require := require.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				defaultBranch.MustDeleteRepositoryFile(
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
					verbose,
				)

				syncBranch.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)

				require.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				syncBranch.MustSyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				require.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
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
				require := require.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				require.True(gitlabProject.MustExists(verbose))

				defaultBranch := gitlabProject.MustGetDefaultBranch()
				syncBranch := gitlabProject.MustCreateBranchFromDefaultBranch("test_sync", verbose)

				const filePath = "abc.txt"

				defaultBranch.MustDeleteRepositoryFile(
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
					verbose,
				)

				syncBranch.MustWriteFileContent(
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						Verbose:       verbose,
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)

				require.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				mergeRequest := syncBranch.MustSyncFilesToBranchUsingMergeRequest(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				require.True(mergeRequest.MustIsOpen())
				require.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				mergeRequest.MustMerge(
					&GitlabMergeOptions{
						Verbose: verbose,
					},
				)

				require.True(mergeRequest.MustIsMerged())
				require.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				require.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
			},
		)
	}
}
