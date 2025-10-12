package ansibleplaybook_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleplaybook"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func TestReadPlaybook(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		ctx := getCtx()
		playbook, err := ansibleplaybook.ReadPlaybook(ctx, "")
		require.Error(t, err)
		require.Nil(t, playbook)
	})

	t.Run("playbook with hosts as list", func(t *testing.T) {
		ctx := getCtx()

		content := "---\n"
		content += "- name: Update web servers\n"
		content += "  hosts:\n"
		content += "    - webservers\n"

		playbookPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)

		isYaml, err := yamlutils.IsYamlFile(ctx, playbookPath, &yamlutils.ValidateOptions{RefuesePureJson: true})
		require.NoError(t, err)
		require.True(t, isYaml)

		playbook, err := ansibleplaybook.ReadPlaybook(ctx, playbookPath)
		require.NoError(t, err)
		require.NotNil(t, playbook)

		require.Len(t, playbook.Plays, 1)
	})

	t.Run("playbook with hosts as string", func(t *testing.T) {
		ctx := getCtx()

		content := "---\n"
		content += "- name: Update web servers\n"
		content += "  hosts: webservers\n"

		playbookPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)

		isYaml, err := yamlutils.IsYamlFile(ctx, playbookPath, &yamlutils.ValidateOptions{RefuesePureJson: true})
		require.NoError(t, err)
		require.True(t, isYaml)

		playbook, err := ansibleplaybook.ReadPlaybook(ctx, playbookPath)
		require.NoError(t, err)
		require.NotNil(t, playbook)

		require.Len(t, playbook.Plays, 1)
	})

	t.Run("example_playbook.yaml", func(t *testing.T) {
		ctx := getCtx()
		playbook, err := ansibleplaybook.ReadPlaybook(ctx, "./testdata/example_playbook.yaml")
		require.NoError(t, err)
		require.NotNil(t, playbook)

		require.Len(t, playbook.Plays, 2)
	})
}
