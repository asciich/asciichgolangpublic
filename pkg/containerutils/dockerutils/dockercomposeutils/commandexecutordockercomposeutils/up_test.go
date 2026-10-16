package commandexecutordockercomposeutils_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

// getContainerIDs returns the sorted IDs of the running containers of the compose file 'yamlPath'.
// If 'serviceNames' are given only the containers of these services are returned.
//
// Calls docker directly (independent of the code under test) so recreation can be verified:
// a recreated container always gets a new ID.
func getContainerIDs(t *testing.T, yamlPath string, serviceNames ...string) []string {
	t.Helper()

	args := append([]string{"compose", "-f", yamlPath, "ps", "--quiet"}, serviceNames...)
	output, err := exec.Command("docker", args...).Output()
	require.NoError(t, err)

	ids := strings.Fields(string(output))
	sort.Strings(ids)

	return ids
}

// requireAllRecreated checks that none of the IDs in 'after' is in 'before'.
func requireAllRecreated(t *testing.T, before []string, after []string) {
	t.Helper()

	require.Len(t, after, len(before))
	for _, id := range after {
		require.NotContains(t, before, id, "container '%s' was not recreated", id)
	}
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
}

func Test_UpByComposeYamlPath_NilOptions(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_nil options behaves like default options", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-niloptions-"+impl.Name, "sleeper")

			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})

			for i := 0; i < 2; i++ {
				err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, nil)
				require.NoError(t, err)

				isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
				require.NoError(t, err)
				require.True(t, isUp)
			}
		})
	}
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

			var idsFirstRun []string
			for i := 0; i < 2; i++ {
				err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{})
				require.NoError(t, err)

				isUp, err = commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
				require.NoError(t, err)
				require.True(t, isUp)

				// Idempotency: the second run must not touch the containers.
				ids := getContainerIDs(t, yamlPath)
				require.Len(t, ids, 1)
				if i == 0 {
					idsFirstRun = ids
				} else {
					require.Equal(t, idsFirstRun, ids)
				}
			}
		})
	}
}

func Test_UpByComposeYamlPath_ForceRecreate(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	for _, impl := range getCommandExecutorImplementations() {
		impl := impl

		t.Run(impl.Name+"_recreates already running services on every call", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-forcerecreate-running-"+impl.Name, "sleeper1", "sleeper2")

			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})

			err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{})
			require.NoError(t, err)

			idsBefore := getContainerIDs(t, yamlPath)
			require.Len(t, idsBefore, 2)

			// ForceRecreate is intentionally NOT idempotent: every call creates new containers.
			for i := 0; i < 2; i++ {
				err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{ForceRecreate: true})
				require.NoError(t, err)

				isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
				require.NoError(t, err)
				require.True(t, isUp)

				idsAfter := getContainerIDs(t, yamlPath)
				requireAllRecreated(t, idsBefore, idsAfter)

				idsBefore = idsAfter
			}
		})

		t.Run(impl.Name+"_starts services which are not running", func(t *testing.T) {
			yamlPath := writeComposeFile(t, t.TempDir(), "asciichtest-forcerecreate-stopped-"+impl.Name, "sleeper")

			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.DownOptions{})

			isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
			require.NoError(t, err)
			require.False(t, isUp)

			err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, yamlPath, &dockercomposeoptions.UpOptions{ForceRecreate: true})
			require.NoError(t, err)

			isUp, err = commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, yamlPath)
			require.NoError(t, err)
			require.True(t, isUp)
			require.Len(t, getContainerIDs(t, yamlPath), 1)
		})

		t.Run(impl.Name+"_combined with RemoveOrphans", func(t *testing.T) {
			projectName := "asciichtest-forcerecreate-orphans-" + impl.Name

			fullYamlPath := writeComposeFile(t, t.TempDir(), projectName, "keep", "orphan")

			err := commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.DownOptions{})
			require.NoError(t, err)
			defer commandexecutordockercomposeutils.DownByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.DownOptions{})

			err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, fullYamlPath, &dockercomposeoptions.UpOptions{})
			require.NoError(t, err)

			idsKeepBefore := getContainerIDs(t, fullYamlPath, "keep")
			require.Len(t, idsKeepBefore, 1)

			reducedYamlPath := writeComposeFile(t, t.TempDir(), projectName, "keep")

			orphans, err := commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
			require.NoError(t, err)
			require.Equal(t, []string{"orphan"}, orphans)

			err = commandexecutordockercomposeutils.UpByComposeYamlPath(ctx, impl.Executor, reducedYamlPath, &dockercomposeoptions.UpOptions{
				ForceRecreate: true,
				RemoveOrphans: true,
			})
			require.NoError(t, err)

			// Orphan is removed ...
			orphans, err = commandexecutordockercomposeutils.ListOrphanServiceNamesByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
			require.NoError(t, err)
			require.Empty(t, orphans)

			// ... and the remaining service was recreated.
			isUp, err := commandexecutordockercomposeutils.IsUpByComposeYamlPath(ctx, impl.Executor, reducedYamlPath)
			require.NoError(t, err)
			require.True(t, isUp)

			requireAllRecreated(t, idsKeepBefore, getContainerIDs(t, reducedYamlPath, "keep"))
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
