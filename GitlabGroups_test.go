package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabGroupsGroupByGroupPathExists(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(
					ctx,
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
					},
				)
				require.NoError(t, err)

				exists, err := gitlab.GroupByGroupPathExists(ctx, tt.groupPath)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExists, exists)
			},
		)
	}
}

func TestGitlabGroupsCreateAndDeleteGroup(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		groupName string
	}{
		{"test_group_create_and_delete"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(
					ctx,
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
					},
				)
				require.NoError(t, err)

				groupUnderTest, err := gitlab.GetGroupByPath(ctx, tt.groupName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = gitlab.DeleteGroupByPath(ctx, tt.groupName)
					require.NoError(t, err)

					exists, err := groupUnderTest.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}

				for i := 0; i < 2; i++ {
					createdGroup, err := gitlab.CreateGroupByPath(ctx, tt.groupName)
					require.NoError(t, err)

					exists, err := createdGroup.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					exists, err = groupUnderTest.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)
				}
				id, err := groupUnderTest.GetId(ctx)
				require.NoError(t, err)
				require.Greater(t, id, 0)

				for i := 0; i < 2; i++ {
					err := gitlab.DeleteGroupByPath(ctx, tt.groupName)
					require.NoError(t, err)

					exists, err := groupUnderTest.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
