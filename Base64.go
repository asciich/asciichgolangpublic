package asciichgolangpublic

import "encoding/base64"

type Base64Service struct{}

func Base64() (b *Base64Service) {
	return NewBase64Service()
}

func NewBase64Service() (b *Base64Service) {
	return new(Base64Service)
}

func (b *Base64Service) DecodeStringAsBytes(input string) (decoded []byte, err error) {
	decoded, err = base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, TracedErrorf("Base64 decoding failed: '%w'", err)
	}

	return decoded, nil
}

func (b *Base64Service) DecodeStringAsString(input string) (decoded string, err error) {
	decodedBytes, err := b.DecodeStringAsBytes(input)
	if err != nil {
		return "", err
	}

	decoded = string(decodedBytes)

	return decoded, nil
}

func (b *Base64Service) EncodeBytesAsString(input []byte) (encoded string, err error) {
	if input == nil {
		return "", TracedErrorNil("input")
	}

	encoded = base64.StdEncoding.EncodeToString(input)

	return encoded, nil
}

func (b *Base64Service) EncodeStringAsBytes(input string) (encoded []byte, err error) {
	encodedString, err := b.EncodeStringAsString(input)
	if err != nil {
		return nil, err
	}

	encoded = []byte(encodedString)

	return encoded, nil
}

func (b *Base64Service) EncodeStringAsString(input string) (encoded string, err error) {
	encoded, err = b.EncodeBytesAsString([]byte(input))
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func (b *Base64Service) MustDecodeStringAsBytes(input string) (decoded []byte) {
	decoded, err := b.DecodeStringAsBytes(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}

func (b *Base64Service) MustDecodeStringAsString(input string) (decoded string) {
	decoded, err := b.DecodeStringAsString(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}

func (b *Base64Service) MustEncodeBytesAsString(input []byte) (encoded string) {
	encoded, err := b.EncodeBytesAsString(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return encoded
}

func (b *Base64Service) MustEncodeStringAsBytes(input string) (encoded []byte) {
	encoded, err := b.EncodeStringAsBytes(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return encoded
}

func (b *Base64Service) MustEncodeStringAsString(input string) (encoded string) {
	encoded, err := b.EncodeStringAsString(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return encoded
}
