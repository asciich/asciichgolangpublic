package bytes

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/float"
)

func GetSizeAsHumanReadableString(sizeBytes int64) (readableSize string, err error) {
	multipliers := map[string]int64{
		"TiB": 1024 * 1024 * 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"MiB": 1024 * 1024,
		"KiB": 1024,
	}
	for _, k := range []string{"TiB", "GiB", "MiB", "KiB"} {
		v, ok := multipliers[k]
		if !ok {
			return "", fmt.Errorf("unable to get size of '%s'", k)
		}
		if sizeBytes >= v {
			const maxDigits int = 2
			readableValue, err := float.ToString(float64(sizeBytes)/float64(v), maxDigits)
			if err != nil {
				return "", fmt.Errorf("failed to format float as string: %w", err)
			}
			readableSize = readableValue + k

			return readableSize, nil
		}
	}

	return strconv.Itoa(int(sizeBytes)), nil
}

func MustGetSizeAsHumanReadableString(sizeBytes int64) (readableSize string) {
	readableSize, err := GetSizeAsHumanReadableString(sizeBytes)
	if err != nil {
		log.Panic(err)
	}

	return readableSize
}

func MustParseSizeStringAsInt64(sizeString string) (sizeBytes int64) {
	sizeBytes, err := ParseSizeStringAsInt64(sizeString)
	if err != nil {
		log.Panic(err)
	}

	return sizeBytes
}

func ParseSizeStringAsInt64(sizeString string) (sizeBytes int64, err error) {
	sizeString = strings.TrimSpace(sizeString)

	if len(sizeString) <= 0 {
		return -1, fmt.Errorf("sizeString is empty string")
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
		return -1, err
	}

	sizeBytes = int64(sizeBytesFloat * float64(multiplier))

	return sizeBytes, nil
}
