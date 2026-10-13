package hostreconciler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

const testComposeContent = "services:\n  hello:\n    image: nginx:alpine\n"

// fakeSource is a Source returning fixed values for testing.
type fakeSource struct {
	beforeScript          []string
	hostname              string
	networkInterfaces     []NetworkInterface
	dockerComposeProjects []DockerComposeProject
}

func (f *fakeSource) GetBeforeScript(ctx context.Context) ([]string, error) {
	return f.beforeScript, nil
}

func (f *fakeSource) GetHostName(ctx context.Context) (string, error) {
	return f.hostname, nil
}

func (f *fakeSource) GetNetworkInterfaces(ctx context.Context) ([]NetworkInterface, error) {
	return f.networkInterfaces, nil
}

func (f *fakeSource) GetDockerComposeProjects(ctx context.Context) ([]DockerComposeProject, error) {
	return f.dockerComposeProjects, nil
}

// writeTestFile ensures the parent directory of 'path' exists and writes 'content' to 'path'.
func writeTestFile(t *testing.T, ctx context.Context, path string, content string) {
	t.Helper()

	require.NoError(t, nativefiles.CreateDirectory(ctx, filepath.Dir(path), &filesoptions.CreateOptions{}))
	require.NoError(t, nativefiles.WriteString(ctx, path, content))
}

func Test_hostConfig_Parse(t *testing.T) {
	t.Run("full config is parsed", func(t *testing.T) {
		content := `---
before_script:
  - echo hello world
hostname: host.example.com
network:
  interfaces:
    - name: enp0
      layer2_bridge: br0
      deactivate_offloading: true
    - name: enp1
      deactivate_offloading: true
    - name: enp2
      layer2_bridge: br2
`
		config := hostConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		require.Equal(t, []string{"echo hello world"}, config.BeforeScript)
		require.Equal(t, "host.example.com", config.Hostname)
		require.Equal(
			t,
			[]NetworkInterface{
				{Name: "enp0", Layer2Bridge: "br0", DeactivateOffloading: true},
				{Name: "enp1", DeactivateOffloading: true},
				{Name: "enp2", Layer2Bridge: "br2"},
			},
			config.Network.Interfaces,
		)
	})

	t.Run("multi line before_script entry is parsed", func(t *testing.T) {
		content := `---
before_script:
  - |
    echo first
    echo second
hostname: host.example.com
`
		config := hostConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		require.Equal(t, []string{"echo first\necho second\n"}, config.BeforeScript)
	})

	t.Run("before_script and network section are optional", func(t *testing.T) {
		config := hostConfig{}
		err := yaml.Unmarshal([]byte("---\nhostname: host.example.com\n"), &config)
		require.NoError(t, err)

		require.Empty(t, config.BeforeScript)
		require.Empty(t, config.Network.Interfaces)
		require.NoError(t, validateBeforeScript(config.BeforeScript))
		require.NoError(t, validateNetworkInterfaces(config.Network.Interfaces))
	})

	t.Run("deactivate_offloading defaults to false", func(t *testing.T) {
		content := `---
interfaces:
  - name: enp0
    layer2_bridge: br0
`
		config := networkConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		require.Len(t, config.Interfaces, 1)
		require.False(t, config.Interfaces[0].DeactivateOffloading)
	})
}

func Test_validateBeforeScript(t *testing.T) {
	t.Run("valid entries", func(t *testing.T) {
		require.NoError(t, validateBeforeScript([]string{"echo hello world", "true"}))
	})

	t.Run("nil is valid", func(t *testing.T) {
		require.NoError(t, validateBeforeScript(nil))
	})

	t.Run("empty entry", func(t *testing.T) {
		require.Error(t, validateBeforeScript([]string{"echo hello world", ""}))
	})

	t.Run("whitespace only entry", func(t *testing.T) {
		require.Error(t, validateBeforeScript([]string{"  \n\t"}))
	})
}

func Test_runBeforeScript(t *testing.T) {
	ctx := context.TODO()

	t.Run("no before_script is a no-op", func(t *testing.T) {
		require.NoError(t, runBeforeScript(ctx, nil))
	})

	t.Run("successful entries", func(t *testing.T) {
		require.NoError(t, runBeforeScript(ctx, []string{"echo hello world", "true"}))
	})

	t.Run("failing entry returns error", func(t *testing.T) {
		require.Error(t, runBeforeScript(ctx, []string{"echo hello world", "exit 1"}))
	})

	t.Run("multi line entry stops at first failing command", func(t *testing.T) {
		require.Error(t, runBeforeScript(ctx, []string{"false\necho not reached"}))
	})

	t.Run("failing command in pipe returns error", func(t *testing.T) {
		require.Error(t, runBeforeScript(ctx, []string{"false | cat"}))
	})

	t.Run("invalid entry is rejected before execution", func(t *testing.T) {
		require.Error(t, runBeforeScript(ctx, []string{""}))
	})
}

func Test_validateNetworkInterfaces(t *testing.T) {
	t.Run("valid network interfaces", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0", DeactivateOffloading: true},
			{Name: "enp1", DeactivateOffloading: true},
			{Name: "enp2", Layer2Bridge: "br2"},
		}
		require.NoError(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("nil is valid", func(t *testing.T) {
		require.NoError(t, validateNetworkInterfaces(nil))
	})

	t.Run("network interface without actions is valid", func(t *testing.T) {
		require.NoError(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0"}}))
	})

	t.Run("settings in separate list items", func(t *testing.T) {
		content := `---
interfaces:
  - name: enp0
  - layer2_bridge: br0
`
		config := networkConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		err = validateNetworkInterfaces(config.Interfaces)
		require.Error(t, err)
		require.Contains(t, err.Error(), "same list item")
	})

	t.Run("missing name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{DeactivateOffloading: true}}))
	})

	t.Run("network interface defined twice", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", DeactivateOffloading: true},
			{Name: "enp0", Layer2Bridge: "br0"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("bridge used by two network interfaces", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0"},
			{Name: "enp1", Layer2Bridge: "br0"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("nested bridges", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0"},
			{Name: "br0", Layer2Bridge: "br1"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("same name for bridge and network interface", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "br0", Layer2Bridge: "br0"}}))
	})

	t.Run("shell unsafe network interface name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0';reboot'", DeactivateOffloading: true}}))
	})

	t.Run("shell unsafe bridge name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0", Layer2Bridge: "br0';reboot'"}}))
	})
}

func Test_reconcileNetworkInterfaces(t *testing.T) {
	ctx := context.TODO()

	t.Run("no network interfaces is a no-op without root", func(t *testing.T) {
		require.NoError(t, reconcileNetworkInterfaces(ctx, nil))
	})

	t.Run("network interface without actions is a no-op without root", func(t *testing.T) {
		require.NoError(t, reconcileNetworkInterfaces(ctx, []NetworkInterface{{Name: "enp0"}}))
	})

	t.Run("invalid network interfaces are rejected before any change", func(t *testing.T) {
		require.Error(t, reconcileNetworkInterfaces(ctx, []NetworkInterface{{Layer2Bridge: "br0"}}))
	})
}

func Test_isDockerComposeFileName(t *testing.T) {
	t.Run("compose.yaml", func(t *testing.T) {
		require.True(t, isDockerComposeFileName("compose.yaml"))
	})

	t.Run("docker-compose.yml", func(t *testing.T) {
		require.True(t, isDockerComposeFileName("docker-compose.yml"))
	})

	t.Run("custom yaml file name", func(t *testing.T) {
		require.False(t, isDockerComposeFileName("my_compose.yaml"))
	})
}

func Test_parseDockerComposeTopLevelName(t *testing.T) {
	t.Run("name is set", func(t *testing.T) {
		name, err := parseDockerComposeTopLevelName("name: mydeployment\n" + testComposeContent)
		require.NoError(t, err)
		require.Equal(t, "mydeployment", name)
	})

	t.Run("name is not set", func(t *testing.T) {
		name, err := parseDockerComposeTopLevelName(testComposeContent)
		require.NoError(t, err)
		require.Equal(t, "", name)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		_, err := parseDockerComposeTopLevelName("name: [unclosed\n")
		require.Error(t, err)
	})
}

func Test_validateDockerComposeProjects(t *testing.T) {
	t.Run("valid projects", func(t *testing.T) {
		projects := []DockerComposeProject{
			{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/compose.yaml"},
			{Name: "beta", ComposeFilePath: "/config/docker_compose/beta/docker-compose.yml"},
		}
		require.NoError(t, validateDockerComposeProjects(projects))
	})

	t.Run("nil is valid", func(t *testing.T) {
		require.NoError(t, validateDockerComposeProjects(nil))
	})

	t.Run("empty name", func(t *testing.T) {
		require.Error(t, validateDockerComposeProjects([]DockerComposeProject{{ComposeFilePath: "/config/docker_compose/alpha/compose.yaml"}}))
	})

	t.Run("invalid name", func(t *testing.T) {
		require.Error(t, validateDockerComposeProjects([]DockerComposeProject{{Name: "MyApp", ComposeFilePath: "/config/docker_compose/MyApp/compose.yaml"}}))
	})

	t.Run("relative path", func(t *testing.T) {
		require.Error(t, validateDockerComposeProjects([]DockerComposeProject{{Name: "alpha", ComposeFilePath: "docker_compose/alpha/compose.yaml"}}))
	})

	t.Run("unsupported compose file name", func(t *testing.T) {
		require.Error(t, validateDockerComposeProjects([]DockerComposeProject{{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/app.yaml"}}))
	})

	t.Run("directory name differs from project name", func(t *testing.T) {
		require.Error(t, validateDockerComposeProjects([]DockerComposeProject{{Name: "alpha", ComposeFilePath: "/config/docker_compose/beta/compose.yaml"}}))
	})

	t.Run("project defined twice", func(t *testing.T) {
		projects := []DockerComposeProject{
			{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/compose.yaml"},
			{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/compose.yml"},
		}
		err := validateDockerComposeProjects(projects)
		require.Error(t, err)
		require.Contains(t, err.Error(), "more than once")
	})
}

func Test_DirectorySource_GetDockerComposeProjects(t *testing.T) {
	ctx := context.TODO()

	t.Run("missing docker_compose directory returns empty list", func(t *testing.T) {
		source := &DirectorySource{dirPath: t.TempDir()}

		projects, err := source.GetDockerComposeProjects(ctx)
		require.NoError(t, err)
		require.Empty(t, projects)
	})

	t.Run("deployments are returned in lexical order and other entries are skipped", func(t *testing.T) {
		configDir := t.TempDir()
		composeDir := filepath.Join(configDir, DockerComposeDirName)

		writeTestFile(t, ctx, filepath.Join(composeDir, "beta", "compose.yaml"), testComposeContent)
		writeTestFile(t, ctx, filepath.Join(composeDir, "alpha", "docker-compose.yml"), testComposeContent)
		writeTestFile(t, ctx, filepath.Join(composeDir, ".hidden", "compose.yaml"), testComposeContent)
		writeTestFile(t, ctx, filepath.Join(composeDir, "README.md"), "# docs\n")

		source := &DirectorySource{dirPath: configDir}

		projects, err := source.GetDockerComposeProjects(ctx)
		require.NoError(t, err)

		absComposeDir, err := filepath.Abs(composeDir)
		require.NoError(t, err)

		require.Equal(
			t,
			[]DockerComposeProject{
				{Name: "alpha", ComposeFilePath: filepath.Join(absComposeDir, "alpha", "docker-compose.yml")},
				{Name: "beta", ComposeFilePath: filepath.Join(absComposeDir, "beta", "compose.yaml")},
			},
			projects,
		)
		require.NoError(t, validateDockerComposeProjects(projects))
	})

	t.Run("matching top level name is accepted", func(t *testing.T) {
		configDir := t.TempDir()
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "alpha", "compose.yaml"), "name: alpha\n"+testComposeContent)

		source := &DirectorySource{dirPath: configDir}

		projects, err := source.GetDockerComposeProjects(ctx)
		require.NoError(t, err)
		require.Len(t, projects, 1)
		require.Equal(t, "alpha", projects[0].Name)
	})

	t.Run("different top level name is rejected", func(t *testing.T) {
		configDir := t.TempDir()
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "alpha", "compose.yaml"), "name: other\n"+testComposeContent)

		source := &DirectorySource{dirPath: configDir}

		_, err := source.GetDockerComposeProjects(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "differs")
	})

	t.Run("deployment without compose file is rejected", func(t *testing.T) {
		configDir := t.TempDir()
		require.NoError(t, nativefiles.CreateDirectory(ctx, filepath.Join(configDir, DockerComposeDirName, "alpha"), &filesoptions.CreateOptions{}))

		source := &DirectorySource{dirPath: configDir}

		_, err := source.GetDockerComposeProjects(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "No compose file")
	})

	t.Run("deployment with multiple compose files is rejected", func(t *testing.T) {
		configDir := t.TempDir()
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "alpha", "compose.yaml"), testComposeContent)
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "alpha", "compose.yml"), testComposeContent)

		source := &DirectorySource{dirPath: configDir}

		_, err := source.GetDockerComposeProjects(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Multiple compose files")
	})

	t.Run("compose file directly in docker_compose is rejected", func(t *testing.T) {
		configDir := t.TempDir()
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "my_compose.yaml"), testComposeContent)

		source := &DirectorySource{dirPath: configDir}

		_, err := source.GetDockerComposeProjects(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "own directory")
	})

	t.Run("invalid deployment directory name is rejected", func(t *testing.T) {
		configDir := t.TempDir()
		writeTestFile(t, ctx, filepath.Join(configDir, DockerComposeDirName, "MyApp", "compose.yaml"), testComposeContent)

		source := &DirectorySource{dirPath: configDir}

		_, err := source.GetDockerComposeProjects(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "MyApp")
	})
}

func Test_reconcileDockerComposeProjects(t *testing.T) {
	ctx := context.TODO()

	t.Run("no docker compose projects is a no-op without docker", func(t *testing.T) {
		require.NoError(t, reconcileDockerComposeProjects(ctx, nil))
	})

	t.Run("invalid docker compose projects are rejected before docker is used", func(t *testing.T) {
		require.Error(t, reconcileDockerComposeProjects(ctx, []DockerComposeProject{{Name: "alpha", ComposeFilePath: "relative/compose.yaml"}}))
	})
}

func Test_ReconcileOnce(t *testing.T) {
	ctx := context.TODO()

	t.Run("nil source", func(t *testing.T) {
		require.Error(t, ReconcileOnce(ctx, nil))
	})

	t.Run("empty hostname", func(t *testing.T) {
		require.Error(t, ReconcileOnce(ctx, &fakeSource{}))
	})

	t.Run("invalid network interfaces are rejected before the hostname is changed", func(t *testing.T) {
		source := &fakeSource{
			hostname:          "host.example.com",
			networkInterfaces: []NetworkInterface{{Layer2Bridge: "br0"}},
		}

		err := ReconcileOnce(ctx, source)
		require.Error(t, err)

		// The validation error proves ReconcileOnce stopped before calling SetHostName.
		require.Contains(t, err.Error(), "same list item")
	})

	t.Run("empty before_script entry is rejected before anything is executed", func(t *testing.T) {
		source := &fakeSource{
			beforeScript: []string{""},
			hostname:     "host.example.com",
		}

		err := ReconcileOnce(ctx, source)
		require.Error(t, err)
		require.Contains(t, err.Error(), "'before_script' is empty")
	})

	t.Run("failing before_script stops before the hostname is changed", func(t *testing.T) {
		source := &fakeSource{
			beforeScript: []string{"exit 1"},
			hostname:     "host.example.com",
		}

		err := ReconcileOnce(ctx, source)
		require.Error(t, err)

		// The before_script error proves ReconcileOnce stopped before calling SetHostName.
		require.Contains(t, err.Error(), "before_script entry '0' failed")
	})

	t.Run("invalid docker compose projects are rejected before anything is executed", func(t *testing.T) {
		source := &fakeSource{
			// Would fail if executed, which proves validation happens before the before_script.
			beforeScript: []string{"exit 1"},
			hostname:     "host.example.com",
			dockerComposeProjects: []DockerComposeProject{
				{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/compose.yaml"},
				{Name: "alpha", ComposeFilePath: "/config/docker_compose/alpha/compose.yml"},
			},
		}

		err := ReconcileOnce(ctx, source)
		require.Error(t, err)
		require.Contains(t, err.Error(), "more than once")
	})
}
