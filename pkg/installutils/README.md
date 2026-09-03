# installutils package

Can be used to install binaries from various sources.

## Subpackages

* [commandexecutorinstall](./commandexecutorinstall/README.md): Install binaries using command executor (works remotely via e.g. SSH).
* [nativeinstall](./nativeinstall/README.md): Install binaries using native Go libraries (local only).

## Functions

The parent package provides convenience functions that delegate to [`nativeinstall`](nativeinstall/README.md):
- `Install(ctx, options)`: Install from URL or local path

For remote installation, use [`commandexecutorinstall.Install`](commandexecutorinstall/README.md).

## Specifications

For specifications see [installutils.spec.md](installutils.spec.md)

## Examples

Examples are documented in the subpackage READMEs.
