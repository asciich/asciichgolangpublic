package kvm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic"
)

func TestKvmHypervisorGetHostNameWhenUsingLocalhost(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
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
