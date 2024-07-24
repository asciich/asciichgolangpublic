package asciichgolangpublic

/* TODO enable again
import (
	"testing"

)

func TestHostCheckReachableBySsh(t *testing.T) {
	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"cerberus3.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				host := MustGetHostByHostname(tt.hostname)
				host.MustCheckReachableBySsh(verbose)
			},
		)
	}
}
*/