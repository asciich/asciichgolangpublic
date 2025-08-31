package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/terminalcolors"
)

var overrideFunctionLogError func(logmessage string)

func OverrideLogError(overrideFunction func(logmessage string)) {
	overrideFunctionLogError = overrideFunction
}

func LogError(logmessage string) {
	if overrideFunctionLogError != nil {
		overrideFunctionLogError(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.CODE_RED + logmessage + terminalcolors.CODE_NO_COLOR
	}
	Log(logmessage)
}

var overrideFunctionLogErrorf func(logmessage string, arg ...interface{})

func OverrideLogErrorf(overrideFunction func(logmessage string, arg ...interface{})) {
	overrideFunctionLogErrorf = overrideFunction
}

func LogErrorf(logmessage string, args ...interface{}) {
	if overrideFunctionLogErrorf != nil {
		overrideFunctionLogErrorf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogError(message)
}

var overrideFunctionLogErrorByCtxf func(ctx context.Context, logmessage string, args ...interface{})

func OverrideLogErrorByCtxf(overrideFunction func(ctx context.Context, logmessage string, args ...interface{})) {
	overrideFunctionLogErrorByCtxf = overrideFunction
}

func LogErrorByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	if overrideFunctionLogErrorByCtxf != nil {
		overrideFunctionLogErrorByCtxf(ctx, logmessage, args...)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	LogErrorf(logmessage, args...)
}

var overrideFunctionLogErrorByCtx func(ctx context.Context, logmessage string)

func OverrideLogErrorByCtx(overrideFunction func(ctx context.Context, logmessage string)) {
	overrideFunctionLogErrorByCtx = overrideFunction
}

func LogErrorByCtx(ctx context.Context, logmessage string, ) {
	if overrideFunctionLogErrorByCtx != nil {
		overrideFunctionLogErrorByCtx(ctx, logmessage)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	LogError(logmessage)
}
