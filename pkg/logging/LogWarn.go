package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/terminalcolors"
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
		logmessage = terminalcolors.CODE_YELLOW + logmessage + terminalcolors.CODE_NO_COLOR
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

var overrideFunctionLogWarnByCtxf func(ctx context.Context, logmessage string, args ...interface{})

func LogWarnByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	if overrideFunctionLogWarnByCtxf != nil {
		overrideFunctionLogWarnByCtxf(ctx, logmessage, args...)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

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

	LogWarnf(logmessage, args...)
}
