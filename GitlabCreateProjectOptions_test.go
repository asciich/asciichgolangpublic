package asciichgolangpublic

import (
	"testing"


	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				projectName := tt.createOptions.MustGetProjectName()
				assert.EqualValues(tt.expectedProjectName, projectName)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				groupNames := tt.createOptions.MustGetGroupNames(verbose)
				assert.EqualValues(tt.expectedGroupNames, groupNames)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				groupPath := tt.createOptions.MustGetGroupPath(verbose)
				assert.EqualValues(tt.expectedGroupPath, groupPath)
			},
		)
	}
}
