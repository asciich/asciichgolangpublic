package kubernetesimplementationindependend_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesimplementationindependend"
)

func Test_GetObjectPlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"secret", "secrets"},
		{"Secret", "secrets"},
		{"secrets", "secrets"},
		{"Secrets", "secrets"},
		{"GitRepository", "gitrepositories"},
		{"gitrepositories", "gitrepositories"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			plural, err := kubernetesimplementationindependend.GetObjectPlural(tt.input)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, plural)
		})
	}
}
