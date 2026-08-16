# nativehost

## Implementation

- `NewNativeHost()` returns a new `NativeHost` representing the `localhost`.
- The `NativeHost` must fulfill the `hostinterfaces.Host` implementation.
- File and directory related functions like ``etDirectoryByPath` must return a struct of the `nativefilesoo` package. This has to be validated by unittests in this package.
