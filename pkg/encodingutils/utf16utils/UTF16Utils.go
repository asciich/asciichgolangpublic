package utf16utils

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func DecodeAsBytes(utf16 []byte) (decoded []byte, err error) {
	if len(utf16) < 2 {
		return utf16, nil
	}

	if len(utf16) > 2 {
		if utf16[1] != 0x00 {
			// no decode needed.
			return utf16, nil
		}
	}

	decoded, _, err = transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), utf16)
	if err != nil {
		return nil, tracederrors.TracedError(err)
	}

	return decoded, nil
}

func DecodeAsString(utf16 []byte) (decoded string, err error) {
	decodedBytes, err := DecodeAsBytes(utf16)
	if err != nil {
		return "", tracederrors.TracedError(err)
	}

	decoded = string(decodedBytes)

	return decoded, nil
}

func MustDecodeAsBytes(utf16 []byte) (decoded []byte) {
	decoded, err := DecodeAsBytes(utf16)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}

func MustDecodeAsString(utf16 []byte) (decoded string) {
	decoded, err := DecodeAsString(utf16)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}
