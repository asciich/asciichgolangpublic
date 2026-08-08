# natfivefiles package

Handle local files and directories in a go native way.

## Implementation Details

This package uses only Go standard library functions and avoids external command execution (no `exec` calls). For example:
- `IsStaticallyLinkedBinary` uses Go's `debug/elf` package to parse ELF binaries instead of calling the external `file` command.

## Examples

* [Copy file](Example_Copy_test.go)
    * [Copy file as root using sudo](Example_CopySudo_test.go)
* [Create directory recursively](Example_CreateDirectoryRecursively_test.go)
* [Create file and it's parent directories recursively](Example_CreateFileRecursively_test.go)
* [Move file](Example_Move_test.go)
    * [Move file as root using sudo](Example_MoveSudo_test.go)

* [Check if statically linked binary](Example_IsStaticallyLinkedBinary_test.go)
## For developers

To run tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/filesutils/nativefiles/..."
```
