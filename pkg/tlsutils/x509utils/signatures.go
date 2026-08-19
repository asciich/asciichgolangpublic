package x509utils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func CheckCertificateChainString(ctx context.Context, toCheck string) error {
	return genericx509utils.CheckCertificateChainString(ctx, toCheck)
}
