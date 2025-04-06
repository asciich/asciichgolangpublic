# asciichgolangpublic

This module helps to write infrastructure and/or automation related microservices and CLIs easier and faster.
By providing a lot of convenience functions, sanity checks during runtime and detailed error messages it can be used to write easy to understand software to automate repeatable work.
The focus is on ease of use and developer speed instead of algorithm speed and computer resource efficiency. 

## Design choices, principles and background information

* Design choices:
	* Function which return an error must not panic.
	* Use `Set` and `Get` functions which allows to validate input and output on every access:
		* Especially when using the provided functions to quickly automatize some stuff validating all inputs as a first step in every function helps to avoid unwanted side effects.
	* While there were a lot of `MustABC` functions introduced this project moves towards removing them except in favor of using `mustutils.Must`
	* `CheckAbc` functions evaluate if `Abc` is given. If given `nil`, otherwise an error is returned.
	* As default there is no log output as silent CLI's are easier to handle if glued together e.g. in Bash scripts. While verbosity was initialy set by `verbose` bool this is projects moves towards replacing it in favor of using `context.Context`. See [`contextutils`](pkg/contextutils)
	* Short cuts and code hacks are not nice but still better than doing things by hand. They are at least a good starting point of what functionality is needed and can be improved over time.
* Releasing:
	* Release often: Every (small) improvement is an improvemnt and will be released as soon as possible.
	* This repository will never reach v1: There will be always be breaking changes if needed to improve the code.
	* Everytime the code base is touched it should look better than before.
* Readability:
	* An end user of this repository should be able to write readable code.
* Background information:
	* Currently this is a one man show.
	* It bases purely on some code I wrote at home in my free time used for automating my personal homelab.
* Multiple levels of automation implementation and where this library can help:
	1. Knowledge in the head of the develpers (worst case):
		- In worst case not even documented at all.
		- Very error prone and a huge truck factor.
		- **This library is no help here!**
	1. Documented instructions:
		- Idealy step-by-step instructions.
		- This approach is also called Wiki-Ops.
		- **This library is no help here!**
	1. Scripting (Bash, Python, Ansible...):
		- Better than any documentation since the steps are reproducible and complete (otherwise it would not run successfully).
		- Often hard to (unit-)test.
		- While Python is a fully fledged programming language able to handle complex things the amount of complexity and reuseabilty is in Bash is limited.
		- Bash scripting reflects the way system administrators work interactively with the systems and is therefore often easy to understand for other team members.
		- External interpreters and tools are needed.
		- Interpreted languages are error prone and often easy to exploid using code injection
		- **While orchestrating Bash commands or python one liners is not how programming works it is often the first starting point and still way better than wiki-ops. This library offers an easy to use interaction with the CLI using Bash() or CommandExecutors in general. But keep in mind: It's a starting point and sometimes needed to get things up and running in time but must be migrated towards native implementations on the long run.**
	1. High level programming languages enriched with convenince functions and shortcuts allowed:
		- Reusable code.
		- Easy to automated new task by combining existing convenience functions.
		- Checking of inputs at every function call already detects malformed input at an early state.
		- A lot of boilerplate code.
		- Unit- and integration tests possible and useful.
		- Still some external dependencies, especially when calling other binaries to achive a shortcut.
		- Putting convenience functions together often leads to inefficient algorithms (e.g. more API requests than needed.).
		- The programming style gives some guard rails which can make it easier for other system administrators to start implementing their own stuff.
		- Easier to debug crashes since full stack trace is provided in errors (see [section errors](#errors)) but can also lead to security issues by exposing internal information.
		- **Most of this library is written this way. Not bad as a starting point and usable on a high level so implementing new tasks is easy but still a lot of improvmentes towards an idiomatic golang codebase.**
	1. Idiomatic Golang code (best case):
		- Reusable code.
		- Unit- and integration tests possible and useful.
		- No external dependencies since everything is natively implemneted.
		- Fastest execution time.
		- **While this library is currently far away from idiomatic go the aim is to move towards idiomatic code in the implementation.**

## Logging

To provide easy readable CLI output its recommended to use the provided logging functions:

```golang
package main

import "github.com/asciich/asciichgolangpublic/logging"

func main() {
	logging.LogInfo("Shown without additional color.")
	logging.LogInfof("Shown without additional color. %s", "Also available with formatting.")

	logging.LogGood("Good messages are green.")
	logging.LogGoodf("Good messages are green. %s", "Also available with formatting.")

	logging.LogChanged("Changes are purple.")
	logging.LogChangedf("Changes are purple. %s", "Also available with formatting.")

	logging.LogWarn("Warnings are yellow.")
	logging.LogWarnf("Warnings are yellow. %s", "Also available with formatting.")

	logging.LogError("Errors are red.")
	logging.LogErrorf("Errors are red. %s", "Also available with formatting.")

	logging.LogFatalf("Fatal will exit with a red error message and exit code %d", 1)
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
    return tracederrors.TracedError("This is an error message") // Use TracedErrors when an error occures.
}

err = inThisFunctionSomethingGoesWrong()
tracederrors.Errors().IsTracedError(err) // returns true for all TracedErrors.
tracederrors.Errors().IsTracedError(fmt.Errorf("another error")) // returns false for all non TracedErrors.

err.Error() // includes the error message and the stack trace as human readable text.
```
