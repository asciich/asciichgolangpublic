package asciichgolangpublic

/* TODO enable again
import (
	"testing"


	"github.com/stretchr/testify/assert"
)


func TestGitlabGroupsGroupByGroupPathExists(t *testing.T) {
	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogWarn("Unavailable in continuous integration pipeline")
		return
	}

	tests := []struct {
		groupPath      string
		expectedExists bool
	}{
		{"this/group/does_not_exist", false},
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

				exists := gitlab.MustGroupByGroupPathExists(tt.groupPath)

				assert.EqualValues(tt.expectedExists, exists)
			},
		)
	}
}

func TestGitlabGroupsCreateGroup(t *testing.T) {

	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogWarn("Unavailable in continuous integration pipeline")
		return
	}

	tests := []struct {
		testcase string
	}{
		{"testcase"},
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

				createdGroup := gitlab.MustCreateGroup(
					&GitlabCreateGroupOptions{
						GroupPath: "test_group",
						Verbose:   verbose,
					},
				)

				assert.Greater(createdGroup.MustGetId(), 0)
			},
		)
	}
}
*/
