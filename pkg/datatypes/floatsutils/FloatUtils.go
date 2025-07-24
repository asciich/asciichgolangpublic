package floatsutils

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// Format the given float64 'input' as string with 'maxDigits' digits after the '.'.
// Tailing zeros after '.' are removed.
func ToString(input float64, maxDigits int) (formatedFloat string, err error) {
	if maxDigits < 0 {
		return "", fmt.Errorf("negative maxDigits='%d' not allowed", maxDigits)
	}

	roundedInput, err := Round(input, maxDigits)
	if err != nil {
		return "", err
	}

	formatedFloat = fmt.Sprintf("%."+strconv.Itoa(maxDigits)+"f", roundedInput)
	if strings.Contains(formatedFloat, ".") {
		for {
			if strings.HasSuffix(formatedFloat, "0") {
				formatedFloat = strings.TrimSuffix(formatedFloat, "0")
				continue
			}

			break
		}
	}
	formatedFloat = strings.TrimSuffix(formatedFloat, ".")

	return formatedFloat, nil
}

func MustRound(input float64, nDigits int) (rounded float64) {
	rounded, err := Round(input, nDigits)
	if err != nil {
		log.Panic(err)
	}

	return rounded
}

func MustToString(input float64, maxDigits int) (formatedFloat string) {
	formatedFloat, err := ToString(input, maxDigits)
	if err != nil {
		log.Panic(err)
	}

	return formatedFloat
}

func Round(input float64, nDigits int) (rounded float64, err error) {
	if nDigits < 0 {
		return -1.0, fmt.Errorf("negative nDigits='%d' not allowed", nDigits)
	}

	multiplier := math.Pow(10.0, float64(nDigits))
	rounded = math.Round(input*multiplier) / multiplier

	return rounded, nil
}
