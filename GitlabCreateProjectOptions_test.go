package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabCreateProjectOptionsGetProjectName(t *testing.T) {
	tests := []struct {
		createOptions       *GitlabCreateProjectOptions
		expectedProjectName string
	}{
		{&GitlabCreateProjectOptions{ProjectPath: "group/hello"}, "hello"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, "helloWorld"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, "hello"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, "helloWorld"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, "hello"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, "helloWorld"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, "hello"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				projectName, err := tt.createOptions.GetProjectName()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedProjectName, projectName)
			},
		)
	}
}

func TestGitlabCreateProjectOptionsGetGroupNames(t *testing.T) {
	tests := []struct {
		createOptions      *GitlabCreateProjectOptions
		expectedGroupNames []string
	}{
		{&GitlabCreateProjectOptions{ProjectPath: "group/hello"}, []string{"group"}},
		{&GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, []string{"group"}},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, []string{"group", "subgroup"}},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, []string{"group", "subgroup"}},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, []string{"group"}},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, []string{"group"}},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, []string{"group", "subgroup"}},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, []string{"group", "subgroup"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				groupNames, err := tt.createOptions.GetGroupNames(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedGroupNames, groupNames)
			},
		)
	}
}

func TestGitlabCreateProjectOptionsGetGroupPath(t *testing.T) {
	tests := []struct {
		createOptions     *GitlabCreateProjectOptions
		expectedGroupPath string
	}{
		{&GitlabCreateProjectOptions{ProjectPath: "group/hello"}, "group"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, "group"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, "group/subgroup"},
		{&GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, "group/subgroup"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, "group"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, "group"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, "group/subgroup"},
		{&GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, "group/subgroup"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				groupPath, err := tt.createOptions.GetGroupPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedGroupPath, groupPath)
			},
		)
	}
}
