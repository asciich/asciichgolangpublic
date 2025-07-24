# tracederrors package

Provides errors with stacktrace and other additional information to help debugging after crashes where only logs are available.


Error wrapping by directly passing errors or using the `%w` format string in `TracedErrorf` is supported.
TracedErrors give you a nice debug output including the stack trace in a human readable form compatiple to VSCode (affected sources can directly be opened from Terminal).

## Examples

* [Example usage](./Example_usage_test.go)

## For developers

To run the tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/tracederrors/..."
```
