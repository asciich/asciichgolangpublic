package urlsutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
				require := require.New(t)

				url := MustGetUrlFromString(tt.url)
				fqdn, path := url.MustGetFqdnWitShemeAndPathAsString()

				require.EqualValues(tt.expectedFqdn, fqdn)
				require.EqualValues(tt.expectedPath, path)
			},
		)
	}
}
