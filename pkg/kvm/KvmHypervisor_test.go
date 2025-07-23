package kvm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestKvmHypervisorGetHostNameWhenUsingLocalhost(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				hypervisor := MustGetKvmHypervisorOnLocalhost()
				require.EqualValues(
					"localhost_connection",
					hypervisor.MustGetHostName(),
				)
			},
		)
	}
}
