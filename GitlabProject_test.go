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
				ctx := getCtx()

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				defaultBranch, err := gitlabProject.GetDefaultBranch(ctx)
				require.NoError(t, err)
				syncBranch, err := gitlabProject.CreateBranchFromDefaultBranch(ctx, "test_sync")
				require.NoError(t, err)

				const filePath = "abc.txt"

				_, err = defaultBranch.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("hello"),
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				content, err := defaultBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "hello", content)

				content, err = syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				err = syncBranch.SyncFilesToBranch(
					ctx,
					&GitlabSyncBranchOptions{
						TargetBranch: defaultBranch,
						PathsToSync:  []string{filePath},
					},
				)
				require.NoError(t, err)

				content, err = defaultBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
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
				ctx := getCtx()

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)
				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)
				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				defaultBranch, err := gitlabProject.GetDefaultBranch(ctx)
				require.NoError(t, err)
				syncBranch, err := gitlabProject.CreateBranchFromDefaultBranch(ctx, "test_sync")
				require.NoError(t, err)

				const filePath = "abc.txt"

				err = defaultBranch.DeleteRepositoryFile(
					ctx,
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				exists, err = defaultBranch.RepositoryFileExists(ctx, filePath)
				require.NoError(t, err)
				require.False(t, exists)

				content, err := syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				err = syncBranch.SyncFilesToBranch(
					ctx,
					&GitlabSyncBranchOptions{
						TargetBranch: defaultBranch,
						PathsToSync:  []string{filePath},
					})
				require.NoError(t, err)

				content, err = defaultBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
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
				ctx := getCtx()

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)
				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				defaultBranch, err := gitlabProject.GetDefaultBranch(ctx)
				require.NoError(t, err)

				syncBranch, err := gitlabProject.CreateBranchFromDefaultBranch(ctx, "test_sync")
				require.NoError(t, err)

				const filePath = "abc.txt"

				err = defaultBranch.DeleteRepositoryFile(
					ctx,
					filePath,
					"Cleanup for testing TestGitlabProjectSyncFilesToBranch_notExistingTargetFile",
				)
				require.NoError(t, err)

				_, err = syncBranch.WriteFileContent(
					ctx,
					&GitlabWriteFileOptions{
						Path:          filePath,
						Content:       []byte("world"),
						CommitMessage: "TestGitlabProjectSyncFilesToBranch",
					},
				)
				require.NoError(t, err)

				exits, err := defaultBranch.RepositoryFileExists(ctx, filePath)
				require.NoError(t, err)
				require.False(t, exits)

				content, err := syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				mergeRequest, err := syncBranch.SyncFilesToBranchUsingMergeRequest(
					ctx,
					&GitlabSyncBranchOptions{
						TargetBranch: defaultBranch,
						PathsToSync:  []string{filePath},
					},
				)
				require.NoError(t, err)

				isOpen, err := mergeRequest.IsOpen(ctx)
				require.NoError(t, err)
				require.True(t, isOpen)

				exists, err = defaultBranch.RepositoryFileExists(ctx, filePath)
				require.NoError(t, err)
				require.False(t, exists)

				content, err = syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				_, err = mergeRequest.Merge(ctx)
				require.NoError(t, err)

				isMerged, err := mergeRequest.IsMerged(ctx)
				require.NoError(t, err)
				require.True(t, isMerged)

				content, err = defaultBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)

				content, err = syncBranch.ReadFileContentAsString(ctx, &GitlabReadFileOptions{Path: filePath})
				require.NoError(t, err)
				require.EqualValues(t, "world", content)
			},
		)
	}
}
