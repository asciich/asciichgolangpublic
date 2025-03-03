package logging

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/contextutils"
	"github.com/asciich/asciichgolangpublic/shell/terminalcolors"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var globalLogSettings LogSettings
var globalLoggers []*log.Logger

func EnableLoggingToUsersHome(applicationName string, verbose bool) (logFilePath string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
	/* TODO enable again
	applicationName = strings.TrimSpace(applicationName)

	if applicationName == "" {
		return nil, tracederrors.TracedErrorEmptyString("applicationName")
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

var overrideFunctionLogByChangeSummary func(changeSummary *changesummary.ChangeSummary, message string)

func OverrideLogByChangeSummary(overrideFunction func(changeSummary *changesummary.ChangeSummary, message string)) {
	overrideFunctionLogByChangeSummary = overrideFunction
}

func LogByChangeSummary(changeSummary *changesummary.ChangeSummary, message string) {
	if overrideFunctionLogByChangeSummary != nil {
		overrideFunctionLogByChangeSummary(changeSummary, message)
		return
	}

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

var overrideFunctionLogByChangeSummaryf func(changeSummary *changesummary.ChangeSummary, message string, args ...interface{})

func OverrideLogByChangeSummaryf(overrideFunction func(changeSummary *changesummary.ChangeSummary, message string, args ...interface{})) {
	overrideFunctionLogByChangeSummaryf = overrideFunction
}

func LogByChangeSummaryf(changeSummary *changesummary.ChangeSummary, message string, args ...interface{}) {
	if overrideFunctionLogByChangeSummaryf != nil {
		overrideFunctionLogByChangeSummaryf(changeSummary, message, args)
		return
	}

	formattedMessage := fmt.Sprintf(message, args...)

	LogByChangeSummary(changeSummary, formattedMessage)
}

var overrideFunctionLogChanged func(logmessage string)

func OverrideLogChanged(overrideFunction func(logmessage string)) {
	overrideFunctionLogChanged = overrideFunction
}

func LogChanged(logmessage string) {
	if overrideFunctionLogChanged != nil {
		overrideFunctionLogChanged(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeMangenta() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

var overrideFunctionLogChangedf func(logmessage string, arg ...interface{})

func OverrideLogChangedf(overrideFunction func(logmessage string, arg ...interface{})) {
	overrideFunctionLogChangedf = overrideFunction
}

func LogChangedByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	verbose := contextutils.GetVerboseFromContext(ctx)

	if verbose {
		LogChangedf(logmessage, args...)
	}
}

func LogChangedf(logmessage string, args ...interface{}) {
	if overrideFunctionLogChangedf != nil {
		overrideFunctionLogChangedf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogChanged(message)
}

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
		logmessage = terminalcolors.GetCodeRed() + logmessage + terminalcolors.GetCodeNoColor()
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

var overrideFunctionLogFatal func(logmessage string)

func OverrideLogFatal(overrideFunction func(logmessage string)) {
	overrideFunctionLogFatal = overrideFunction
}

func LogFatal(logmessage string) {
	if overrideFunctionLogFatal != nil {
		overrideFunctionLogFatal(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeRed() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
	os.Exit(1)
}

var overrideFunctionLogFatalWithTrace func(errorMessageOrError interface{})

func OverrideLogFatalWithTrace(overrideFunction func(errorMessageOrError interface{})) {
	overrideFunctionLogFatalWithTrace = overrideFunction
}

func LogFatalWithTrace(errorMessageOrError interface{}) {
	if overrideFunctionLogFatalWithTrace != nil {
		overrideFunctionLogFatalWithTrace(errorMessageOrError)
		return
	}
	LogGoErrorFatal(tracederrors.TracedError(errorMessageOrError))
}

var overrideFunctionLogFatalWithTracef func(logmessage string, args ...interface{})

func OverrideLogFatalWithTracef(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogFatalWithTracef = overrideFunction
}

func LogFatalWithTracef(logmessage string, args ...interface{}) {
	if overrideFunctionLogFatalWithTracef != nil {
		overrideFunctionLogFatalWithTracef(logmessage, args...)
	}

	message := fmt.Sprintf(logmessage, args...)
	LogFatalWithTrace(message)
}

var overrideFunctionLogFatalf func(logmessage string, args ...interface{})

func OverrideLogFatalf(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogFatalf = overrideFunction
}

func LogFatalf(logmessage string, args ...interface{}) {
	if overrideFunctionLogFatalf != nil {
		overrideFunctionLogFatalf(logmessage, args...)
		return
	}

	message := fmt.Sprintf(logmessage, args...)
	LogFatal(message)
}

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

func LogGoodf(logmessage string, args ...interface{}) {
	message := fmt.Sprintf(logmessage, args...)
	LogGood(message)
}

var overrideFunctionLogInfo func(logmessage string)

func OverrideLogInfo(overrideFunction func(logmessage string)) {
	overrideFunctionLogInfo = overrideFunction
}

func LogInfoByCtx(ctx context.Context, logmessage string) {
	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	LogInfo(logmessage)
}

func LogInfoByCtxf(ctx context.Context, logmessage string, args ...interface{}) {
	if !contextutils.GetVerboseFromContext(ctx) {
		return
	}

	LogInfof(logmessage, args...)
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

func LogTurnOfColorOutput() {
	globalLogSettings.SetColorEnabled(false)
}

func LogTurnOnColorOutput() {
	globalLogSettings.SetColorEnabled(true)
}

var overrideFunctionLogWarn func(logmessage string)

func OverrideLogWarn(overrideFunction func(logmessage string)) {
	overrideFunctionLogWarn = overrideFunction
}

func LogWarn(logmessage string) {
	if overrideFunctionLogWarn != nil {
		overrideFunctionLogWarn(logmessage)
		return
	}

	if globalLogSettings.IsColorEnabled() {
		logmessage = terminalcolors.GetCodeYellow() + logmessage + terminalcolors.GetCodeNoColor()
	}
	Log(logmessage)
}

var overrideFunctionLogWarnf func(logmessage string, args ...interface{})

func OverrideLogWarnf(overrideFunction func(logmessage string, args ...interface{})) {
	overrideFunctionLogWarnf = overrideFunction
}

func LogWarnf(logmessage string, args ...interface{}) {
	if overrideFunctionLogWarnf != nil {
		overrideFunctionLogWarnf(logmessage, args...)
		return
	}

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
