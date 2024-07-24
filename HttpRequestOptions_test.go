package asciichgolangpublic

/* TODO enable again
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpRequestOptionsGetOutputFilePath(t *testing.T) {
	tests := []struct {
		urlString        string
		expectedBasename string
	}{
		{"http://cerberus3.asciich.ch:8080/data/credentials/asciich_ssh_keys/authorized_keys.asc", "authorized_keys.asc"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				options := &HttpRequestOptions{
					URL: tt.urlString,
				}

				assert.EqualValues(tt.expectedBasename, options.MustGetOutputFilePath())
			},
		)
	}
}
*/
