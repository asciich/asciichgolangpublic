package asciichgolangpublic

import "runtime"

type OsService struct{}

func NewOsService() (o *OsService) {
	return new(OsService)
}

func OS() (o *OsService) {
	return NewOsService()
}

func (o *OsService) IsRunningOnWindows() (isRunningOnWindows bool) {
	return runtime.GOOS == "windows"
}
