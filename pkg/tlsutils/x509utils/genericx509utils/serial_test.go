package genericx509utils_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func Test_GenerateCertificateSerialNumber(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("returns non nil serial number", func(t *testing.T) {
		serial, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)
		require.NotNil(t, serial)
	})

	t.Run("serial number is greater than or equal to minimum", func(t *testing.T) {
		serial, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)

		minNumber := big.NewInt(256 * 256 * 256)
		require.True(t, serial.Cmp(minNumber) >= 0, "serial number %s should be >= %s", serial.String(), minNumber.String())
	})

	t.Run("serial number is less than maximum", func(t *testing.T) {
		serial, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)

		maxNumber := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
		require.True(t, serial.Cmp(maxNumber) < 0, "serial number %s should be < %s", serial.String(), maxNumber.String())
	})

	t.Run("serial numbers are unique", func(t *testing.T) {
		serial1, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)

		serial2, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)

		require.NotEqual(t, 0, serial1.Cmp(serial2), "two generated serial numbers should not be equal")
	})

	t.Run("serial number is positive", func(t *testing.T) {
		serial, err := genericx509utils.GenerateCertificateSerialNumber(ctx)
		require.NoError(t, err)

		require.True(t, serial.Sign() > 0, "serial number should be positive")
	})
}

func Test_GenerateCertificateSerialNumberAsString(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("returns non empty string", func(t *testing.T) {
		serialStr, err := genericx509utils.GenerateCertificateSerialNumberAsString(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, serialStr)
	})

	t.Run("returned string is a valid decimal number", func(t *testing.T) {
		serialStr, err := genericx509utils.GenerateCertificateSerialNumberAsString(ctx)
		require.NoError(t, err)

		parsed, ok := new(big.Int).SetString(serialStr, 10)
		require.True(t, ok, "serial number string '%s' should be a valid decimal number", serialStr)
		require.NotNil(t, parsed)
		require.True(t, parsed.Sign() > 0, "parsed serial number should be positive")
	})

	t.Run("returned string represents a number within valid range", func(t *testing.T) {
		serialStr, err := genericx509utils.GenerateCertificateSerialNumberAsString(ctx)
		require.NoError(t, err)

		parsed, ok := new(big.Int).SetString(serialStr, 10)
		require.True(t, ok)

		minNumber := big.NewInt(256 * 256 * 256)
		maxNumber := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

		require.True(t, parsed.Cmp(minNumber) >= 0, "serial number %s should be >= %s", parsed.String(), minNumber.String())
		require.True(t, parsed.Cmp(maxNumber) < 0, "serial number %s should be < %s", parsed.String(), maxNumber.String())
	})

	t.Run("serial number strings are unique", func(t *testing.T) {
		serial1, err := genericx509utils.GenerateCertificateSerialNumberAsString(ctx)
		require.NoError(t, err)

		serial2, err := genericx509utils.GenerateCertificateSerialNumberAsString(ctx)
		require.NoError(t, err)

		require.NotEqual(t, serial1, serial2, "two generated serial number strings should not be equal")
	})
}
