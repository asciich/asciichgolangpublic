# containerimagehandler specifications

This are the specifications for the [`containerimagehandler` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation:

- The `func PushImage(ctx context.Context, image v1.Image, nameAndTag string, options *containeroptions.PushImageOptions) (string, error)` pushes an image to a repository.
    - A registry can be specified either as part of `nameAndTag` or in `options.Registry`. If both are set, they must match, otherwise an error is returned.
    - `options.Username` / `options.Password` can be set to authenticate with basic auth. If `options.Username` is empty, the default keychain (e.g. local docker credential store) is used instead.
    - Return the full image name and tag of the uploaded container. Including registry if needed. It must be possible to run `docker run --rm -it <returned_string>` afterwards.
