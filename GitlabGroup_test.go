package asciichgolangpublic

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				assert.EqualValues(
					tt.expectedGroupName,
					testGroup.MustGetGroupPath(),
				)
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
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				assert.EqualValues(
					tt.expectedGroupName,
					testGroup.MustGetGroupName(),
				)
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
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				assert.EqualValues(
					tt.expectedIsSubgroup,
					testGroup.MustIsSubgroup(),
				)
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
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				assert.EqualValues(
					tt.expectedParentGroupPath,
					testGroup.MustGetParentGroupPath(verbose),
				)
			},
		)
	}
}

// Validate if getting the gitlab group by path and by id works.
func TestGitlabGroupByPathAndId(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		groupPath string
	}{
		{"test_group_id/hello"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
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

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				parentGroup := testGroup.MustGetParentGroup(verbose)

				parentGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))
				assert.False(parentGroup.MustExists(verbose))

				testGroup.MustCreate(
					&GitlabCreateGroupOptions{
						Verbose: verbose,
					},
				)

				assert.True(testGroup.MustExists(verbose))
				assert.True(parentGroup.MustExists(verbose))

				testGroupId := testGroup.MustGetId(verbose)
				testGroupById := gitlab.MustGetGroupById(testGroupId, verbose)

				assert.True(testGroupById.MustExists(verbose))
				assert.EqualValues(
					tt.groupPath,
					testGroupById.MustGetGroupPath(),
				)
				assert.EqualValues(
					testGroup.MustGetGroupName(),
					testGroupById.MustGetGroupName(),
				)
				assert.EqualValues(
					parentGroup.MustGetGroupPath(),
					testGroupById.MustGetParentGroupPath(verbose),
				)

				parentGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))
				assert.False(testGroupById.MustExists(verbose))
				assert.False(parentGroup.MustExists(verbose))

			},
		)
	}
}

func TestGitlabGroupListProjects(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
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

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
						Verbose:                verbose,
					},
				)

				const testGroupName string = "test_projects_in_group"
				testGroup := gitlab.MustGetGroupByPath(testGroupName, verbose)

				testGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))

				nProjects := 25
				projectPaths := []string{}
				for i := 0; i < nProjects; i++ {
					projectPath := fmt.Sprintf("%s/project_%d", testGroupName, i)
					projectPaths = append(projectPaths, projectPath)

					gitlab.MustCreateProject(
						&GitlabCreateProjectOptions{
							ProjectPath: projectPath,
							Verbose:     verbose,
						},
					)
				}

				time.Sleep(3 * time.Second)
				listedProjectPaths := testGroup.MustListProjectPaths(
					&GitlabListProjectsOptions{
						Verbose: verbose,
					},
				)

				assert.Len(listedProjectPaths, nProjects)

				for _, toCheck := range projectPaths {
					assert.True(
						slices.Contains(
							listedProjectPaths,
							toCheck,
						),
					)
				}

				testGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))
			},
		)
	}
}
