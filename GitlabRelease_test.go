package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/randomgenerator"

	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabReleaseCreateAndDelete(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				const releaseName string = "test_release"
				const releaseDescription string = "Release description."

				project := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				release := project.MustGetReleaseByName(releaseName)

				for i := 0; i < 2; i++ {
					release.MustDelete(
						&GitlabDeleteReleaseOptions{
							Verbose:                verbose,
							DeleteCorrespondingTag: true,
						},
					)

					assert.False(release.MustExists(verbose))
				}

				var tag *GitlabTag

				for i := 0; i < 2; i++ {
					release = project.MustCreateReleaseFromLatestCommitInDefaultBranch(
						&GitlabCreateReleaseOptions{
							Name:        releaseName,
							Verbose:     verbose,
							Description: releaseDescription,
						},
					)
					tag = release.MustGetTag()

					assert.True(release.MustExists(verbose))
					assert.True(tag.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					release.MustDelete(
						&GitlabDeleteReleaseOptions{
							Verbose:                verbose,
							DeleteCorrespondingTag: true,
						},
					)

					assert.False(release.MustExists(verbose))
					assert.False(tag.MustExists(verbose))
				}
			},
		)
	}
}

func TestGitlabRelease_ReleaseLinks(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				const releaseName string = "test_release"
				const releaseDescription string = "Release description."

				project := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)
				release := project.MustGetReleaseByName(releaseName)

				release.MustDelete(
					&GitlabDeleteReleaseOptions{
						Verbose:                verbose,
						DeleteCorrespondingTag: true,
					},
				)

				release = project.MustCreateReleaseFromLatestCommitInDefaultBranch(
					&GitlabCreateReleaseOptions{
						Name:        releaseName,
						Verbose:     verbose,
						Description: releaseDescription,
					},
				)
				assert.True(release.MustExists(verbose))

				const releaseLink = "https://asciich.ch/release/link/1"

				assert.False(release.MustHasReleaseLinks(verbose))

				for i := 0; i < 2; i++ {
					release.MustCreateReleaseLink(
						&GitlabCreateReleaseLinkOptions{
							Url:     releaseLink,
							Name:    "testReleaseLink",
							Verbose: verbose,
						},
					)
					assert.True(release.MustHasReleaseLinks(verbose))
				}

				assert.EqualValues(
					[]string{releaseLink},
					release.MustListReleaseLinkUrls(verbose),
				)

				release.MustDelete(
					&GitlabDeleteReleaseOptions{
						Verbose:                verbose,
						DeleteCorrespondingTag: true,
					},
				)

				assert.False(release.MustExists(verbose))
			},
		)
	}
}

func TestGitlabRelease_CreateNewPatchRelease(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlabFQDN := "gitlab.asciich.ch"

				gitlab := MustGetGitlabByFqdn(gitlabFQDN)
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				const projectPath string = "test_group/testproject"

				const releaseDescription string = "Release description."

				project := gitlab.MustGetGitlabProjectByPath(projectPath, verbose)

				project.MustDeleteAllReleases(
					&GitlabDeleteReleaseOptions{
						Verbose:                verbose,
						DeleteCorrespondingTag: true,
					},
				)

				release := project.MustCreateReleaseFromLatestCommitInDefaultBranch(
					&GitlabCreateReleaseOptions{
						Name:        "v1.2.3",
						Description: releaseDescription,
						Verbose:     verbose,
					},
				)
				assert.True(release.MustExists(verbose))

				project.MustWriteFileContentInDefaultBranch(
					&GitlabWriteFileOptions{
						Path:          "random.txt",
						Content:       []byte(randomgenerator.MustGetRandomString(50)),
						CommitMessage: "Dummy change to test release.",
						Verbose:       verbose,
					},
				)

				nextPatchRelease := project.MustCreateNextPatchReleaseFromLatestCommitInDefaultBranch("next patch release", verbose)

				assert.EqualValues(
					"v1.2.4",
					nextPatchRelease.MustGetName(),
				)

				nextMinorRelease := project.MustCreateNextMinorReleaseFromLatestCommitInDefaultBranch("next minor release", verbose)

				assert.EqualValues(
					"v1.3.0",
					nextMinorRelease.MustGetName(),
				)

				nextMajorRelease := project.MustCreateNextMajorReleaseFromLatestCommitInDefaultBranch("next minor release", verbose)

				assert.EqualValues(
					"v2.0.0",
					nextMajorRelease.MustGetName(),
				)
			},
		)
	}
}
