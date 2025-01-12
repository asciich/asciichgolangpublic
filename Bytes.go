package asciichgolangpublic

import (
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/float"
)

type BytesService struct{}

func Bytes() (bytesService *BytesService) {
	return NewBytesService()
}

func NewBytesService() (bytesService *BytesService) {
	return new(BytesService)
}

func (b *BytesService) GetSizeAsHumanReadableString(sizeBytes int64) (readableSize string, err error) {
	multipliers := map[string]int64{
		"TiB": 1024 * 1024 * 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"MiB": 1024 * 1024,
		"KiB": 1024,
	}
	for _, k := range []string{"TiB", "GiB", "MiB", "KiB"} {
		v, ok := multipliers[k]
		if !ok {
			return "", TracedErrorf("Unable to get size of '%s'", k)
		}
		if sizeBytes >= v {
			const maxDigits int = 2
			readableValue, err := float.ToString(float64(sizeBytes)/float64(v), maxDigits)
			if err != nil {
				return "", err
			}
			readableSize = readableValue + k

			return readableSize, nil
		}
	}

	return strconv.Itoa(int(sizeBytes)), nil
}

func (b *BytesService) MustGetSizeAsHumanReadableString(sizeBytes int64) (readableSize string) {
	readableSize, err := b.GetSizeAsHumanReadableString(sizeBytes)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return readableSize
}

func (b *BytesService) MustParseSizeStringAsInt64(sizeString string) (sizeBytes int64) {
	sizeBytes, err := b.ParseSizeStringAsInt64(sizeString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sizeBytes
}

func (b *BytesService) ParseSizeStringAsInt64(sizeString string) (sizeBytes int64, err error) {
	sizeString = strings.TrimSpace(sizeString)

	if len(sizeString) <= 0 {
		return -1, TracedError("sizeString is empty string")
	}

	var multiplier int64 = 1
	multipliers := map[string]int64{
		"kB":  1000,
		"KB":  1024,
		"KiB": 1024,
		"MB":  1024 * 1024,
		"MiB": 1024 * 1024,
		"GB":  1024 * 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TB":  1024 * 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
	}

	for k, v := range multipliers {
		if strings.Contains(sizeString, k) {
			sizeString = strings.ReplaceAll(sizeString, k, "")
			multiplier = v
			break
		}
	}

	sizeBytesFloat, err := strconv.ParseFloat(sizeString, 64)
	if err != nil {
		return -1, TracedError(err.Error())
	}

	sizeBytes = int64(sizeBytesFloat * float64(multiplier))

	return sizeBytes, nil
}
