package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getGitTagToTest(implementationName string) (gitTag GitTag) {
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
				require := require.New(t)

				toTest := getGitTagToTest(tt.implementationName)

				toTest.MustSetName(tt.tagName)

				require.EqualValues(
					tt.expectedIsVersionTag,
					toTest.MustIsVersionTag(),
				)
			},
		)
	}
}
