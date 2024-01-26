package github.com/asciich/asciichgolangpublic

import (
	"fmt"
	"log"
	"os"
)

func Log(logmessage string) {
	log.Println(logmessage)
}

func LogBold(logmessage string) {
	Log(logmessage)
}

func LogChanged(logmessage string) {
	logmessage = TerminalColors().GetCodeMangenta() + logmessage + TerminalColors().GetCodeNoColor()
	Log(logmessage)
}

func LogChangedf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogChanged(message)
}

func LogError(logmessage string) {
	logmessage = TerminalColors().GetCodeRed() + logmessage + TerminalColors().GetCodeNoColor()
	Log(logmessage)
}

func LogErrorf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogError(message)
}

func LogFatal(logmessage string) {
	logmessage = TerminalColors().GetCodeRed() + logmessage + TerminalColors().GetCodeNoColor()
	Log(logmessage)
	os.Exit(1)
}

func LogFatalWithTrace(message string) {
	LogGoErrorFatal(TracedError(message))
}

func LogFatalWithTracef(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogFatalWithTrace(message)
}

func LogFatalf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogFatal(message)
}

func LogGoError(err error) {
	LogErrorf("%v", err)
}

func LogGoErrorFatal(err error) {
	LogFatalf(err.Error())
}

func LogGoErrorFatalWithTrace(err error) {
	LogGoErrorFatal(TracedErrorf("%v", err))
}

func LogGood(logmessage string) {
	logmessage = TerminalColors().GetCodeGreen() + logmessage + TerminalColors().GetCodeNoColor()
	LogInfo(logmessage)
}

func LogGoodf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogGood(message)
}

func LogInfo(logmessage string) {
	Log(logmessage)
}

func LogInfof(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogInfo(message)
}

func LogWarn(logmessage string) {
	logmessage = TerminalColors().GetCodeYellow() + logmessage + TerminalColors().GetCodeNoColor()
	Log(logmessage)
}

func LogWarnf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogWarn(message)
}
