package asciichgolangpublic

import (
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestKvmHypervisorGetHostnameWhenUsingLocalhost(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				hypervisor := MustGetKvmHypervisorOnLocalhost()
				assert.EqualValues(
					"localhost_connection",
					hypervisor.MustGetHostName(),
				)
			},
		)
	}
}
