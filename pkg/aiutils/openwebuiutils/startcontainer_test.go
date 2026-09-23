package openwebuiutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOpenWebUI(t *testing.T) {
	t.Run("valid url", func(t *testing.T) {
		openwebui, err := NewOpenWebUI("http://localhost:8080")
		require.NoError(t, err)
		require.NotNil(t, openwebui)
		require.Equal(t, "http://localhost:8080", openwebui.Url)
	})

	t.Run("invalid url", func(t *testing.T) {
		openwebui, err := NewOpenWebUI("not-a-url")
		require.Error(t, err)
		require.Nil(t, openwebui)
	})

	t.Run("empty url", func(t *testing.T) {
		openwebui, err := NewOpenWebUI("")
		require.Error(t, err)
		require.Nil(t, openwebui)
	})
}

func TestOpenWebUISetUrl(t *testing.T) {
	t.Run("set valid url", func(t *testing.T) {
		openwebui := &OpenWebUI{}
		err := openwebui.SetUrl("http://localhost:3000")
		require.NoError(t, err)
		require.Equal(t, "http://localhost:3000", openwebui.Url)
	})

	t.Run("set invalid url", func(t *testing.T) {
		openwebui := &OpenWebUI{}
		err := openwebui.SetUrl("invalid")
		require.Error(t, err)
	})
}

func TestOpenWebUIGetUrl(t *testing.T) {
	t.Run("get url when set", func(t *testing.T) {
		openwebui := &OpenWebUI{Url: "http://localhost:8080"}
		url, err := openwebui.GetUrl()
		require.NoError(t, err)
		require.Equal(t, "http://localhost:8080", url)
	})

	t.Run("get url when not set", func(t *testing.T) {
		openwebui := &OpenWebUI{}
		url, err := openwebui.GetUrl()
		require.Error(t, err)
		require.Empty(t, url)
	})
}

func TestStartOpenWebUIContainerOptions(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test-openwebui",
			Port:          8080,
			DataPath:      "/tmp/openwebui-data",
		}

		name, err := opts.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "test-openwebui", name)

		port, err := opts.GetPort()
		require.NoError(t, err)
		require.Equal(t, 8080, port)
	})

	t.Run("missing container name", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			Port: 8080,
		}

		_, err := opts.GetContainerName()
		require.Error(t, err)
	})

	t.Run("invalid port", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test",
			Port:          0,
		}

		_, err := opts.GetPort()
		require.Error(t, err)
	})

	t.Run("bind address localhost", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName:            "test",
			Port:                     8080,
			ReachableByOtherMachines: false,
		}

		addr, err := opts.GetBindAddress()
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1:8080", addr)
	})

	t.Run("bind address 0.0.0.0", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName:            "test",
			Port:                     8080,
			ReachableByOtherMachines: true,
		}

		addr, err := opts.GetBindAddress()
		require.NoError(t, err)
		require.Equal(t, "0.0.0.0:8080", addr)
	})

	t.Run("get ollama base url", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test",
			Port:          8080,
			OllamaBaseUrl: "http://localhost:11434",
		}

		require.Equal(t, "http://localhost:11434", opts.GetOllamaBaseUrl())
	})

	t.Run("get additional env vars", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test",
			Port:          8080,
			AdditionalEnvVars: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		}

		envVars := opts.GetAdditionalEnvVars()
		require.Len(t, envVars, 2)
		require.Equal(t, "value1", envVars["KEY1"])
	})

	t.Run("get additional env vars when nil", func(t *testing.T) {
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test",
			Port:          8080,
		}

		envVars := opts.GetAdditionalEnvVars()
		require.NotNil(t, envVars)
		require.Empty(t, envVars)
	})
}

func TestStartOpenWebUIAsDockerContainer(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		ctx := context.Background()
		container, err := StartOpenWebUIAsDockerContainer(ctx, nil)
		require.Error(t, err)
		require.Nil(t, container)
	})

	t.Run("missing container name", func(t *testing.T) {
		ctx := context.Background()
		opts := &StartOpenWebUIContainerOptions{
			Port: 8080,
		}
		container, err := StartOpenWebUIAsDockerContainer(ctx, opts)
		require.Error(t, err)
		require.Nil(t, container)
	})

	t.Run("missing port", func(t *testing.T) {
		ctx := context.Background()
		opts := &StartOpenWebUIContainerOptions{
			ContainerName: "test",
		}
		container, err := StartOpenWebUIAsDockerContainer(ctx, opts)
		require.Error(t, err)
		require.Nil(t, container)
	})

	// Note: Full integration test would require Docker to be running
	// and would actually start a container, so we skip it in unit tests
	t.Skip("Integration test requires Docker")
}
