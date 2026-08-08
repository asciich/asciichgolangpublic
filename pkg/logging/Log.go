package logging

import (
	"context"
	"log"
)

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

// LogByCtx logs a message using context. If the context contains a LogRecorder,
// the message will be captured in addition to being emitted normally.
func LogByCtx(ctx context.Context, logmessage string) {
	recorder := getLogRecorderFromCtx(ctx)
	if recorder != nil {
		recorder.Write([]byte(logmessage + "\n"))
	}

	Log(logmessage)
}
