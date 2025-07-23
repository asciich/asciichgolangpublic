package logging

import (
	"context"
	"fmt"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/shell/terminalcolors"
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
		logmessage = terminalcolors.GetCodeGreen() + logmessage + terminalcolors.GetCodeNoColor()
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

	if contextutils.GetVerboseFromContext(ctx) {
		LogGoodf(logmessage, args...)
	}
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

	if contextutils.GetVerboseFromContext(ctx) {
		LogGood(logmessage)
	}
}

func LogGoodf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogGood(message)
}
