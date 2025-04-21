package logging

var overrideFunctionLogBold func(logmessage string)

func OverrideLogBold(overrideFunction func(logmessage string)) {
	overrideFunctionLogBold = overrideFunction
}

func LogBold(logmessage string) {
	if overrideFunctionLogBold != nil {
		overrideFunctionLogBold(logmessage)
		return
	}

	Log(logmessage)
}
