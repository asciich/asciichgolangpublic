package asciichgolangpublic

type ExitCodesService struct {
}

func ExitCodes() (exitCodes *ExitCodesService) {
	return new(ExitCodesService)
}

func NewExitCodesService() (e *ExitCodesService) {
	return new(ExitCodesService)
}

func (e *ExitCodesService) ExitCodeOK() (exitCode int) {
	return 0
}

func (e *ExitCodesService) ExitCodeTimeout() (exitCode int) {
	return 124
}
