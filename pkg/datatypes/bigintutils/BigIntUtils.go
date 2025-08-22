package bigintutils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Equals(i1 string, i2 string) (bool, error) {
	i1int, err := GetFromDecimalString(i1)
	if err != nil {
		return false, err
	}

	i2int, err := GetFromDecimalString(i2)
	if err != nil {
		return false, err
	}

	return EqualsInts(i1int, i2int), nil
}

func EqualsInts(i1 *big.Int, i2 *big.Int) bool {
	if i1 == nil {
		return false
	}

	if i2 == nil {
		return false
	}

	return i1.Cmp(i2) == 0
}

func GetFromDecimalString(decimal string) (bigInt *big.Int, err error) {
	if decimal == "" {
		return nil, tracederrors.TracedErrorEmptyString("decimal")
	}

	bigInt = new(big.Int)

	_, ok := bigInt.SetString(decimal, 10)
	if !ok {
		return nil, tracederrors.TracedErrorf("failed to parse decimal string '%s' as *big.Int", bigInt)
	}

	return bigInt, nil
}

func GreatherThanInts(i1 *big.Int, i2 *big.Int) bool {
	if i1 == nil {
		return false
	}

	if i2 == nil {
		return false
	}

	return i1.Cmp(i2) > 0
}

// Int returns a uniform random value in [min, max) so min is included and max not.
func GetRandomBigIntByInts(min *big.Int, max *big.Int) (random *big.Int, err error) {
	if min == nil {
		return nil, tracederrors.TracedErrorNil("min")
	}

	if max == nil {
		return nil, tracederrors.TracedErrorNil("max")
	}

	if !GreatherThanInts(max, min) {
		minStr, err := ToDecimalString(min)
		if err != nil {
			return nil, err
		}

		maxStr, err := ToDecimalString(max)
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf("unable to generate random big.Int: min '%s' is not lower than max '%s'", minStr, maxStr)
	}

	randomRange := max.Sub(max, min)

	generated, err := rand.Int(rand.Reader, randomRange)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to generate serial number: %w", err)
	}
	random = generated.Add(generated, min)

	return random, nil
}

func ToDecimalString(bigInt *big.Int) (decimal string, err error) {
	if bigInt == nil {
		return "", tracederrors.TracedErrorNil("bigInt")
	}

	decimal = bigInt.String()

	return decimal, nil
}

func AddIntToDecimalString(decimal string, toAdd int) (result string, err error) {
	if decimal == "" {
		return "", tracederrors.TracedErrorEmptyString("decimal")
	}

	bigInt, err := GetFromDecimalString(decimal)
	if err != nil {
		return "", err
	}

	bigToAdd := big.NewInt(int64(toAdd))

	bigResult := bigInt.Add(bigInt, bigToAdd)

	return ToDecimalString(bigResult)
}

func IncrementDecimalString(decimal string) (incremented string, err error) {
	if decimal == "" {
		return "", tracederrors.TracedErrorEmptyString("decimal")
	}

	incremented, err = AddIntToDecimalString(decimal, 1)
	if err != nil {
		return "", err
	}

	return incremented, nil
}

func ToHexStringColonSeparated(input *big.Int) (out string, err error) {
	if input == nil {
		return "", tracederrors.TracedErrorNil("input")
	}

	serialBytes := input.Bytes()
	if len(serialBytes) <= 0 {
		return "00", nil
	}

	hexStrings := make([]string, len(serialBytes))
	for i, b := range serialBytes {
		hexStrings[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(hexStrings, ":"), nil
}
