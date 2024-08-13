package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabProjectsGetFileList(t *testing.T) {
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

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				assert.True(gitlabProject.MustExists(verbose))

				branchName := gitlabProject.MustGetDefaultBranchName()

				gitlabProject.MustDeleteAllRepositoryFiles(branchName, verbose)
				assert.True(gitlabProject.MustHasNoRepositoryFiles(branchName, verbose))

				gitlabProject.MustCreateEmptyFile("a.txt", branchName, verbose)
				gitlabProject.MustCreateEmptyFile("b.txt", branchName, verbose)
				gitlabProject.MustCreateEmptyFile("aa/a.txt", branchName, verbose)
				gitlabProject.MustCreateEmptyFile("aa/b.txt", branchName, verbose)
				gitlabProject.MustCreateEmptyFile("evenMore/aa/b.txt", branchName, verbose)
				gitlabProject.MustCreateEmptyFile("evenMore/aa/a.txt", branchName, verbose)

				fileList := gitlabProject.MustGetFilesNames(branchName, verbose)
				expectedFileList := []string{
					"a.txt",
					"aa/a.txt",
					"aa/b.txt",
					"b.txt",
					"evenMore/aa/a.txt",
					"evenMore/aa/b.txt",
				}
				assert.EqualValues(expectedFileList, fileList)

				directoryList := gitlabProject.MustGetDirectoryNames(branchName, verbose)
				exepctedDirectoryList := []string{
					"aa",
					"evenMore",
					"evenMore/aa",
				}
				assert.EqualValues(exepctedDirectoryList, directoryList)
			},
		)
	}
}

func TestGitlabProjectsGetFileList_pagination(t *testing.T) {
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

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				assert.True(gitlabProject.MustExists(verbose))

				branchName := gitlabProject.MustGetDefaultBranchName()

				gitlabProject.MustDeleteAllRepositoryFiles(branchName, verbose)
				assert.True(gitlabProject.MustHasNoRepositoryFiles(branchName, verbose))

				expectedFileList := []string{}
				for i := 0; i < 21; i++ {
					expectedFileList = append(
						expectedFileList,
						fmt.Sprintf("asd%03d.txt", i),
					)
				}

				for _, toCreate := range expectedFileList {
					gitlabProject.MustCreateEmptyFile(toCreate, branchName, verbose)
				}

				fileList := gitlabProject.MustGetFilesNames(branchName, verbose)
				assert.EqualValues(expectedFileList, fileList)
			},
		)
	}
}
