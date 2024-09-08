package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabGroupsGroupByGroupPathExists(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in continuous integration pipeline")
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
						Verbose:                verbose,
					},
				)

				exists := gitlab.MustGroupByGroupPathExists(tt.groupPath, verbose)

				assert.EqualValues(tt.expectedExists, exists)
			},
		)
	}
}

func TestGitlabGroupsCreateAndDeleteGroup(t *testing.T) {

	if ContinuousIntegration().IsRunningInGithub() {
		LogWarn("Unavailable in continuous integration pipeline")
		return
	}

	tests := []struct {
		groupName string
	}{
		{"test_group_create_and_delete"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

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
					assert.False(groupUnderTest.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					createdGroup := gitlab.MustCreateGroupByPath(
						tt.groupName,
						&GitlabCreateGroupOptions{
							Verbose: verbose,
						},
					)
					assert.True(createdGroup.MustExists(verbose))
					assert.True(groupUnderTest.MustExists(verbose))
				}
				assert.Greater(groupUnderTest.MustGetId(verbose), 0)

				for i := 0; i < 2; i++ {
					gitlab.MustDeleteGroupByPath(
						tt.groupName,
						verbose,
					)
					assert.False(groupUnderTest.MustExists(verbose))
				}
			},
		)
	}
}
