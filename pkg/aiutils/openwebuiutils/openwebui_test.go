package openwebuiutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenWebUIGetHealthStatus(t *testing.T) {
	t.Run("get health status", func(t *testing.T) {
		openwebui := &OpenWebUI{Url: "http://localhost:8080"}

		// This is a unit test - actual health check would require running service
		// We just test that the method exists and can be called
		ctx := context.Background()
		_, err := openwebui.GetHealthStatus(ctx)

		// Expect error since service is not running, but method should be callable
		assert.Error(t, err)
	})
}

func TestOpenWebUIGetVersion(t *testing.T) {
	t.Run("get version", func(t *testing.T) {
		openwebui := &OpenWebUI{Url: "http://localhost:8080"}

		ctx := context.Background()
		_, err := openwebui.GetVersion(ctx)

		// Expect error since service is not running, but method should be callable
		assert.Error(t, err)
	})
}
