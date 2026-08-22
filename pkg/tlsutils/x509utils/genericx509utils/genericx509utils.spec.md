# genericx509utils specifications

## Implementation:

- The `X509CertKeyPair` is used to hold the private key and the corresponding certificate.
    - At least these functions must be implemented:
        - `keyPair.WriteCertificatePemToFile(ctx context.Context, toWrite filesinterfaces.File) error`
        - `keyPair.WritePrivateKeyToFile(ctx context.Context, toWrite filesinterfaces.File) error`

## Testing

- All functions in this package must be validated by a unit test.
