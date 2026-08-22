# httputils specifications

This are the specifications for the [`httputils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The `httputils` package itself contains only convenience functions to using the `httpnativeclientoo`. The actual implementation has to be done in:
    - `httpnativeclientoo` using go native http libraries. This is used on local machines.
    - `commandexecutorclientoo` using exec commands and mostly `curl`. This is used mostly on jumphosts and other remote machines.
- SSL/TLS certificates:
    - The `httpoptions.RequestOptions` must contain an optional bool `CollectCertificates` to collect store the certificate used by the webserver.
        - If set the response must keep a copy of the certificate.
        - The `httputilsinterfaces.Response` must implement a `GetServerEndEntitiyCertificate(ctx context.Context) (*x509.Certificate, error)` and a `GetServerCertificateChain(ctx context.Context) ([]*x509.Certificate, error)` to access the copy of the certificate(s) returned by the server.
        - The `httputilsinterfaces.Response` must implement a `LogCertInfo(ctx context.Context) error` to log the servers certificate infos.
            - Print the serial number separated by `:` to make it more human readable.
        - This requires an `Exammple_GetRequestAndLogCertificateInfos_test.go` for documentation purposes.

## Testing

- Tests in the `httputils` package are ment to run the same logic against all implementations to ensure they behave in the same way.
    - Every function must be tested this way.
- `Example_*_test.go` packages are ment for documentation purposes and can use the convenience functions in `httputils` directly. If not asked explicitly they must not run against all implementations.
