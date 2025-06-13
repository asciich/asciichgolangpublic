package kubernetesimplementationindependend_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
)

func Test_GetResourcePlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"secret", "secrets"},
		{"Secret", "secrets"},
		{"secrets", "secrets"},
		{"Secrets", "secrets"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			plural, err := kubernetesimplementationindependend.GetResourcePlural(tt.input)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, plural)
		})
	}
}
