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

func (b *Base64Service) MustDecodeStringAsBytes(input string) (decoded []byte) {
	decoded, err := b.DecodeStringAsBytes(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return decoded
}
