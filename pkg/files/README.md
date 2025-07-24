# files package

The goal is to move everyithing from this pacakge to the [filesutils](/pkg/filesutils/)

## For developers

To run the tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/files/..."
```