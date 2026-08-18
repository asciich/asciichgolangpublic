package genericx509utils

import (
	"context"
	"math/big"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bigintutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GenerateCertificateSerialNumberAsString(ctx context.Context) (string, error) {
	serial, err := GenerateCertificateSerialNumber(ctx)
	if err != nil {
		return "", err
	}

	ret, err := bigintutils.ToDecimalString(serial)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func GenerateCertificateSerialNumber(ctx context.Context) (serialNumber *big.Int, err error) {
	logging.LogInfoByCtx(ctx, "Generate certificate serial number started.")

	minNumber := big.NewInt(256 * 256 * 256)
	maxNumber := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

	serialNumber, err = bigintutils.GetRandomBigIntByInts(minNumber, maxNumber)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	hexRepresentation, err := bigintutils.ToHexStringColonSeparated(serialNumber)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	logging.LogInfoByCtxf(ctx, "Generate certificate serial number finished. Generated serial number is '%s'.", hexRepresentation)

	return serialNumber, nil
}
