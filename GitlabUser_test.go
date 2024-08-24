package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabUserGetCurrentUserName(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Not implemented to run on github")
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
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				userName := gitlab.MustGetCurrentUserName(verbose)
				assert.EqualValues("reto", userName)
			},
		)
	}
}
