package asciichgolangpublic

import "runtime"

type OsService struct{}

// Provides Windows (the operating system) related functions.
func OS() (o *OsService) {
	return NewOsService()
}

func NewOsService() (o *OsService) {
	return new(OsService)
}

func (o *OsService) IsRunningOnWindows() (isRunningOnWindows bool) {
	return runtime.GOOS == "windows"
}
