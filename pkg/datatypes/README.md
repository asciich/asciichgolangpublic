# asciichgolangpublic/datatypes

Additional convenience functions to handle datatypes in go.

## Coding style rules

* Since handling low level datatypes there should be as less dependencies as possible.
* Do not use logging since we are not handling high level stuff here.
    * The only exception to this rule are the `MustXXX` functions: Use `log.Panic(err)` to print the error and stack trace and abort the execution.
* Do not use `TracedErrors` to act in the same way as the std go library does and avoid unneeded dependencies.
