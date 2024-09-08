package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
			MustFormatAsTestname(tt),
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
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Not available in gitlab CI")
		return
	}

	tests := []struct {
		groupPath string
	}{
		{"test_group_id/hello"},
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

				testGroup := gitlab.MustGetGroupByPath(tt.groupPath, verbose)

				partenGroup := testGroup.MustGetParentGroup(verbose)

				partenGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))
				assert.False(partenGroup.MustExists(verbose))

				testGroup.MustCreate(
					&GitlabCreateGroupOptions{
						Verbose: verbose,
					},
				)

				assert.True(testGroup.MustExists(verbose))
				assert.True(partenGroup.MustExists(verbose))

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
					partenGroup.MustGetGroupPath(),
					testGroupById.MustGetParentGroupPath(verbose),
				)

				partenGroup.MustDelete(verbose)
				assert.False(testGroup.MustExists(verbose))
				assert.False(testGroupById.MustExists(verbose))
				assert.False(partenGroup.MustExists(verbose))

			},
		)
	}
}
