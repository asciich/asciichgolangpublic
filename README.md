# asciichgolangpublic

Helper functions for writing infrastructure related CLIs easier and faster.

## Errors

It's recommended to use `TracedError` whenever an error occurs with a custom error message.
Error wrapping by directly passing errors or using the `%w` format string in `TracedErrorf` is supported.
TracedErrors give you a nice debug output including the stack trace in a human readable form compatiple to VSCode (affected sources can directly be opened from Terminal).

Example usage:
```golang
func inThisFunctionSomethingGoesWrong() (err error) {
    return asciichgolangpublic.TracedError("This is an error message") // Use TracedErrors when an error occures.
}

err = inThisFunctionSomethingGoesWrong()
asciichgolangpublic.Errors().IsTracedError(err) // returns true for all TracedErrors.
asciichgolangpublic.Errors().IsTracedError(fmt.Errorf("another error")) // returns false for all non TracedErrors.

err.Error() // includes the error message and the stack trace as human readable text.
```

