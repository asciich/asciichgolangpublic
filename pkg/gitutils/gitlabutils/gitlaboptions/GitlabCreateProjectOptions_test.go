package gitlaboptions_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlaboptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGitlabCreateProjectOptionsGetProjectName(t *testing.T) {
	tests := []struct {
		createOptions       *gitlaboptions.GitlabCreateProjectOptions
		expectedProjectName string
	}{
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/hello"}, "hello"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, "helloWorld"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, "hello"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, "helloWorld"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, "hello"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, "helloWorld"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, "hello"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, "helloWorld"},
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
		createOptions      *gitlaboptions.GitlabCreateProjectOptions
		expectedGroupNames []string
	}{
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/hello"}, []string{"group"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, []string{"group"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, []string{"group", "subgroup"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, []string{"group", "subgroup"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, []string{"group"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, []string{"group"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, []string{"group", "subgroup"}},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, []string{"group", "subgroup"}},
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
		createOptions     *gitlaboptions.GitlabCreateProjectOptions
		expectedGroupPath string
	}{
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/hello"}, "group"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/helloWorld"}, "group"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/hello"}, "group/subgroup"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "group/subgroup/helloWorld"}, "group/subgroup"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/hello"}, "group"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/helloWorld"}, "group"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/hello"}, "group/subgroup"},
		{&gitlaboptions.GitlabCreateProjectOptions{ProjectPath: "/group/subgroup/helloWorld"}, "group/subgroup"},
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
