# nativehost specifications

This are the specifications for the [`nativehost` package](README.md).

This document extends the parent specifications [hostsutils.spec.md](../hostsutils.spec.md).

## Implementation

- `NewNativeHost()` returns a new `NativeHost` representing the `localhost`.
- The `NativeHost` must fulfill the `hostinterfaces.Host` implementation.
- File and directory related functions like `GetDirectoryByPath` must return a struct of the [`nativefilesoo`](../../filesutils/nativefilesoo/README.md) package. This has to be validated by unittests in this package.
