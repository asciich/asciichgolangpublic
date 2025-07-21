package logging

import "log"

var overrideFunctionLog func(logmessage string)

func OverrideLog(overrideFunction func(logmessage string)) {
	overrideFunctionLog = overrideFunction
}

func Log(logmessage string) {
	if overrideFunctionLog != nil {
		overrideFunctionLog(logmessage)
		return
	}

	log.Println(logmessage)

	for _, l := range globalLoggers {
		l.Println(logmessage)
	}
}
