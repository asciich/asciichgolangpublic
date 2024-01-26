# asciichgolangpublic

Helper functions for writing infrastructure related CLIs easier and faster.

## Logging

To provide easy readable CLI output its recommended to use the provided logging functions:

```golang
package main

import "github.com/asciich/asciichgolangpublic"

func main() {
	asciichgolangpublic.LogInfo("Shown without additional color.")
	asciichgolangpublic.LogInfof("Shown without additional color. %s", "Also available with formatting.")

	asciichgolangpublic.LogGood("Good messages are green.")
	asciichgolangpublic.LogGoodf("Good messages are green. %s", "Also available with formatting.")

	asciichgolangpublic.LogChanged("Changes are purple.")
	asciichgolangpublic.LogChangedf("Changes are purple. %s", "Also available with formatting.")

	asciichgolangpublic.LogWarn("Warnings are yellow.")
	asciichgolangpublic.LogWarnf("Warnings are yellow. %s", "Also available with formatting.")

	asciichgolangpublic.LogError("Errors are red.")
	asciichgolangpublic.LogErrorf("Errors are red. %s", "Also available with formatting.")

	asciichgolangpublic.LogFatalf("Fatal will exit with a red error message and exit code %d", 1)
}
```

Output produced by this example code:

![](docs/log_example.png)

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

