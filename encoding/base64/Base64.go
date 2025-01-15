package base64

import (
	"encoding/base64"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func DecodeStringAsBytes(input string) (decoded []byte, err error) {
	decoded, err = base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Base64 decoding failed: '%w'", err)
	}

	return decoded, nil
}

func DecodeStringAsString(input string) (decoded string, err error) {
	decodedBytes, err := DecodeStringAsBytes(input)
	if err != nil {
		return "", err
	}

	decoded = string(decodedBytes)

	return decoded, nil
}

func EncodeBytesAsString(input []byte) (encoded string, err error) {
	if input == nil {
		return "", tracederrors.TracedErrorNil("input")
	}

	encoded = base64.StdEncoding.EncodeToString(input)

	return encoded, nil
}

func EncodeStringAsBytes(input string) (encoded []byte, err error) {
	encodedString, err := EncodeStringAsString(input)
	if err != nil {
		return nil, err
	}

	encoded = []byte(encodedString)

	return encoded, nil
}

func EncodeStringAsString(input string) (encoded string, err error) {
	encoded, err = EncodeBytesAsString([]byte(input))
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func MustDecodeStringAsBytes(input string) (decoded []byte) {
	decoded, err := DecodeStringAsBytes(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}

func MustDecodeStringAsString(input string) (decoded string) {
	decoded, err := DecodeStringAsString(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return decoded
}

func MustEncodeBytesAsString(input []byte) (encoded string) {
	encoded, err := EncodeBytesAsString(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return encoded
}

func MustEncodeStringAsBytes(input string) (encoded []byte) {
	encoded, err := EncodeStringAsBytes(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return encoded
}

func MustEncodeStringAsString(input string) (encoded string) {
	encoded, err := EncodeStringAsString(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return encoded
}
