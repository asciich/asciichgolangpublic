package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabGroupsGroupByGroupPathExists(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogWarn("Unavailable in continuous integration pipeline")
		return
	}

	tests := []struct {
		groupPath      string
		expectedExists bool
	}{
		{"this/group/does_not_exist", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
						Verbose:                verbose,
					},
				)

				exists := gitlab.MustGroupByGroupPathExists(tt.groupPath, verbose)

				require.EqualValues(tt.expectedExists, exists)
			},
		)
	}
}

func TestGitlabGroupsCreateAndDeleteGroup(t *testing.T) {

	if continuousintegration.IsRunningInGithub() {
		logging.LogWarn("Unavailable in continuous integration pipeline")
		return
	}

	tests := []struct {
		groupName string
	}{
		{"test_group_create_and_delete"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
						Verbose:                verbose,
					},
				)

				groupUnderTest := gitlab.MustGetGroupByPath(tt.groupName, verbose)

				for i := 0; i < 2; i++ {
					gitlab.MustDeleteGroupByPath(
						tt.groupName,
						verbose,
					)
					require.False(groupUnderTest.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					createdGroup := gitlab.MustCreateGroupByPath(
						tt.groupName,
						&GitlabCreateGroupOptions{
							Verbose: verbose,
						},
					)
					require.True(createdGroup.MustExists(verbose))
					require.True(groupUnderTest.MustExists(verbose))
				}
				require.Greater(groupUnderTest.MustGetId(verbose), 0)

				for i := 0; i < 2; i++ {
					gitlab.MustDeleteGroupByPath(
						tt.groupName,
						verbose,
					)
					require.False(groupUnderTest.MustExists(verbose))
				}
			},
		)
	}
}
