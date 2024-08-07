package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				credential := MustGetGopassCredentialByName(tt.name)
				assert.EqualValues(tt.name, credential.MustGetName())
			},
		)
	}
}
