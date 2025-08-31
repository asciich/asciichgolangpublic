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
				ctx := getCtx()
				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				const releaseName string = "test_release"
				const releaseDescription string = "Release description."

				project, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)
				release, err := project.GetReleaseByName(releaseName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = release.Delete(
						ctx,
						&GitlabDeleteReleaseOptions{
							DeleteCorrespondingTag: true,
						},
					)
					require.NoError(t, err)

					exists, err := release.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}

				var tag *GitlabTag

				for i := 0; i < 2; i++ {
					release, err = project.CreateReleaseFromLatestCommitInDefaultBranch(
						ctx,
						&GitlabCreateReleaseOptions{
							Name:        releaseName,
							Description: releaseDescription,
						},
					)
					require.NoError(t, err)

					tag, err = release.GetTag()
					require.NoError(t, err)

					exists, err := release.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					exists, err = tag.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)
				}

				for i := 0; i < 2; i++ {
					err = release.Delete(
						ctx,
						&GitlabDeleteReleaseOptions{
							DeleteCorrespondingTag: true,
						},
					)

					exists, err := release.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)

					exists, err = tag.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
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
				ctx := getCtx()
				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})

				const projectPath string = "test_group/testproject"

				const releaseName string = "test_release"
				const releaseDescription string = "Release description."

				project, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)
				require.NoError(t, err)

				release, err := project.GetReleaseByName(releaseName)
				require.NoError(t, err)

				err = release.Delete(ctx, &GitlabDeleteReleaseOptions{DeleteCorrespondingTag: true})

				release, err = project.CreateReleaseFromLatestCommitInDefaultBranch(
					ctx,
					&GitlabCreateReleaseOptions{
						Name:        releaseName,
						Description: releaseDescription,
					},
				)

				exists, err := release.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				const releaseLink = "https://asciich.ch/release/link/1"

				hasReleaseLinks, err := release.HasReleaseLinks(ctx)
				require.NoError(t, err)
				require.False(t, hasReleaseLinks)

				for i := 0; i < 2; i++ {
					_, err = release.CreateReleaseLink(
						ctx,
						&GitlabCreateReleaseLinkOptions{
							Url:  releaseLink,
							Name: "testReleaseLink",
						},
					)
					hasLinks, err := release.HasReleaseLinks(ctx)
					require.NoError(t, err)
					require.True(t, hasLinks)
				}

				linkUrls, err := release.ListReleaseLinkUrls(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []string{releaseLink}, linkUrls)

				err = release.Delete(
					ctx,
					&GitlabDeleteReleaseOptions{
						DeleteCorrespondingTag: true,
					},
				)
				require.NoError(t, err)

				exists, err = release.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
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
				ctx := getCtx()
				gitlabFQDN := "gitlab.asciich.ch"

				gitlab, err := GetGitlabByFQDN(gitlabFQDN)
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				const projectPath string = "test_group/testproject"

				const releaseDescription string = "Release description."

				project, err := gitlab.GetGitlabProjectByPath(ctx, projectPath)

				err = project.DeleteAllReleases(
					ctx,
					&GitlabDeleteReleaseOptions{
						DeleteCorrespondingTag: true,
					},
				)

				release, err := project.CreateReleaseFromLatestCommitInDefaultBranch(
					ctx,
					&GitlabCreateReleaseOptions{
						Name:        "v1.2.3",
						Description: releaseDescription,
					},
				)
				exists, err := release.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				_, err = project.WriteFileContentInDefaultBranch(
					ctx,
					&GitlabWriteFileOptions{
						Path:          "random.txt",
						Content:       []byte(mustutils.Must(randomgenerator.GetRandomString(50))),
						CommitMessage: "Dummy change to test release.",
					},
				)
				require.NoError(t, err)

				nextPatchRelease, err := project.CreateNextPatchReleaseFromLatestCommitInDefaultBranch(ctx, "next patch release")
				require.NoError(t, err)

				name, err := nextPatchRelease.GetName()
				require.NoError(t, err)
				require.EqualValues(t, "v1.2.4", name)

				nextMinorRelease, err := project.CreateNextMinorReleaseFromLatestCommitInDefaultBranch(ctx, "next minor release")
				require.NoError(t, err)

				name, err = nextMinorRelease.GetName()
				require.NoError(t, err)
				require.EqualValues(t, "v1.3.0", name)

				nextMajorRelease, err := project.CreateNextMajorReleaseFromLatestCommitInDefaultBranch(ctx, "next minor release")
				require.NoError(t, err)

				name, err = nextMajorRelease.GetName()
				require.NoError(t, err)
				require.EqualValues(t, "v2.0.0", name)
			},
		)
	}
}
