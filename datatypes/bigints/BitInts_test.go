package bigints

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigIntes_FromAndToDecimalString(t *testing.T) {
	tests := []struct {
		inputString   string
		expectedValue *big.Int
	}{
		{"-1", big.NewInt(-1)},
		{"0", big.NewInt(0)},
		{"1", big.NewInt(1)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				bigInt := MustGetFromDecimalString(tt.inputString)

				require.EqualValues(
					t,
					tt.expectedValue,
					bigInt,
				)

				require.EqualValues(
					t,
					tt.inputString,
					MustToDecimalString(bigInt),
				)
			},
		)
	}
}

func TestBigIntes_ToDecimalString_nonInitialized(t *testing.T) {
	bigInt := new(big.Int)

	require.EqualValues(
		t,
		"0",
		MustToDecimalString(bigInt),
	)
}

func TestBigIntes_IncrementDecimalString(t *testing.T) {
	for i := 0; i < 10; i++ {
		require.EqualValues(
			t,
			strconv.Itoa(i+1),
			MustIncrementDecimalString(strconv.Itoa(i)),
		)
	}
}
