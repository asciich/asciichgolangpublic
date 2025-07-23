# commandexectuorexec package

Bash implementation to run commands inside a bash.

For the object oriented Bash implementation see [commandexecutorbashoo](/pkg/commandexecutor/commandexecutorbashoo/)

## Examples

* [Run a simple command (echo hello world)](./Example_echoHelloWorld_test.go)

## For developers

To run the tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/commandexecutor/commandexecutorbash/..."
```
