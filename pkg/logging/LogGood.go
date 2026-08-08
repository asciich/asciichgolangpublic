package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/terminalcolors"
)

var overrideFunctionLogGood func(logmessage string)

func OverrideLogGood(overrideFunction func(logmessage string)) {
	overrideFunctionLogGood = overrideFunction
}

func LogGood(logmessage string) {
	if overrideFunctionLogGood != nil {
		overrideFunctionLogGood(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.CODE_GREEN + logmessage + terminalcolors.CODE_NO_COLOR
	}
	LogInfo(logmessage)
}

var overrideFunctionLogGoodByCtxf func(ctx context.Context, logmessage string, args ...interface{})

func OverrideLogGoodByCtxf(overrideFunction func(ctx context.Context, logmessage string, args ...interface{})) {
	overrideFunctionLogGoodByCtxf = overrideFunction
}

func LogGoodByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	if overrideFunctionLogGoodByCtxf != nil {
		overrideFunctionLogGoodByCtxf(ctx, logmessage, args...)
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

	LogGoodf(logmessage, args...)
}

var overrideFunctionLogGoodByCtx func(ctx context.Context, logmessage string)

func OverrideLogGoodByCtx(overrideFunction func(ctx context.Context, logmessage string)) {
	overrideFunctionLogGoodByCtx = overrideFunction
}

func LogGoodByCtx(ctx context.Context, logmessage string) {
	if overrideFunctionLogGoodByCtx != nil {
		overrideFunctionLogGoodByCtx(ctx, logmessage)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	recorder := getLogRecorderFromCtx(ctx)
	if recorder != nil {
		recorder.Write([]byte(logmessage + "\n"))
	}

	LogGood(logmessage)
}

func LogGoodf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogGood(message)
}
