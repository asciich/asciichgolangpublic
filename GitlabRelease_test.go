package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"

	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabReleaseCreateAndDelete(t *testing.T) {
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

					require.False(release.MustExists(verbose))
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

					require.True(release.MustExists(verbose))
					require.True(tag.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					release.MustDelete(
						&GitlabDeleteReleaseOptions{
							Verbose:                verbose,
							DeleteCorrespondingTag: true,
						},
					)

					require.False(release.MustExists(verbose))
					require.False(tag.MustExists(verbose))
				}
			},
		)
	}
}

func TestGitlabRelease_ReleaseLinks(t *testing.T) {
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
				require.True(release.MustExists(verbose))

				const releaseLink = "https://asciich.ch/release/link/1"

				require.False(release.MustHasReleaseLinks(verbose))

				for i := 0; i < 2; i++ {
					release.MustCreateReleaseLink(
						&GitlabCreateReleaseLinkOptions{
							Url:     releaseLink,
							Name:    "testReleaseLink",
							Verbose: verbose,
						},
					)
					require.True(release.MustHasReleaseLinks(verbose))
				}

				require.EqualValues(
					[]string{releaseLink},
					release.MustListReleaseLinkUrls(verbose),
				)

				release.MustDelete(
					&GitlabDeleteReleaseOptions{
						Verbose:                verbose,
						DeleteCorrespondingTag: true,
					},
				)

				require.False(release.MustExists(verbose))
			},
		)
	}
}

func TestGitlabRelease_CreateNewPatchRelease(t *testing.T) {
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
				require.True(release.MustExists(verbose))

				project.MustWriteFileContentInDefaultBranch(
					&GitlabWriteFileOptions{
						Path:          "random.txt",
						Content:       []byte(mustutils.Must(randomgenerator.GetRandomString(50))),
						CommitMessage: "Dummy change to test release.",
						Verbose:       verbose,
					},
				)

				nextPatchRelease := project.MustCreateNextPatchReleaseFromLatestCommitInDefaultBranch("next patch release", verbose)

				require.EqualValues(
					"v1.2.4",
					nextPatchRelease.MustGetName(),
				)

				nextMinorRelease := project.MustCreateNextMinorReleaseFromLatestCommitInDefaultBranch("next minor release", verbose)

				require.EqualValues(
					"v1.3.0",
					nextMinorRelease.MustGetName(),
				)

				nextMajorRelease := project.MustCreateNextMajorReleaseFromLatestCommitInDefaultBranch("next minor release", verbose)

				require.EqualValues(
					"v2.0.0",
					nextMajorRelease.MustGetName(),
				)
			},
		)
	}
}
