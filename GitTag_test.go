package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getGitTagToTest(implementationName string) (gitTag gitinterfaces.GitTag) {
	if implementationName == "gitRepositoryTag" {
		return NewGitRepositoryTag()
	}

	if implementationName == "gitlabTag" {
		return NewGitlabTag()
	}

	logging.LogFatalWithTracef(
		"Unknown implementation name: '%s'",
		implementationName,
	)

	return nil
}

func TestGitlabTag_IsVersionTag(t *testing.T) {
	tests := []struct {
		implementationName   string
		tagName              string
		expectedIsVersionTag bool
	}{
		{"gitRepositoryTag", "v0.1.2", true},
		{"gitRepositoryTag", "abc", false},
		{"gitRepositoryTag", "v20241229_140707", true},
		{"gitlabTag", "v0.1.2", true},
		{"gitlabTag", "abc", false},
		{"gitlabTag", "v20241229_140707", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				toTest := getGitTagToTest(tt.implementationName)

				err := toTest.SetName(tt.tagName)
				require.NoError(t, err)

				isVersionTag, err := toTest.IsVersionTag()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedIsVersionTag, isVersionTag)
			},
		)
	}
}
