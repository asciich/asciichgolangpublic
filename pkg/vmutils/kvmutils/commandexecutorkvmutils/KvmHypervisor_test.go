package commandexecutorkvmutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/commandexecutorkvmutils"
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
				hypervisor, err := commandexecutorkvmutils.GetKvmHypervisorOnLocalhost()
				require.NoError(t, err)
				hostname, err := hypervisor.GetHostName()
				require.EqualValues(t, "localhost_connection", hostname)
			},
		)
	}
}
