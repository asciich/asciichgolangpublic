# genericx509utils specifications

This are the specifications for the [`genericx509utils` package](README.md).

This document extends the parent specifications [x509utils.spec.md](../x509utils.spec.md).

## Implementation:

- The `X509CertKeyPair` is used to hold the private key and the corresponding certificate.
    - At least these functions must be implemented:
        - `keyPair.WriteCertificatePemToFile(ctx context.Context, toWrite filesinterfaces.File) error`
        - `keyPair.WritePrivateKeyToFile(ctx context.Context, toWrite filesinterfaces.File) error`

## Testing

- All functions in this package must be validated by a unit test.
