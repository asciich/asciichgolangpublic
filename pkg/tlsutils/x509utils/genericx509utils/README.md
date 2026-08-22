# genericx509utils

Generic X509 certificate utilities providing a pure Go implementation for handling X509 certificates and key pairs.

## Types

### X509CertKeyPair

Holds a private key and its corresponding certificate.

**Fields:**
- `Cert *x509.Certificate` - The X509 certificate
- `Key crypto.PrivateKey` - The private key

**Methods:**
- `GetCertificateAsPEMString() (string, error)` - Returns the certificate as a PEM-encoded string
- `GetCertificateAsPEMBytes() ([]byte, error)` - Returns the certificate as PEM-encoded bytes
- `GetX509Certificate() (*x509.Certificate, error)` - Returns a deep copy of the certificate
- `GetPrivateKeyAsPEMString() (string, error)` - Returns the private key as a PEM-encoded string
- `GetPrivateKey() (crypto.PrivateKey, error)` - Returns the private key
- `GetPublicKey() (crypto.PublicKey, error)` - Extracts the public key from the private key
- `IsKeyMatchingCert() (bool, error)` - Checks if the private key matches the certificate
- `CheckKeyMatchingCertificate() error` - Returns an error if the key doesn't match the certificate
- `WriteCertificatePemToFile(ctx context.Context, toWrite filesinterfaces.File) error` - Writes the certificate PEM to a file
- `WritePrivateKeyToFile(ctx context.Context, toWrite filesinterfaces.File) error` - Writes the private key PEM to a file

## Testing

Run tests with:
```bash
go test -v ./...
```

## Specifications

For specifications see [genericx509utils.spec.md](genericx509utils.spec.md)
