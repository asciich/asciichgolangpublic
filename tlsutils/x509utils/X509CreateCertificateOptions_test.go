package x509utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/tlsutils/x509utils"
)

func Test_GetValidityDuration(t *testing.T) {
	t.Run("default duration", func(t *testing.T) {
		options := &x509utils.X509CreateCertificateOptions{}
		duration, err := options.GetValidityDurationAsString()
		require.NoError(t,err)
		require.EqualValues(t, "1months15d", duration)
	})
}
