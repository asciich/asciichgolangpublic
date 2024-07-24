package asciichgolangpublic

import (
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestGitlabCreateGroupOptionsGetGroupName(t *testing.T) {
	tests := []struct {
		createOptions       *GitlabCreateGroupOptions
		expectedProjectName string
	}{
		{&GitlabCreateGroupOptions{GroupPath: "group/hello"}, "hello"},
		{&GitlabCreateGroupOptions{GroupPath: "group/helloWorld"}, "helloWorld"},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/hello"}, "hello"},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/helloWorld"}, "helloWorld"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/hello"}, "hello"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/helloWorld"}, "helloWorld"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/hello"}, "hello"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/helloWorld"}, "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				projectName := tt.createOptions.MustGetGroupName()
				assert.EqualValues(tt.expectedProjectName, projectName)
			},
		)
	}
}

func TestGitlabCreateGroupOptionsIsSubgroup(t *testing.T) {
	tests := []struct {
		createOptions      *GitlabCreateGroupOptions
		expectedIsSubgroup bool
	}{
		{&GitlabCreateGroupOptions{GroupPath: "group/hello"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "group/helloWorld"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/hello"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/helloWorld"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "/group/hello"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "/group/helloWorld"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/hello"}, true},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/helloWorld"}, true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				projectName := tt.createOptions.MustIsSubgroup()
				assert.EqualValues(tt.expectedIsSubgroup, projectName)
			},
		)
	}
}

func TestGitlabCreateGroupOptionsGetParentGroupPath(t *testing.T) {
	tests := []struct {
		createOptions           *GitlabCreateGroupOptions
		expectedParentGroupPath string
	}{
		{&GitlabCreateGroupOptions{GroupPath: "group/hello"}, "group"},
		{&GitlabCreateGroupOptions{GroupPath: "group/helloWorld"}, "group"},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/hello"}, "group/subgroup"},
		{&GitlabCreateGroupOptions{GroupPath: "group/subgroup/helloWorld"}, "group/subgroup"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/hello"}, "group"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/helloWorld"}, "group"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/hello"}, "group/subgroup"},
		{&GitlabCreateGroupOptions{GroupPath: "/group/subgroup/helloWorld"}, "group/subgroup"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				parentGroupPath := tt.createOptions.MustGetParentGroupPath()
				assert.EqualValues(tt.expectedParentGroupPath, parentGroupPath)
			},
		)
	}
}
