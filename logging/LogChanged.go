package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/shell/terminalcolors"
)

var overrideFunctionLogChanged func(logmessage string)

func OverrideLogChanged(overrideFunction func(logmessage string)) {
	overrideFunctionLogChanged = overrideFunction
}

func LogChanged(logmessage string) {
	if overrideFunctionLogChanged != nil {
		overrideFunctionLogChanged(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeMangenta() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

var overrideFunctionLogChangedf func(logmessage string, arg ...interface{})

func OverrideLogChangedf(overrideFunction func(logmessage string, arg ...interface{})) {
	overrideFunctionLogChangedf = overrideFunction
}

func LogChangedByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	verbose := contextutils.GetVerboseFromContext(ctx)

	if verbose {
		LogChangedf(logmessage, args...)
	}
}

func LogChangedf(logmessage string, args ...interface{}) {
	if overrideFunctionLogChangedf != nil {
		overrideFunctionLogChangedf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogChanged(message)
}
