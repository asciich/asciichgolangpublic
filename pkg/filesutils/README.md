# filesutils package

Contains various implementations to work with files:

## Subpackages

* [commandexecutorfile](./commandexecutorfile/README.md): File operations using command executor.
* [commandexecutorfileoo](./commandexecutorfileoo/README.md): File operations using command executor (object oriented).
* [commandexecutortempfile](./commandexecutortempfile/README.md): Temporary file operations using command executor.
* [commandexecutortempfilesoo](./commandexecutortempfilesoo/README.md): Temporary file operations using command executor (object oriented).
* [fileinfo](./fileinfo/README.md): File information utilities.
* [nativefiles](./nativefiles/README.md): Handle local files using go native/ std library commands.
* [nativefilesoo](./nativefilesoo/README.md): Object oriented native file operations.
* [tempfiles](./tempfiles/README.md): Create temporary files and directories.
* [tempfilesoo](./tempfilesoo/README.md): Create temporary files and directories in a object oriented way.
* [filesoptions](./filesoptions/README.md): File operation options.
* [filesinterfaces](./filesinterfaces/README.md): File interfaces.
* [filesgeneric](./filesgeneric/README.md): Generic file utilities (pure in-memory operations).

## Specifications

For specifications see [filesutils.spec.md](filesutils.spec.md)

## Examples

* [Copy file](./nativefiles/Example_Copy_test.go)
    * [Copy file as root using sudo](./nativefiles/Example_CopySudo_test.go)
* [Create directory recursively](./nativefiles/Example_CreateDirectoryRecursively_test.go)
* [Create file and it's parent directories recursively](./nativefiles/Example_CreateFileRecursively_test.go)
* [Get nativefileoo by path](Example_NewFileByPath_test.go)
* [Is file a statically linked binary](Example_IsStaticallyLinkedBinary_test.go)
* [Move file](./nativefiles/Example_Move_test.go)
    * [Move file as root using sudo](./nativefiles/Example_MoveSudo_test.go)


## For developers

To run the tests of filesutils use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/filesutils/..."
```
