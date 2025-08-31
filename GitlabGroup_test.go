package asciichgolangpublic

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabGroupGetGroupPath(t *testing.T) {
	tests := []struct {
		groupPath         string
		expectedGroupName string
	}{
		{"group/hello", "group/hello"},
		{"group/hello/", "group/hello"},
		{"group/helloWorld", "group/helloWorld"},
		{"group/subgroup/hello", "group/subgroup/hello"},
		{"group/subgroup/helloWorld", "group/subgroup/helloWorld"},
		{"/group/hello", "group/hello"},
		{"/group/hello/", "group/hello"},
		{"/group/helloWorld", "group/helloWorld"},
		{"/group/subgroup/hello", "group/subgroup/hello"},
		{"/group/subgroup/helloWorld", "group/subgroup/helloWorld"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				testGroup, err := gitlab.GetGroupByPath(ctx, tt.groupPath)
				require.NoError(t, err)

				groupPath, err := testGroup.GetGroupPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedGroupName, groupPath)
			},
		)
	}
}

func TestGitlabGroupGetGroupName(t *testing.T) {
	tests := []struct {
		groupPath         string
		expectedGroupName string
	}{
		{"group/hello", "hello"},
		{"group/hello/", "hello"},
		{"group/helloWorld", "helloWorld"},
		{"group/subgroup/hello", "hello"},
		{"group/subgroup/helloWorld", "helloWorld"},
		{"/group/hello", "hello"},
		{"/group/hello/", "hello"},
		{"/group/helloWorld", "helloWorld"},
		{"/group/subgroup/hello", "hello"},
		{"/group/subgroup/helloWorld", "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				testGroup, err := gitlab.GetGroupByPath(ctx, tt.groupPath)
				require.NoError(t, err)

				groupName, err := testGroup.GetGroupName(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedGroupName, groupName)
			},
		)
	}
}

func TestGitlabGroupIsSubgroup(t *testing.T) {
	tests := []struct {
		groupPath          string
		expectedIsSubgroup bool
	}{
		{"hello", false},
		{"/hello", false},
		{"hello/", false},
		{"/hello/", false},
		{"group/hello", true},
		{"group/helloWorld", true},
		{"group/subgroup/hello", true},
		{"group/subgroup/helloWorld", true},
		{"/group/hello", true},
		{"/group/helloWorld", true},
		{"/group/subgroup/hello", true},
		{"/group/subgroup/helloWorld", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				testGroup, err := gitlab.GetGroupByPath(ctx, tt.groupPath)
				require.NoError(t, err)

				isSubGroup, err := testGroup.IsSubgroup(ctx)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedIsSubgroup, isSubGroup)
			},
		)
	}
}

func TestGitlabGroupGetParentGroupPath(t *testing.T) {
	tests := []struct {
		groupPath               string
		expectedParentGroupPath string
	}{
		{"group/hello", "group"},
		{"group/helloWorld", "group"},
		{"group/subgroup/hello", "group/subgroup"},
		{"group/subgroup/helloWorld", "group/subgroup"},
		{"/group/hello", "group"},
		{"/group/helloWorld", "group"},
		{"/group/subgroup/hello", "group/subgroup"},
		{"/group/subgroup/helloWorld", "group/subgroup"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				testGroup, err := gitlab.GetGroupByPath(ctx, tt.groupPath)
				require.NoError(t, err)

				parentGroupPath, err := testGroup.GetParentGroupPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedParentGroupPath, parentGroupPath)
			},
		)
	}
}

// Validate if getting the gitlab group by path and by id works.
func TestGitlabGroupByPathAndId(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		groupPath string
	}{
		{"test_group_id/hello"},
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

				testGroup, err := gitlab.GetGroupByPath(ctx, tt.groupPath)
				require.NoError(t, err)

				parentGroup, err := testGroup.GetParentGroup(ctx)
				require.NoError(t, err)

				err = parentGroup.Delete(ctx)
				require.NoError(t, err)

				exists, err := testGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				exists, err = parentGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				err = testGroup.Create(ctx)
				require.NoError(t, err)

				exists, err = testGroup.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
				exists, err = parentGroup.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				testGroupId, err := testGroup.GetId(ctx)
				require.NoError(t, err)

				testGroupById, err := gitlab.GetGroupById(ctx, testGroupId)
				require.NoError(t, err)

				exists, err = testGroupById.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				groupPath, err := testGroupById.GetGroupPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.groupPath, groupPath)

				testGroupName, err := testGroup.GetGroupName(ctx)
				require.NoError(t, err)

				testGroupByIdName, err := testGroupById.GetGroupName(ctx)
				require.NoError(t, err)

				require.EqualValues(t, testGroupName, testGroupByIdName)

				parentGroupPath, err := parentGroup.GetGroupPath(ctx)
				require.NoError(t, err)
				testGroupByIdPath, err := testGroupById.GetGroupPath(ctx)
				require.NoError(t, err)

				require.EqualValues(t, parentGroupPath, testGroupByIdPath)

				err = parentGroup.Delete(ctx)
				require.NoError(t, err)

				exists, err = testGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
				exists, err = testGroupById.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
				exists, err = parentGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

			},
		)
	}
}

func TestGitlabGroupListProjects(t *testing.T) {
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

				err = gitlab.Authenticate(
					ctx,
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
					},
				)
				require.NoError(t, err)

				const testGroupName string = "test_projects_in_group"
				testGroup, err := gitlab.GetGroupByPath(ctx, testGroupName)
				require.NoError(t, err)

				err = testGroup.Delete(ctx)
				require.NoError(t, err)

				exists, err := testGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				nProjects := 25
				projectPaths := []string{}
				for i := 0; i < nProjects; i++ {
					projectPath := fmt.Sprintf("%s/project_%d", testGroupName, i)
					projectPaths = append(projectPaths, projectPath)

					_, err = gitlab.CreateProject(
						ctx,
						&GitlabCreateProjectOptions{
							ProjectPath: projectPath,
						},
					)
					require.NoError(t, err)
				}

				time.Sleep(3 * time.Second)
				listedProjectPaths, err := testGroup.ListProjectPaths(ctx, &GitlabListProjectsOptions{})
				require.NoError(t, err)

				require.Len(t, listedProjectPaths, nProjects)

				for _, toCheck := range projectPaths {
					require.True(t,
						slices.Contains(
							listedProjectPaths,
							toCheck,
						),
					)
				}

				err = testGroup.Delete(ctx)
				require.NoError(t, err)

				exists, err = testGroup.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}
