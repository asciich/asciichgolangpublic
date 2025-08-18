# natfivefiles package

Handle local files and directories in a go native way.

## Examples

* [Create directory recursively](Example_CreateDirectoryRecursively_test.go)
* [Create file and it's parent directories recursively](Example_CreateFileRecursively_test.go)

## For developers

To run tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/filesutils/nativefiles/..."
```