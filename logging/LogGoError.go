package logging

import "github.com/asciich/asciichgolangpublic/tracederrors"

var overrideFunctionLogGoError func(err error)

func OverrideLogGoError(overrideFunction func(err error)) {
	overrideFunctionLogGoError = overrideFunction
}

func LogGoError(err error) {
	if overrideFunctionLogGoError != nil {
		overrideFunctionLogGoError(err)
		return
	}

	LogError(err.Error())
}

var overrideFunctionLogGoErrorFatal func(err error)

func OverrideFunctionLogGoErrorFatal(overrideFunction func(err error)) {
	overrideFunctionLogGoErrorFatal = overrideFunction
}

func LogGoErrorFatal(err error) {
	if overrideFunctionLogGoErrorFatal != nil {
		overrideFunctionLogGoErrorFatal(err)
		return
	}

	LogFatal(err.Error())
}

var overrideFunctionLogGoErrorFatalWithTrace func(err error)

func OverrideLogGoErrorFatalWithTrace(overrideFunction func(err error)) {
	overrideFunctionLogGoErrorFatalWithTrace = overrideFunction
}

func LogGoErrorFatalWithTrace(err error) {
	if overrideFunctionLogGoErrorFatalWithTrace != nil {
		overrideFunctionLogGoErrorFatalWithTrace(err)
		return
	}

	LogGoErrorFatal(tracederrors.TracedErrorf("%v", err))
}
