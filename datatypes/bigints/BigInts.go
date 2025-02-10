package bigints

import (
	"fmt"
	"log"
	"math/big"
)

func MustGetFromDecimalString(decimal string) (bigInt *big.Int) {
	bigInt, err := GetFromDecimalString(decimal)
	if err != nil {
		log.Panic(err)
	}

	return bigInt
}

func GetFromDecimalString(decimal string) (bigInt *big.Int, err error) {
	if decimal == "" {
		return nil, fmt.Errorf("decimal is empty string")
	}

	bigInt = new(big.Int)

	_, ok := bigInt.SetString(decimal, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse decimal string '%s' as *big.Int", bigInt)
	}

	return bigInt, nil
}

func MustToDecimalString(bigInt *big.Int) (decimal string) {
	decimal, err := ToDecimalString(bigInt)
	if err != nil {
		log.Panic(err)
	}

	return decimal
}

func ToDecimalString(bigInt *big.Int) (decimal string, err error) {
	if bigInt == nil {
		return "", fmt.Errorf("bigInt is nil pointer")
	}

	decimal = bigInt.String()

	return decimal, nil
}

func MustIncrementDecimalString(decimal string) (incremented string) {
	incremented, err := IncrementDecimalString(decimal)
	if err != nil {
		log.Panic(err)
	}

	return incremented
}

func AddIntToDecimalString(decimal string, toAdd int) (result string, err error) {
	if decimal == "" {
		return "", fmt.Errorf("decimal is empty string")
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
		return "", fmt.Errorf("decimal is empty string")
	}

	incremented, err = AddIntToDecimalString(decimal, 1)
	if err != nil {
		return "", err
	}

	return incremented, nil
}
