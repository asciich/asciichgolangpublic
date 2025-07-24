package bigintutils_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/bigintutils"
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
				bigInt, err := bigintutils.GetFromDecimalString(tt.inputString)
				require.NoError(t, err)

				require.EqualValues(
					t,
					tt.expectedValue,
					bigInt,
				)

				decimalString, err := bigintutils.ToDecimalString(bigInt)
				require.NoError(t, err)

				require.EqualValues(
					t,
					tt.inputString,
					decimalString,
				)
			},
		)
	}
}

func TestBigIntes_ToDecimalString_nonInitialized(t *testing.T) {
	bigInt := new(big.Int)

	decimalString, err := bigintutils.ToDecimalString(bigInt)
	require.NoError(t, err)

	require.EqualValues(t, "0", decimalString)
}

func TestBigIntes_IncrementDecimalString(t *testing.T) {
	for i := 0; i < 10; i++ {
		incremented, err := bigintutils.IncrementDecimalString(strconv.Itoa(i))
		require.NoError(t, err)
		require.EqualValues(
			t,
			strconv.Itoa(i+1),
			incremented,
		)
	}
}

func TestGetAsHexColonSeparatedString(t *testing.T) {
	tests := []struct {
		input    *big.Int
		expected string
	}{
		{big.NewInt(0), "00"},
		{big.NewInt(1), "01"},
		{big.NewInt(10), "0A"},
		{big.NewInt(15), "0F"},
		{big.NewInt(255), "FF"},
		{big.NewInt(256), "01:00"},
		{big.NewInt(256 + 255), "01:FF"},
		{big.NewInt(1024), "04:00"},
		{big.NewInt(256 * 256), "01:00:00"},
		{big.NewInt(256 * 256 * 2), "02:00:00"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				out, err := bigintutils.ToHexStringColonSeparated(tt.input)
				require.NoError(t, err)

				require.EqualValues(
					t,
					tt.expected,
					out,
				)
			},
		)
	}
}

func Test_EqualInts(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		require.False(t, bigintutils.EqualsInts(nil, nil))
	})

	t.Run("i1 nil", func(t *testing.T) {
		require.False(t, bigintutils.EqualsInts(nil, big.NewInt(1)))
	})

	t.Run("i2 nil", func(t *testing.T) {
		require.False(t, bigintutils.EqualsInts(big.NewInt(1), nil))
	})

	t.Run("equals", func(t *testing.T) {
		require.True(t, bigintutils.EqualsInts(big.NewInt(1), big.NewInt(1)))
	})

	t.Run("not equal", func(t *testing.T) {
		require.False(t, bigintutils.EqualsInts(big.NewInt(1), big.NewInt(3)))
	})
}

func Test_GreaterThanInts(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		require.False(t, bigintutils.GreatherThanInts(nil, nil))
	})

	t.Run("i1 nil", func(t *testing.T) {
		require.False(t, bigintutils.GreatherThanInts(nil, big.NewInt(1)))
	})

	t.Run("i2 nil", func(t *testing.T) {
		require.False(t, bigintutils.GreatherThanInts(big.NewInt(1), nil))
	})

	t.Run("equals", func(t *testing.T) {
		require.False(t, bigintutils.GreatherThanInts(big.NewInt(1), big.NewInt(1)))
	})

	t.Run("lower", func(t *testing.T) {
		require.False(t, bigintutils.GreatherThanInts(big.NewInt(1), big.NewInt(3)))
	})

	t.Run("greater", func(t *testing.T) {
		require.True(t, bigintutils.GreatherThanInts(big.NewInt(2), big.NewInt(1)))
	})
}

func Test_Equals(t *testing.T) {
	t.Run("equals", func(t *testing.T) {
		equals, err := bigintutils.Equals("1234", "1234")
		require.NoError(t, err)
		require.True(t, equals)
	})

	t.Run("not equal", func(t *testing.T) {
		equals, err := bigintutils.Equals("12345", "1234")
		require.NoError(t, err)
		require.False(t, equals)
	})
}

func Test_GetRandomBigIntByInts(t *testing.T) {
	t.Run("All args nil", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(nil, nil)
		require.Error(t, err)
		require.Nil(t, bigInt)
	})

	t.Run("min nil", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(1), nil)
		require.Error(t, err)
		require.Nil(t, bigInt)
	})

	t.Run("max nil", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(nil, big.NewInt(1))
		require.Error(t, err)
		require.Nil(t, bigInt)
	})

	t.Run("min max equal", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(1), big.NewInt(1))
		require.Error(t, err)
		require.Nil(t, bigInt)
	})

	t.Run("max greater than min", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(2), big.NewInt(1))
		require.Error(t, err)
		require.Nil(t, bigInt)
	})

	t.Run("Only one possible value 1", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(1), big.NewInt(2))
		require.NoError(t, err)
		require.True(t, bigintutils.EqualsInts(big.NewInt(1), bigInt))
	})

	t.Run("Only one possible value 12345", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(1234), big.NewInt(1235))
		require.NoError(t, err)
		require.True(t, bigintutils.EqualsInts(big.NewInt(1234), bigInt))
	})

	t.Run("Only one possible value -3", func(t *testing.T) {
		bigInt, err := bigintutils.GetRandomBigIntByInts(big.NewInt(-3), big.NewInt(-2))
		require.NoError(t, err)
		require.True(t, bigintutils.EqualsInts(big.NewInt(-3), bigInt))
	})
}
