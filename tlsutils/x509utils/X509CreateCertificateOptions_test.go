package x509utils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tlsutils/x509utils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateCertOptions_GetValidityDuration(t *testing.T) {
	t.Run("default duration", func(t *testing.T) {
		options := &x509utils.X509CreateCertificateOptions{}
		duration, err := options.GetValidityDurationAsString()
		require.NoError(t, err)
		require.EqualValues(t, "1months15d", duration)
	})
}

func Test_GetPrivateKeySizeOrDefaultIfUnset(t *testing.T) {
	ctx := getCtx()
	t.Run("key size unset", func(t *testing.T) {
		options := &x509utils.X509CreateCertificateOptions{}
		require.EqualValues(t, 4096, options.GetPrivateKeySizeOrDefaultIfUnset(ctx))
	})

	t.Run("key size set to 1024", func(t *testing.T) {
		options := &x509utils.X509CreateCertificateOptions{
			PrivateKeySize: 1024,
		}
		require.EqualValues(t, 1024, options.GetPrivateKeySizeOrDefaultIfUnset(ctx))
	})

	t.Run("key size set to 4096", func(t *testing.T) {
		options := &x509utils.X509CreateCertificateOptions{
			PrivateKeySize: 4096,
		}
		require.EqualValues(t, 4096, options.GetPrivateKeySizeOrDefaultIfUnset(ctx))
	})
}
