# commandexecutorbashoo package

Object oriented bash [commandexecutor](/pkg/commandexecutor/) implementation.

For the non object oriented implementations see [commandexecutorbash](/pkg/commandexecutor/commandexecutorexec/).

## Examples

* [Run a simple command (echo hello world)](./Example_echoHelloWorld_test.go)

## For developers

To run the tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/commandexecutor/commandexecutorbashoo/..."
```
