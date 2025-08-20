# filesutils package

Contains various implementations to work with files:
* [nativefiles](./nativefiles/): Handle local files using go native/ std library commands.
* [tempfile](./tempfiles/): Create temporary files and directories.
* [tempfileoo](./tempfilesoo/): Create temporary files and directories in a object oriented way.

## Examples

* [Create directory recursively](./nativefiles/Example_CreateDirectoryRecursively_test.go)
* [Create file and it's parent directories recursively](./nativefiles/Example_CreateFileRecursively_test.go)

## For developers

To run the tests of filesutils use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/filesutils/..."
```
