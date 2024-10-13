package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabProjectSyncFilesToBranch(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
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

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				assert.True(gitlabProject.MustExists(verbose))

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

				assert.EqualValues(
					"hello",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				syncBranch.MustSyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				assert.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
			},
		)
	}
}

func TestGitlabProjectSyncFilesToBranch_notExistingTargetFile(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
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

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				assert.True(gitlabProject.MustExists(verbose))

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

				assert.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				syncBranch.MustSyncFilesToBranch(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				assert.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
			},
		)
	}
}

func TestGitlabProjectSyncFilesToBranch_notExistingTargetFile_usingMR(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Unavailable in Github CI")
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

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				gitlabProject := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				assert.True(gitlabProject.MustExists(verbose))

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

				assert.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				mergeRequest := syncBranch.MustSyncFilesToBranchUsingMergeRequest(&GitlabSyncBranchOptions{
					TargetBranch: defaultBranch,
					Verbose:      verbose,
					PathsToSync:  []string{filePath},
				})

				assert.True(mergeRequest.MustIsOpen())
				assert.False(
					defaultBranch.MustRepositoryFileExists(filePath, verbose),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)

				mergeRequest.MustMerge(
					&GitlabMergeOptions{
						Verbose: verbose,
					},
				)

				assert.True(mergeRequest.MustIsMerged())
				assert.EqualValues(
					"world",
					defaultBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
				assert.EqualValues(
					"world",
					syncBranch.MustReadFileContentAsString(&GitlabReadFileOptions{Path: filePath, Verbose: verbose}),
				)
			},
		)
	}
}
