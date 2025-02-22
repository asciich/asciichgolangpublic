package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabPersonalProjectsCreateAndDelete(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				privateProject := gitlab.MustGetPersonalProjectByName(tt.projectName, verbose)
				require.True(privateProject.MustIsPersonalProject())

				for i := 0; i < 2; i++ {
					privateProject.MustDelete(verbose)
					require.False(privateProject.MustExists(verbose))

					require.True(privateProject.MustIsPersonalProject())
				}

				for i := 0; i < 2; i++ {
					privateProject.MustCreate(verbose)
					require.True(privateProject.MustExists(verbose))

					require.True(privateProject.MustIsPersonalProject())
				}

				for i := 0; i < 2; i++ {
					privateProject.MustDelete(verbose)
					require.False(privateProject.MustExists(verbose))

					require.True(privateProject.MustIsPersonalProject())
				}
			},
		)
	}
}
