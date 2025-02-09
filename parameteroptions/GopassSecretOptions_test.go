package parameteroptions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGopassSecretOptions_GetPath(t *testing.T) {
	tests := []struct {
		path         string
		expectedPath string
	}{
		{"this/is/my/path", "this/is/my/path"},

		// Leading slashes "/" are removed automatically.
		// This is expected behaviour since gopass does not work with leading slashes "/"
		{"/this/is/my/path", "this/is/my/path"},
		{"//this/is/my/path", "this/is/my/path"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				secretOptions := &GopassSecretOptions{
					SecretPath: tt.path,
				}

				require.EqualValues(
					t,
					tt.expectedPath,
					secretOptions.MustGetSecretPath(),
				)
			},
		)
	}
}

func TestGopassSecretOptions_SetAndGetBaseName(t *testing.T) {
	tests := []struct {
		path                        string
		expectedBaseName            string
		newBaseName                 string
		expectedPathWithNewBaseName string
	}{
		{"this/is/my/path", "path", "abc", "this/is/my/abc"},
		{"/this/is/my/path", "path", "abc", "this/is/my/abc"},
		{"/this/is/my/path", "path", "baseName", "this/is/my/baseName"},
		{"/this/is/my/baseName", "baseName", "abc", "this/is/my/abc"},
		{"/this/is/my/baseName", "baseName", "abc.key", "this/is/my/abc.key"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				secretOptions := &GopassSecretOptions{
					SecretPath: tt.path,
				}

				require.EqualValues(
					tt.expectedBaseName,
					secretOptions.MustGetBaseName(),
				)

				secretOptions.MustSetBaseName(tt.newBaseName)

				require.EqualValues(
					tt.newBaseName,
					secretOptions.MustGetBaseName(),
				)

				require.EqualValues(
					tt.expectedPathWithNewBaseName,
					secretOptions.MustGetSecretPath(),
				)

			},
		)
	}
}
