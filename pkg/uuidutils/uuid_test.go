package uuidutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGenerate(t *testing.T) {
	generated := uuidutils.Generate(getCtx())
	require.Len(t, generated, len("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"))
}

func TestIsUuid(t *testing.T) {
	tests := []struct{
		name string
		input string
		expected bool
	}{
		{"empty", "", false},
		{"random", "abc123", false},
		{"k8s ns uid", "b913691f-eaf8-4ac0-a872-f1e2881c3ec9", true},
		{"generated", uuidutils.Generate(getCtx()), true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, uuidutils.IsUuid(tt.input))
		})
	}
}