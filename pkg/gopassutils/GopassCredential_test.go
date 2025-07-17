package gopassutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/gopassutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGopassCredentialSetAndGetName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"a"},
		{"a/b"},
		{"a/c"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				credential, err := gopassutils.GetGopassCredentialByName(tt.name)
				require.NoError(t, err)
				name, err := credential.GetName()
				require.NoError(t, err)
				require.EqualValues(t, tt.name, name)
			},
		)
	}
}
