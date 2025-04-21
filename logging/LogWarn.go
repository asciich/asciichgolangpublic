package logging

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/shell/terminalcolors"
)

var overrideFunctionLogWarn func(logmessage string)

func OverrideLogWarn(overrideFunction func(logmessage string)) {
	overrideFunctionLogWarn = overrideFunction
}

func LogWarn(logmessage string) {
	if overrideFunctionLogWarn != nil {
		overrideFunctionLogWarn(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeYellow() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

var overrideFunctionLogWarnf func(logmessage string, args ...interface{})

func OverrideLogWarnf(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogWarnf = overrideFunction
}

func LogWarnf(logmessage string, args ...interface{}) {
	if overrideFunctionLogWarnf != nil {
		overrideFunctionLogWarnf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogWarn(message)
}
