package ansibleplaybook_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleplaybook"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_WriteMininalPlaybookExecutingRoles(t *testing.T) {
	t.Run("example-host without remote user", func(t *testing.T) {
		ctx := getCtx()
		playbookPath, err := ansibleplaybook.WriteTemporaryMinimalPlaybookExecutingRoles(
			ctx,
			&ansibleplaybook.MinimalPlaybookOptions{
				Hostname: "example-host",
				Roles:    []string{"role1", "role2"},
			},
		)
		require.NoError(t, err)
		require.NotEmpty(t, playbookPath)

		isYaml, err := yamlutils.IsYamlFile(ctx, playbookPath, &yamlutils.ValidateOptions{
			RefuesePureJson: true,
		})
		require.NoError(t, err)
		require.True(t, isYaml)

		playbook, err := ansibleplaybook.ReadPlaybook(ctx, playbookPath)
		require.NoError(t, err)
		require.Len(t, playbook.Plays, 1)

		play := playbook.Plays[0]
		require.EqualValues(t, "Minimal playbook executing roles", play.Name)
		require.EqualValues(t, []string{"example-host"}, play.Hosts)
		require.EqualValues(t, []string{"role1", "role2"}, play.Roles)
		require.EqualValues(t, "", play.RemoteUser) // Remote user not set if not explicitly defined.
	})

	t.Run("example-host with remote user", func(t *testing.T) {
		ctx := getCtx()
		playbookPath, err := ansibleplaybook.WriteTemporaryMinimalPlaybookExecutingRoles(
			ctx,
			&ansibleplaybook.MinimalPlaybookOptions{
				Hostname: "example-host",
				Roles:    []string{"role1", "role2"},
				RemoteUser: "root",
			},
		)
		require.NoError(t, err)
		require.NotEmpty(t, playbookPath)

		isYaml, err := yamlutils.IsYamlFile(ctx, playbookPath, &yamlutils.ValidateOptions{
			RefuesePureJson: true,
		})
		require.NoError(t, err)
		require.True(t, isYaml)

		playbook, err := ansibleplaybook.ReadPlaybook(ctx, playbookPath)
		require.NoError(t, err)
		require.Len(t, playbook.Plays, 1)

		play := playbook.Plays[0]
		require.EqualValues(t, "Minimal playbook executing roles", play.Name)
		require.EqualValues(t, []string{"example-host"}, play.Hosts)
		require.EqualValues(t, []string{"role1", "role2"}, play.Roles)
		require.EqualValues(t, "root", play.RemoteUser)
	})
}
