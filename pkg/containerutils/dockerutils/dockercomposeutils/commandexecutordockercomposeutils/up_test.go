package commandexecutordockercomposeutils_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/commandexecutordockercomposeutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockercomposeutils/dockercomposeoptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

type commandExecutorImplementation struct {
	Name     string
	Executor commandexecutorinterfaces.CommandExecutor
}

func getCommandExecutorImplementations() []commandExecutorImplementation {
	return []commandExecutorImplementation{
		{Name: "exec", Executor: commandexecutorexecoo.Exec()},
		{Name: "bash", Executor: commandexecutorbashoo.Bash()},
	}
}

// writeComposeFile writes a compose file with the given project name and one sleeping alpine container
// per service name into 'directory'. Returns the path of the written compose file.
func writeComposeFile(t *testing.T, directory string, projectName string, serviceNames ...string) string {
	t.Helper()

	require.NotEmpty(t, serviceNames)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("name: %s\nservices:\n", projectName))
	for _, serviceName := range serviceNames {
		sb.WriteString(fmt.Sprintf("  %s:\n    image: alpine:3\n    init: true\n    command: [\"sleep\", \"3600\"]\n", serviceName))
	}

	ret := filepath.Join(directory, "compose.yaml")
	err := os.WriteFile(ret, []byte(sb.String()), 0644)
	require.NoError(t, err)

	return ret
}

func Test_UpByComposeYamlPath_InvalidInput(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil commandExecutor", func(t *testing.T) {
		err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, nil, "compose.yaml", &dockercomposeoptions.UpOptions{})
		require.Error(t, err)
	})

	t.Run("empty yamlPath", func(t *testing.T) {
		err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, commandexecutorexecoo.Exec(), "", &dockercomposeoptions.UpOptions{})
		require.Error(t, err)
	})

	t.Run("nil options", func(t *testing.T) {
		err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, commandexecutorexecoo.Exec(), "compose.yaml", nil)
		require.Error(t, err)
	})
}

func Test_ListOrphanServiceNamesByComposeYamlPath_InvalidInput(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil commandExecutor", func(t *testing.T) {
		_, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, nil, "compose.yaml")
		require.Error(t, err)
	})

	t.Run("empty yamlPath", func(t *testing.T) {
		_, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
	})
}

func Test_UpByComposeYamlPath_ComposeAvailable(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_docker compose is available", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-available-"+impl.Name, "sleeper")

			isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
			require.NoError(t, err)
			require.False(t, isUp)
		})
	}
}

func Test_UpByComposeYamlPath(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_starts services and is idempotent", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-up-"+impl.Name, "sleeper")

			// Ensure precondition: nothing running.
			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})

			isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
			require.NoError(t, err)
			require.False(t, isUp)

			for i := 0; i < 2; i++ {
				err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{})
				require.NoError(t, err)

				isUp, err = commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
				require.NoError(t, err)
				require.True(t, isUp)
			}
		})
	}
}

func Test_UpByComposeYamlPath_RemoveOrphans(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_removes orphans only if requested", func(t *testing.T) {
			projectName := "asciichtest-orphans-" + impl.Name

			// The full compose file contains all services ever used in this test.
			// Using it for Down ensures no container is left behind, even if the test fails.
			fullYamlPath := writeComposeFile(t, t.TempDir(), projectName, "keep", "orphan")

			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.DownOptions{})

			err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.UpOptions{})
			require.NoError(t, err)

			// Same project but the 'orphan' service is no longer defined.
			reducedYamlPath := writeComposeFile(t, t.TempDir(), projectName, "keep")

			t.Run("orphan is detected", func(t *testing.T) {
				orphans, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
				require.NoError(t, err)
				require.Equal(t, []string{"orphan"}, orphans)
			})

			t.Run("up without RemoveOrphans keeps orphan", func(t *testing.T) {
				err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, reducedYamlPath, &dockercomposeoptions.UpOptions{})
				require.NoError(t, err)

				orphans, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
				require.NoError(t, err)
				require.Equal(t, []string{"orphan"}, orphans)
			})

			t.Run("up with RemoveOrphans removes orphan and is idempotent", func(t *testing.T) {
				for i := 0; i < 2; i++ {
					err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, reducedYamlPath, &dockercomposeoptions.UpOptions{RemoveOrphans: true})
					require.NoError(t, err)

					orphans, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
					require.NoError(t, err)
					require.Empty(t, orphans)

					isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
					require.NoError(t, err)
					require.True(t, isUp)
				}
			})
		})
	}
}

func Test_DownByComposeYamlPath(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_stops services and is idempotent", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-down-"+impl.Name, "sleeper")
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})

			err := commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{})
			require.NoError(t, err)

			for i := 0; i < 2; i++ {
				err = commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})
				require.NoError(t, err)

				isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
				require.NoError(t, err)
				require.False(t, isUp)
			}
		})
	}
}
