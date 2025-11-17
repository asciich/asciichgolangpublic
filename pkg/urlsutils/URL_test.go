package urlsutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

func TestUrlGetFqdnAndPath(t *testing.T) {
	tests := []struct {
		url          string
		expectedFqdn string
		expectedPath string
	}{
		{"https://gitlab.asciich.ch", "https://gitlab.asciich.ch", ""},
		{"https://gitlab.asciich.ch/", "https://gitlab.asciich.ch", ""},
		{"https://gitlab.asciich.ch/gitlab_management", "https://gitlab.asciich.ch", "gitlab_management"},
		{"https://gitlab.asciich.ch/gitlab_management/", "https://gitlab.asciich.ch", "gitlab_management"},
		{"https://gitlab.asciich.ch/gitlab_management/pre-commit", "https://gitlab.asciich.ch", "gitlab_management/pre-commit"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				url, err := urlsutils.GetUrlFromString(tt.url)
				require.NoError(t, err)

				fqdn, path, err := url.GetFqdnWitShemeAndPathAsString()
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedFqdn, fqdn)
				require.EqualValues(t, tt.expectedPath, path)
			},
		)
	}
}
