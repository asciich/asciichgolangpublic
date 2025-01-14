package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/shell/terminalcolors"
)

var globalLogSettings LogSettings
var globalLoggers []*log.Logger

func EnableLoggingToUsersHome(applicationName string, verbose bool) (logFilePath string, err error) {
	return "", errors.TracedErrorNotImplemented()
	/* TODO enable again
	applicationName = strings.TrimSpace(applicationName)

	if applicationName == "" {
		return nil, errors.TracedErrorEmptyString("applicationName")
	}

	homeDir, err := Users().GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	logsDir, err := homeDir.GetSubDirectory("logs")
	if err != nil {
		return nil, err
	}

	err = logsDir.Create(verbose)
	if err != nil {
		return nil, err
	}

	applicationLogsDir, err := logsDir.GetSubDirectory(applicationName)
	if err != nil {
		return nil, err
	}

	err = applicationLogsDir.Create(verbose)
	if err != nil {
		return nil, err
	}

	logFileName := Time().GetCurrentTimeAsSortableString() + ".log"

	logFile, err = applicationLogsDir.GetFileInDirectory(logFileName)
	if err != nil {
		return nil, err
	}

	logFilePath, err := logFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		return nil, err
	}

	loggerToAdd := log.New(file, "", log.LstdFlags|log.Lshortfile)

	globalLoggers = append(globalLoggers, loggerToAdd)

	if verbose {
		logging.LogInfof("All logs are now written to the log file '%s'.", logFilePath)
	}

	return logFile, nil
	*/
}

func Log(logmessage string) {
	log.Println(logmessage)

	for _, l := range globalLoggers {
		l.Println(logmessage)
	}
}

func LogBold(logmessage string) {
	Log(logmessage)
}

func LogByChangeSummary(changeSummary *changesummary.ChangeSummary, message string) {
	isChanged := false

	if changeSummary != nil {
		isChanged = changeSummary.IsChanged()
	}

	if isChanged {
		LogChanged(message)
	} else {
		LogInfo(message)
	}
}

func LogByChangeSummaryf(changeSummary *changesummary.ChangeSummary, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)

	LogByChangeSummary(changeSummary, formattedMessage)
}

func LogChanged(logmessage string) {
	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeMangenta() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

func LogChangedf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogChanged(message)
}

func LogError(logmessage string) {
	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeRed() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

func LogErrorf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogError(message)
}

func LogFatal(logmessage string) {
	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeRed() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
	os.Exit(1)
}

func LogFatalWithTrace(errorMessageOrError interface{}) {
	LogGoErrorFatal(errors.TracedError(errorMessageOrError))
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
	LogGoErrorFatal(errors.TracedErrorf("%v", err))
}

func LogGood(logmessage string) {
	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeGreen() + logmessage + terminalcolors.GetCodeNoColor()
	}
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

func LogTurnOfColorOutput() {
	globalLogSettings.SetColorEnabled(false)
}

func LogTurnOnColorOutput() {
	globalLogSettings.SetColorEnabled(true)
}

func LogWarn(logmessage string) {
	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeYellow() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

func LogWarnf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogWarn(message)
}

func MustEnableLoggingToUsersHome(applicationName string, verbose bool) (logFilePath string) {
	LogFatalWithTrace("NotImplemented")
	/*
		logFile, err := EnableLoggingToUsersHome(applicationName, verbose)
		if err != nil {
			logging.LogGoErrorFatal(err)
		}

		return logFile
	*/
	return ""
}
