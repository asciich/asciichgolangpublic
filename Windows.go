package asciichgolangpublic

type WindowsService struct{}

// Provides Windows (the operating system) related functions.
func Windows() (w *WindowsService) {
	return NewWindowsService()
}

func NewWindowsService() (w *WindowsService) {
	return new(WindowsService)
}

func (w *WindowsService) DecodeAsBytes(windowsUtf16 []byte) (decoded []byte, err error) {
	decoded, err = UTF16().DecodeAsBytes(windowsUtf16)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (w *WindowsService) DecodeAsString(windowsUtf16 []byte) (decoded string, err error) {
	decoded, err = UTF16().DecodeAsString(windowsUtf16)
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func (w *WindowsService) DecodeStringAsString(windowsUtf16 string) (decoded string, err error) {
	return w.DecodeAsString([]byte(windowsUtf16))
}

func (w *WindowsService) IsRunningOnWindows() (isRunningOnWindows bool) {
	return OS().IsRunningOnWindows()
}

func (w *WindowsService) MustDecodeAsBytes(windowsUtf16 []byte) (decoded []byte) {
	decoded, err := w.DecodeAsBytes(windowsUtf16)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}

func (w *WindowsService) MustDecodeAsString(windowsUtf16 []byte) (decoded string) {
	decoded, err := w.DecodeAsString(windowsUtf16)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}

func (w *WindowsService) MustDecodeStringAsString(windowsUtf16 string) (decoded string) {
	decoded, err := w.DecodeStringAsString(windowsUtf16)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}
