package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/terminalcolors"
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
		logmessage = terminalcolors.CODE_MANGENTA + logmessage + terminalcolors.CODE_NO_COLOR
	}
	Log(logmessage)
}

var overrideFunctionLogChangedf func(logmessage string, arg ...interface{})

func OverrideLogChangedf(overrideFunction func(logmessage string, arg ...interface{})) {
	overrideFunctionLogChangedf = overrideFunction
}

func LogChangedByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	verbose := contextutils.GetVerboseFromContext(ctx)

	contextutils.SetChangeIndicator(ctx, true)

	if verbose {
		recorder := getLogRecorderFromCtx(ctx)
		var formattedMessage string
		if len(args) > 0 {
			formattedMessage = fmt.Sprintf(logmessage, args...)
		} else {
			formattedMessage = logmessage
		}

		if recorder != nil {
			recorder.Write([]byte(formattedMessage + "\n"))
		}

		LogChangedf(logmessage, args...)
	}
}

func LogChangedByCtx(ctx context.Context, logmessage string) {
	verbose := contextutils.GetVerboseFromContext(ctx)

	contextutils.SetChangeIndicator(ctx, true)

	if verbose {
		recorder := getLogRecorderFromCtx(ctx)
		if recorder != nil {
			recorder.Write([]byte(logmessage + "\n"))
		}

		LogChanged(logmessage)
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
