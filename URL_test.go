package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				url := MustGetUrlFromString(tt.url)
				fqdn, path := url.MustGetFqdnWitShemeAndPathAsString()

				assert.EqualValues(tt.expectedFqdn, fqdn)
				assert.EqualValues(tt.expectedPath, path)
			},
		)
	}
}
