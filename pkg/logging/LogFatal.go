package logging

import (
	"fmt"
	"os"

	"github.com/asciich/asciichgolangpublic/shellutils/terminalcolors"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var overrideFunctionLogFatal func(logmessage string)

func OverrideLogFatal(overrideFunction func(logmessage string)) {
	overrideFunctionLogFatal = overrideFunction
}

func LogFatal(logmessage string) {
	if overrideFunctionLogFatal != nil {
		overrideFunctionLogFatal(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeRed() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
	os.Exit(1)
}

var overrideFunctionLogFatalWithTrace func(errorMessageOrError interface{})

func OverrideLogFatalWithTrace(overrideFunction func(errorMessageOrError interface{})) {
	overrideFunctionLogFatalWithTrace = overrideFunction
}

func LogFatalWithTrace(errorMessageOrError interface{}) {
	if overrideFunctionLogFatalWithTrace != nil {
		overrideFunctionLogFatalWithTrace(errorMessageOrError)
		return
	}
	LogGoErrorFatal(tracederrors.TracedError(errorMessageOrError))
}

var overrideFunctionLogFatalWithTracef func(logmessage string, args ...interface{})

func OverrideLogFatalWithTracef(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogFatalWithTracef = overrideFunction
}

func LogFatalWithTracef(logmessage string, args ...interface{}) {
	if overrideFunctionLogFatalWithTracef != nil {
		overrideFunctionLogFatalWithTracef(logmessage, args...)
	}

	message := fmt.Sprintf(logmessage, args...)
	LogFatalWithTrace(message)
}

var overrideFunctionLogFatalf func(logmessage string, args ...interface{})

func OverrideLogFatalf(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogFatalf = overrideFunction
}

func LogFatalf(logmessage string, args ...interface{}) {
	if overrideFunctionLogFatalf != nil {
		overrideFunctionLogFatalf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogFatal(message)
}
