package windows

import (
	"github.com/asciich/asciichgolangpublic/encoding/utf16utils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/os/osutils"
)

func DecodeAsBytes(windowsUtf16 []byte) (decoded []byte, err error) {
	decoded, err = utf16utils.DecodeAsBytes(windowsUtf16)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func DecodeAsString(windowsUtf16 []byte) (decoded string, err error) {
	decoded, err = utf16utils.DecodeAsString(windowsUtf16)
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func DecodeStringAsString(windowsUtf16 string) (decoded string, err error) {
	return DecodeAsString([]byte(windowsUtf16))
}

func IsRunningOnWindows() (isRunningOnWindows bool) {
	return osutils.IsRunningOnWindows()
}

func MustDecodeAsBytes(windowsUtf16 []byte) (decoded []byte) {
	decoded, err := DecodeAsBytes(windowsUtf16)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}

func MustDecodeAsString(windowsUtf16 []byte) (decoded string) {
	decoded, err := DecodeAsString(windowsUtf16)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}

func MustDecodeStringAsString(windowsUtf16 string) (decoded string) {
	decoded, err := DecodeStringAsString(windowsUtf16)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}
