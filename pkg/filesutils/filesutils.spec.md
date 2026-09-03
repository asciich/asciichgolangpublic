# filesutils specifications

This are the specifications for the [`filesutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The `filesinterfaces.Directory` must implement at least this functions:
    - The `GetFileInfoOfFilesInDirectory(ctx context.Context, options *parametersoptions.ListFileOptions) ([]*FileInfo, error)` must return the `fileinfo.FileInfo` for every file in the directory:
        - To respect all options `parametersoptions.ListFileOptions` provides this function must reuse an already existing `List...` function.
        - To avoid code duplication this function should be added to the base class in the `filesgeneric` package.

## Testing

- For all exported functions of the `filesinterfaces` there must be tests for all available implementations in the `filesutils` package.
