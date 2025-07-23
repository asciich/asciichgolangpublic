package parameteroptions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGopassSecretOptions_GetPath(t *testing.T) {
	tests := []struct {
		path         string
		expectedPath string
	}{
		{"this/is/my/path", "this/is/my/path"},
		{"this/is/my/path2", "this/is/my/path2"},
		{"this/is/my/path_2", "this/is/my/path_2"},
		{"this/is/my_my/path_2", "this/is/my_my/path_2"},

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

				secretPath, err := secretOptions.GetSecretPath()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedPath, secretPath)
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
		{"/this/is/my/base-Name", "base-Name", "abc", "this/is/my/abc"},
		{"/this/is/my/base-Name", "base-Name", "abc.key", "this/is/my/abc.key"},
		{"/this/is/my/base_Name", "base_Name", "abc", "this/is/my/abc"},
		{"/this/is/my/base_Name", "base_Name", "abc.key", "this/is/my/abc.key"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				secretOptions := &GopassSecretOptions{
					SecretPath: tt.path,
				}

				baseName, err := secretOptions.GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedBaseName, baseName)

				err = secretOptions.SetBaseName(tt.newBaseName)
				require.NoError(t, err)

				baseName, err = secretOptions.GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, tt.newBaseName, baseName)

				path, err := secretOptions.GetSecretPath()
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedPathWithNewBaseName, path)
			},
		)
	}
}
