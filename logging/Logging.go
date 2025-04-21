package logging

import (
	"log"

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

func LogTurnOfColorOutput() {
	globalLogSettings.SetColorEnabled(false)
}

func LogTurnOnColorOutput() {
	globalLogSettings.SetColorEnabled(true)
}
