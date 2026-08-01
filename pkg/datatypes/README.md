# asciichgolangpublic/datatypes

Additional convenience functions to handle datatypes in go.

## Subpackages

* [bigintutils](./bigintutils/): Handle big integers.
* [gettypename](./gettypename/): Get type names from values.
* [pointerutils](./pointerutils/): Helper functions to work with pointers.

## Coding style rules

* Since handling low level datatypes there should be as less dependencies as possible.
* Do not use logging since we are not handling high level stuff here.
    * The only exception to this rule are the `MustXXX` functions: Use `log.Panic(err)` to print the error and stack trace and abort the execution.
* Use `TracedErrors` to give a clear stack trace in the error message to simplify debugging. [This is automatically validated.](./decisions_UseTracedErrors_test.go)
