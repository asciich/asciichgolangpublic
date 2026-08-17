# nativehost

## Implementation

- `NewNativeHost()` returns a new `NativeHost` representing the `localhost`.
- The `NativeHost` must fulfill the `hostinterfaces.Host` implementation.
- File and directory related functions like `GetDirectoryByPath` must return a struct of the `nativefilesoo` package. This has to be validated by unittests in this package.
