package asciichgolangpublic

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type FloatService struct {
}

func Float() (floatService *FloatService) {
	return NewFloatService()
}

func NewFloatService() (floatService *FloatService) {
	return new(FloatService)
}

// Format the given float64 'input' as string with 'maxDigits' digits after the '.'.
// Tailing zeros are removed.
func (f *FloatService) ToString(input float64, maxDigits int) (formatedFloat string, err error) {
	if maxDigits < 0 {
		return "", TracedErrorf("Negative maxDigits='%d' not allowed", maxDigits)
	}

	roundedInput, err := f.Round(input, maxDigits)
	if err != nil {
		return "", err
	}

	formatedFloat = fmt.Sprintf("%."+strconv.Itoa(maxDigits)+"f", roundedInput)
	if strings.Contains(formatedFloat, ".") {
		formatedFloat = Strings().TrimAllSuffix(formatedFloat, "0")
	}
	formatedFloat = strings.TrimSuffix(formatedFloat, ".")

	return formatedFloat, nil
}

func (f *FloatService) MustRound(input float64, nDigits int) (rounded float64) {
	rounded, err := f.Round(input, nDigits)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rounded
}

func (f *FloatService) MustToString(input float64, maxDigits int) (formatedFloat string) {
	formatedFloat, err := f.ToString(input, maxDigits)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return formatedFloat
}

func (f *FloatService) Round(input float64, nDigits int) (rounded float64, err error) {
	if nDigits < 0 {
		return -1.0, TracedErrorf("Negative nDigits='%d' not allowed", nDigits)
	}

	multiplier := math.Pow(10.0, float64(nDigits))
	rounded = math.Round(input*multiplier) / multiplier

	return rounded, nil
}
