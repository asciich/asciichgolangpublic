package genericx509utils

import "crypto/x509"

func GetX509CertificateDeepCopy(in *x509.Certificate) (out *x509.Certificate) {
	if in == nil {
		return nil
	}

	out = new(x509.Certificate)
	*out = *in

	return out
}
