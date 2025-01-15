package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabPersonalProjectsCreateAndDelete(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not implemented to run on github")
		return
	}

	tests := []struct {
		projectName string
	}{
		{"testproject1"},
		{"testproject2"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				privateProject := gitlab.MustGetPersonalProjectByName(tt.projectName, verbose)
				assert.True(privateProject.MustIsPersonalProject())

				for i := 0; i < 2; i++ {
					privateProject.MustDelete(verbose)
					assert.False(privateProject.MustExists(verbose))

					assert.True(privateProject.MustIsPersonalProject())
				}

				for i := 0; i < 2; i++ {
					privateProject.MustCreate(verbose)
					assert.True(privateProject.MustExists(verbose))

					assert.True(privateProject.MustIsPersonalProject())
				}

				for i := 0; i < 2; i++ {
					privateProject.MustDelete(verbose)
					assert.False(privateProject.MustExists(verbose))

					assert.True(privateProject.MustIsPersonalProject())
				}
			},
		)
	}
}
