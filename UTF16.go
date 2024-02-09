package asciichgolangpublic

import (
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type UTF16Service struct{}

func NewUTF16Service() (u *UTF16Service) {
	return new(UTF16Service)
}

func UTF16() (u *UTF16Service) {
	return NewUTF16Service()
}

func (u *UTF16Service) DecodeAsBytes(utf16 []byte) (decoded []byte, err error) {
	if len(utf16) < 2 {
		return utf16, err
	}

	if len(utf16) > 2 {
		if utf16[1] != 0x00 {
			// no decode needed.
			return utf16, nil
		}
	}

	decoded, _, err = transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), utf16)
	if err != nil {
		return nil, TracedError(err)
	}

	return decoded, nil
}

func (u *UTF16Service) DecodeAsString(utf16 []byte) (decoded string, err error) {
	decodedBytes, err := u.DecodeAsBytes(utf16)
	if err != nil {
		return "", TracedError(err)
	}

	decoded = string(decodedBytes)

	return decoded, nil
}

func (u *UTF16Service) MustDecodeAsBytes(utf16 []byte) (decoded []byte) {
	decoded, err := u.DecodeAsBytes(utf16)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}

func (u *UTF16Service) MustDecodeAsString(utf16 []byte) (decoded string) {
	decoded, err := u.DecodeAsString(utf16)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}
