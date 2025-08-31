package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabProjectsGetFileList(t *testing.T) {
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

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, "test_group/testproject")
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.True(t, exists)

				branchName, err := gitlabProject.GetDefaultBranchName(ctx)
				require.NoError(t, err)

				err = gitlabProject.DeleteAllRepositoryFiles(ctx, branchName)
				require.NoError(t, err)

				hasRepoFile, err := gitlabProject.HasNoRepositoryFiles(ctx, branchName)
				require.NoError(t, err)
				require.True(t, hasRepoFile)

				_, err = gitlabProject.CreateEmptyFile(ctx, "a.txt", branchName)
				require.NoError(t, err)
				_, err = gitlabProject.CreateEmptyFile(ctx, "b.txt", branchName)
				require.NoError(t, err)
				_, err = gitlabProject.CreateEmptyFile(ctx, "aa/a.txt", branchName)
				require.NoError(t, err)
				_, err = gitlabProject.CreateEmptyFile(ctx, "aa/b.txt", branchName)
				require.NoError(t, err)
				_, err = gitlabProject.CreateEmptyFile(ctx, "evenMore/aa/b.txt", branchName)
				require.NoError(t, err)
				_, err = gitlabProject.CreateEmptyFile(ctx, "evenMore/aa/a.txt", branchName)
				require.NoError(t, err)

				fileList, err := gitlabProject.GetFilesNames(ctx, branchName)
				require.NoError(t, err)
				expectedFileList := []string{
					"a.txt",
					"aa/a.txt",
					"aa/b.txt",
					"b.txt",
					"evenMore/aa/a.txt",
					"evenMore/aa/b.txt",
				}
				require.EqualValues(t, expectedFileList, fileList)

				directoryList, err := gitlabProject.GetDirectoryNames(ctx, branchName)
				exepctedDirectoryList := []string{
					"aa",
					"evenMore",
					"evenMore/aa",
				}
				require.EqualValues(t, exepctedDirectoryList, directoryList)
			},
		)
	}
}

func TestGitlabProjectsGetFileList_pagination(t *testing.T) {
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

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, "test_group/testproject")

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				branchName, err := gitlabProject.GetDefaultBranchName(ctx)
				require.NoError(t, err)

				err = gitlabProject.DeleteAllRepositoryFiles(ctx, branchName)
				require.NoError(t, err)

				hasNoRepositoryFiles, err := gitlabProject.HasNoRepositoryFiles(ctx, branchName)
				require.NoError(t, err)
				require.True(t, hasNoRepositoryFiles)

				expectedFileList := []string{}
				for i := 0; i < 21; i++ {
					expectedFileList = append(
						expectedFileList,
						fmt.Sprintf("asd%03d.txt", i),
					)
				}

				for _, toCreate := range expectedFileList {
					_, err = gitlabProject.CreateEmptyFile(ctx, toCreate, branchName)
					require.NoError(t, err)
				}

				fileList, err := gitlabProject.GetFilesNames(ctx, branchName)
				require.NoError(t, err)
				require.EqualValues(t, expectedFileList, fileList)
			},
		)
	}
}
