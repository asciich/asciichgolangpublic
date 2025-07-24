# filesutils package

Contains various implementations to work with files:
* [tempfileoo](./tempfilesoo/): Create temporary files and directories in a object oriented way.

## For developers

To run the tests of filesutils use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/filesutils/..."
```