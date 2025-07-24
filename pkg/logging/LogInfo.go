package logging

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
)

var overrideFunctionLogInfo func(logmessage string)

func OverrideLogInfo(overrideFunction func(logmessage string)) {
	overrideFunctionLogInfo = overrideFunction
}

func LogInfo(logmessage string) {
	if overrideFunctionLogInfo != nil {
		overrideFunctionLogInfo(logmessage)
		return
	}

	Log(logmessage)
}

var overrideFunctionLogInfof func(logmessage string, args ...interface{})

func OverrideLogInfof(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogInfof = overrideFunction
}

func LogInfof(logmessage string, args ...interface{}) {
	if overrideFunctionLogInfof != nil {
		overrideFunctionLogInfof(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogInfo(message)
}

var overrideFunctionLogInfoByCtx func(ctx context.Context, logmessage string)

func OverrideLogInfoByCtx(overrideFunction func(ctx context.Context, logmessage string)) {
	overrideFunctionLogInfoByCtx = overrideFunction
}

func LogInfoByCtx(ctx context.Context, logmessage string) {
	if overrideFunctionLogInfoByCtx != nil {
		overrideFunctionLogInfoByCtx(ctx, logmessage)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	logLinePrefix := contextutils.GetLogLinePrefixFromCtx(ctx)
	if logLinePrefix == "" {
		LogInfo(logmessage)
	} else {
		LogInfoWithLinePrefix(logmessage, logLinePrefix)
	}
}

var overrideFunctionLogInfoByCtxf func(ctx context.Context, logmessage string, args ...interface{})

func OverrideLogInfoByCtxf(overrideFunction func(ctx context.Context, logmessage string, args ...interface{})) {
	overrideFunctionLogInfoByCtxf = overrideFunction
}

func LogInfoByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	if overrideFunctionLogInfoByCtxf != nil {
		overrideFunctionLogInfoByCtxf(ctx, logmessage, args...)
		return
	}

	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	LogInfof(logmessage, args...)
}

var overrideFunctionLogInfoWithLinePrefix func(logmessage string, logLinePrefix string)

func OverrideLogInfoWithLinePrefix(overrideFunction func(logmessage string, logLinePrefix string)) {
	overrideFunctionLogInfoWithLinePrefix = overrideFunction
}

func LogInfoWithLinePrefix(logmessage string, logLinePrefix string) {
	if overrideFunctionLogInfoWithLinePrefix != nil {
		overrideFunctionLogInfoWithLinePrefix(logmessage, logLinePrefix)
		return
	}

	LogInfo(stringsutils.AddLinePrefix(logmessage, logLinePrefix))
}
