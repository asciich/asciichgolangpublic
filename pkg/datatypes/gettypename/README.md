# gettypename package

The `gettypename` package contains the separate function to get the name of a data type because of:
- This is needed since the functionality is used in different places to avoid cyclic imports.
- Avoid a dependency to tracederrors.

## Specifications

For specifications see [gettypename.spec.md](gettypename.spec.md)
