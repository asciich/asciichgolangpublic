package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				privateProject, err := gitlab.GetPersonalProjectByName(ctx, tt.projectName)
				require.NoError(t, err)

				isPersonalProject, err := privateProject.IsPersonalProject(ctx)
				require.NoError(t, err)
				require.True(t, isPersonalProject)

				for i := 0; i < 2; i++ {
					err = privateProject.Delete(ctx)
					require.NoError(t, err)

					exists, err := privateProject.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)

					isPersonalProject, err := privateProject.IsPersonalProject(ctx)
					require.NoError(t, err)
					require.True(t, isPersonalProject)
				}

				for i := 0; i < 2; i++ {
					err := privateProject.Create(ctx)
					require.NoError(t, err)

					exists, err := privateProject.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					isPersonalProject, err := privateProject.IsPersonalProject(ctx)
					require.NoError(t, err)
					require.True(t, isPersonalProject)
				}

				for i := 0; i < 2; i++ {
					err := privateProject.Delete(ctx)
					require.NoError(t, err)

					exists, err := privateProject.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)

					isPersonalProject, err := privateProject.IsPersonalProject(ctx)
					require.NoError(t, err)
					require.True(t, isPersonalProject)
				}
			},
		)
	}
}
