# x509utils specifications

This are the specifications for the [`x509utils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- For functions interacting with files or network there are two implementations, both in a dedicated subpackage. Every function must implemented in both so the functionality stays the same, regardless if executed locally or on a remote machine:
    - [`commandexecutorx509utils`](commandexecutorx509utils/README.md):
        - Using `commandexecutor` in combination with shell commands to implement the logic.
        - This also works on other machines over SSH for example.
    - [`nativex509utils`](nativex509utils/README.md):
        - Use native golang to handle the certificates.
        - Works only on the local machine.
    - Furthermore, for convenience every function implemented in these packages is as well added to this `x509utils`, calling the [`nativex509utils`](nativex509utils/README.md) implementation.
    - For every convenience function add a well commented `Example_<function>_test.go` containing a test for documentation purposes how to use the function
- Generic functions like `GetCommonName(cert *x509.Certificate) (string, error)` which do not require network or file access are implemented in the subpackage [`genericx509utils`](genericx509utils/README.md).
- For the inmemory representation of a certificate use `*x509.Certificate`.
- The `ReadCertificateFromFile(ctx context.Context, pathToRead string) (*x509.Certificate, error)` reads a file from disk.
- Do not shorten `Certificate` to `Cert` in function names. Instead of `IsKeyMatchingCert` use `IsKeyMatchingCertificate`.

## Tests

- Test in this directory are used to test all implementations behave the same way. This is why they require a for loop which runs the same test for each implementation.
