package slicesutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
)

func TestSlicesContainsInt(t *testing.T) {
	tests := []struct {
		intToSearch      int
		inputSlice       []int
		expectedContains bool
	}{
		{0, []int{}, false},
		{0, []int{1}, false},
		{0, []int{1, 2}, false},
		{0, []int{0, 1, 2}, true},
		{0, []int{0, 1, 2, 0}, true},
		{0, []int{1, 2, 0}, true},
		{1, []int{1, 2, 0}, true},
		{2, []int{1, 2, 0}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedContains, slicesutils.ContainsInt(tt.inputSlice, tt.intToSearch))
			},
		)
	}
}

func TestSlicesContainsString(t *testing.T) {
	tests := []struct {
		stringToSearch   string
		inputSlice       []string
		expectedContains bool
	}{
		{"0", []string{}, false},
		{"0", []string{"1"}, false},
		{"0", []string{"1", "2"}, false},
		{"0", []string{"0", "1", "2"}, true},
		{"0", []string{"0", "1", "2", "0"}, true},
		{"0", []string{"1", "2", "0"}, true},
		{"1", []string{"1", "2", "0"}, true},
		{"2", []string{"1", "2", "0"}, true},
		{"hello", []string{"hello", "Hello", "world", "World"}, true},
		{"Hello", []string{"hello", "Hello", "world", "World"}, true},
		{"HellO", []string{"hello", "Hello", "world", "World"}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedContains, slicesutils.ContainsString(tt.inputSlice, tt.stringToSearch))
			},
		)
	}
}

func TestSlicesContainsStringIgnoreCase(t *testing.T) {
	tests := []struct {
		stringToSearch   string
		inputSlice       []string
		expectedContains bool
	}{
		{"0", []string{}, false},
		{"0", []string{"1"}, false},
		{"0", []string{"1", "2"}, false},
		{"0", []string{"0", "1", "2"}, true},
		{"0", []string{"0", "1", "2", "0"}, true},
		{"0", []string{"1", "2", "0"}, true},
		{"1", []string{"1", "2", "0"}, true},
		{"2", []string{"1", "2", "0"}, true},
		{"hello", []string{"hello", "Hello", "world", "World"}, true},
		{"Hello", []string{"hello", "Hello", "world", "World"}, true},
		{"HellO", []string{"hello", "Hello", "world", "World"}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedContains, slicesutils.ContainsStringIgnoreCase(tt.inputSlice, tt.stringToSearch))
			},
		)
	}
}

func TestSlicesTrimSpace(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", "\nc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \nc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc\n"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc\n\n"}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				trimmed := slicesutils.TrimSpace(tt.input)
				require.EqualValues(t, tt.expectedOutput, trimmed)
			},
		)
	}
}

func TestSlicesRemoveMatchingStrings(t *testing.T) {
	tests := []struct {
		input          []string
		removeMatching string
		expectedOutput []string
	}{
		{[]string{}, "", []string{}},
		{[]string{}, "a", []string{}},
		{[]string{""}, "", []string{}},
		{[]string{"a"}, "a", []string{}},
		{[]string{"a", "b"}, "a", []string{"b"}},
		{[]string{"a", "b", "a"}, "a", []string{"b"}},
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "c"}, "[ab]", []string{"c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				removedMatching := slicesutils.RemoveMatchingStrings(tt.input, tt.removeMatching)
				require.EqualValues(t, tt.expectedOutput, removedMatching)
			},
		)
	}
}

func TestSlicesRemoveStringsWhichContains(t *testing.T) {
	tests := []struct {
		input          []string
		searchString   string
		expectedOutput []string
	}{
		{[]string{}, "a", []string{}},
		{[]string{"a"}, "a", []string{}},
		{[]string{"a", "b"}, "a", []string{"b"}},
		{[]string{"a", "b", "a"}, "a", []string{"b"}},
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "ca"}, "a", []string{"b"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				removedContains, err := slicesutils.RemoveStringsWhichContains(tt.input, tt.searchString)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, removedContains)
			},
		)
	}
}

func TestSlicesMaxIntValuePerIndex(t *testing.T) {
	tests := []struct {
		input1         []int
		input2         []int
		expectedOutput []int
	}{
		{nil, nil, []int{}},
		{[]int{0}, []int{1}, []int{1}},
		{[]int{0}, []int{-1}, []int{0}},
		{[]int{-10}, []int{-1}, []int{-1}},
		{[]int{-10}, []int{-1, 1, 2, 3}, []int{-1, 1, 2, 3}},
		{[]int{-10, 0, 0, 0, 4}, []int{-1, 1, 2, 3}, []int{-1, 1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				maxValues := slicesutils.MaxIntValuePerIndex(tt.input1, tt.input2)
				require.EqualValues(t, tt.expectedOutput, maxValues)
			},
		)
	}
}

func TestSlicesRemoveLastElementIfEmptyString(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{nil, []string{}},
		{[]string{""}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "b", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", "c"}, []string{"a", "b", "", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				output := slicesutils.RemoveLastElementIfEmptyString(tt.input)
				require.EqualValues(t, tt.expectedOutput, output)
			},
		)
	}
}

func TestSlicesRemoveDuplicatedEntries(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{nil, []string{}},
		{[]string{""}, []string{""}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a", ""}},
		{[]string{"a", "b", ""}, []string{"a", "b", ""}},
		{[]string{"a", "b", "", "c"}, []string{"a", "b", "", "c"}},
		{[]string{"a", "a", "", "c"}, []string{"a", "", "c"}},
		{[]string{"a", "a", "a", "c"}, []string{"a", "c"}},
		{[]string{"a", "a", "a", "a"}, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				output := slicesutils.RemoveDuplicatedStrings(tt.input)
				require.EqualValues(t, tt.expectedOutput, output)
			},
		)
	}
}

func TestSlicesStringSlicesEqual(t *testing.T) {
	tests := []struct {
		input1        []string
		input2        []string
		expectedEqual bool
	}{
		{nil, nil, false},
		{nil, []string{}, false},
		{[]string{}, nil, false},
		{[]string{}, []string{}, true},
		{[]string{}, []string{"a"}, false},
		{[]string{"a"}, []string{}, false},
		{[]string{"a"}, []string{"A"}, false},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a"}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedEqual, slicesutils.StringSlicesEqual(tt.input1, tt.input2))
			},
		)
	}
}

func TestSlicesDiffStringSlices(t *testing.T) {
	tests := []struct {
		input1          []string
		input2          []string
		expectedANotInB []string
		expectedBNotInA []string
	}{
		{nil, nil, []string{}, []string{}},
		{[]string{}, nil, []string{}, []string{}},
		{nil, []string{}, []string{}, []string{}},
		{[]string{}, []string{}, []string{}, []string{}},
		{[]string{"a"}, []string{}, []string{"a"}, []string{}},
		{[]string{"a"}, []string{"b"}, []string{"a"}, []string{"b"}},
		{[]string{""}, []string{"b"}, []string{""}, []string{"b"}},
		{[]string{""}, []string{"b", "a"}, []string{""}, []string{"a", "b"}},
		{[]string{"c"}, []string{"b", "a"}, []string{"c"}, []string{"a", "b"}},
		{[]string{"a", "c"}, []string{"b", "a"}, []string{"c"}, []string{"b"}},
		{[]string{"a", "c"}, []string{"a"}, []string{"c"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				aNotInB, bNotInA := slicesutils.DiffStringSlices(tt.input1, tt.input2)

				require.EqualValues(t, tt.expectedANotInB, aNotInB)
				require.EqualValues(t, tt.expectedBNotInA, bNotInA)
			},
		)
	}
}

func TestSlicesGetDeepCopyOfByteSlice(t *testing.T) {
	tests := []struct {
		input           []byte
		expected_output []byte
	}{
		{[]byte{}, []byte{}},
		{nil, nil},
		{[]byte("a"), []byte("a")},
		{[]byte("ab"), []byte("ab")},
		{[]byte("abc"), []byte("abc")},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				copy := slicesutils.GetDeepCopyOfByteSlice(tt.input)
				require.EqualValues(t, tt.expected_output, copy)

				for i := 0; i < len(tt.input); i++ {
					tt.input[i] = 0x00
				}

				require.EqualValues(t, tt.expected_output, copy)
			},
		)
	}
}

func TestSlicesGetDeepCopyOfStringSlice(t *testing.T) {
	tests := []struct {
		input           []string
		expected_output []string
	}{
		{nil, nil},
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a", ""}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				copy := slicesutils.GetDeepCopyOfStringsSlice(tt.input)
				require.EqualValues(t, tt.expected_output, copy)
			},
		)
	}
}

func TestSlices_GetSortedDeepCopyString(t *testing.T) {
	inputSlice := []string{"c", "b", "a"}

	sorted := slicesutils.GetSortedDeepCopyOfStringsSlice(inputSlice)
	require.EqualValues(t, []string{"a", "b", "c"}, sorted)
	require.NotEqual(t, inputSlice, sorted)
	require.EqualValues(t, []string{"c", "b", "a"}, inputSlice)
}

func TestSlices_RemoveEmptyStringsAtEnd(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{[]string{""}, []string{}},
		{[]string{"", ""}, []string{}},
		{[]string{"", "", ""}, []string{}},
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "", ""}, []string{"a"}},
		{[]string{"a", "", "", ""}, []string{"a"}},
		{[]string{"a", "b", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", "", ""}, []string{"a", "b"}},
		{[]string{"a", "", "b", ""}, []string{"a", "", "b"}},
		{[]string{"a", "", "b", "", ""}, []string{"a", "", "b"}},
		{[]string{"a", "", "b", "", "", ""}, []string{"a", "", "b"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				copy := slicesutils.RemoveEmptyStringsAtEnd(tt.input)
				require.EqualValues(t, tt.expectedOutput, copy)
			},
		)
	}
}

func Test_GetInitializedIntSlice(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		require.EqualValues(t, []int{}, slicesutils.GetInitializedIntSlice(0, 0))
	})

	t.Run("two zeros", func(t *testing.T) {
		require.EqualValues(t, []int{0, 0}, slicesutils.GetInitializedIntSlice(2, 0))
	})

	t.Run("two threes", func(t *testing.T) {
		require.EqualValues(t, []int{3, 3}, slicesutils.GetInitializedIntSlice(2, 3))
	})

	t.Run("two minus threes", func(t *testing.T) {
		require.EqualValues(t, []int{-3, -3}, slicesutils.GetInitializedIntSlice(2, -3))
	})

	t.Run("minus two minus threes", func(t *testing.T) {
		require.EqualValues(t, []int{}, slicesutils.GetInitializedIntSlice(-2, -3))
	})
}

func Test_GetInitializedIntSliceWithZeros(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		require.EqualValues(t, []int{}, slicesutils.GetInitializedIntSliceWithZeros(0))
	})

	t.Run("two", func(t *testing.T) {
		require.EqualValues(t, []int{0, 0}, slicesutils.GetInitializedIntSliceWithZeros(2))
	})

	t.Run("five", func(t *testing.T) {
		require.EqualValues(t, []int{0, 0, 0, 0, 0}, slicesutils.GetInitializedIntSliceWithZeros(5))
	})

	t.Run("minus 1", func(t *testing.T) {
		require.EqualValues(t, []int{}, slicesutils.GetInitializedIntSliceWithZeros(-5))
	})
}

func TestSlices_AddPrefixToEachString(t *testing.T) {
	tests := []struct {
		input    []string
		prefix   string
		expected []string
	}{
		{[]string{"a", "b"}, "pre", []string{"prea", "preb"}},
		{[]string{}, "pre", []string{}},
		{nil, "pre", []string{}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			result := slicesutils.AddPrefixToEachString(tt.input, tt.prefix)
			require.EqualValues(t, tt.expected, result)
		})
	}
}

func TestSlices_AddSuffixToEachString(t *testing.T) {
	tests := []struct {
		input    []string
		suffix   string
		expected []string
	}{
		{[]string{"a", "b"}, "suf", []string{"asuf", "bsuf"}},
		{[]string{}, "suf", []string{}},
		{nil, "suf", []string{}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			result := slicesutils.AddSuffixToEachString(tt.input, tt.suffix)
			require.EqualValues(t, tt.expected, result)
		})
	}
}

func TestSlices_AtLeastOneElementStartsWith(t *testing.T) {
	t.Run("at least one matches", func(t *testing.T) {
		result := slicesutils.AtLeastOneElementStartsWith([]string{"hello", "world"}, "hel")
		require.True(t, result)
	})

	t.Run("none matches", func(t *testing.T) {
		result := slicesutils.AtLeastOneElementStartsWith([]string{"hello", "world"}, "foo")
		require.False(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := slicesutils.AtLeastOneElementStartsWith([]string{}, "foo")
		require.False(t, result)
	})
}

func TestSlices_ByteSlicesEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		result := slicesutils.ByteSlicesEqual([]byte{1, 2}, []byte{1, 2})
		require.True(t, result)
	})

	t.Run("not equal", func(t *testing.T) {
		result := slicesutils.ByteSlicesEqual([]byte{1, 2}, []byte{1, 3})
		require.False(t, result)
	})

	t.Run("both nil", func(t *testing.T) {
		result := slicesutils.ByteSlicesEqual(nil, nil)
		require.False(t, result)
	})

	t.Run("one nil", func(t *testing.T) {
		result := slicesutils.ByteSlicesEqual([]byte{1}, nil)
		require.False(t, result)
	})
}

func TestSlices_ContainsAllStrings(t *testing.T) {
	t.Run("contains all", func(t *testing.T) {
		result := slicesutils.ContainsAllStrings([]string{"a", "b", "c"}, []string{"a", "b"})
		require.True(t, result)
	})

	t.Run("missing one", func(t *testing.T) {
		result := slicesutils.ContainsAllStrings([]string{"a", "b"}, []string{"a", "c"})
		require.False(t, result)
	})

	t.Run("empty toCheck", func(t *testing.T) {
		result := slicesutils.ContainsAllStrings([]string{"a", "b"}, []string{})
		require.False(t, result)
	})
}

func TestSlices_ContainsEmptyString(t *testing.T) {
	t.Run("has empty", func(t *testing.T) {
		result := slicesutils.ContainsEmptyString([]string{"a", "", "b"})
		require.True(t, result)
	})

	t.Run("no empty", func(t *testing.T) {
		result := slicesutils.ContainsEmptyString([]string{"a", "b"})
		require.False(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := slicesutils.ContainsEmptyString([]string{})
		require.False(t, result)
	})
}

func TestSlices_ContainsNoEmptyStrings(t *testing.T) {
	t.Run("no empty", func(t *testing.T) {
		result := slicesutils.ContainsNoEmptyStrings([]string{"a", "b"})
		require.True(t, result)
	})

	t.Run("has empty", func(t *testing.T) {
		result := slicesutils.ContainsNoEmptyStrings([]string{"a", "", "b"})
		require.False(t, result)
	})
}

func TestSlices_ContainsOnlyUniqeStrings(t *testing.T) {
	t.Run("all unique", func(t *testing.T) {
		result := slicesutils.ContainsOnlyUniqeStrings([]string{"a", "b", "c"})
		require.True(t, result)
	})

	t.Run("has duplicate", func(t *testing.T) {
		result := slicesutils.ContainsOnlyUniqeStrings([]string{"a", "b", "a"})
		require.False(t, result)
	})
}

func TestSlices_CountStrings(t *testing.T) {
	t.Run("count occurrences", func(t *testing.T) {
		result := slicesutils.CountStrings([]string{"a", "b", "a", "a"}, "a")
		require.Equal(t, 3, result)
	})

	t.Run("not found", func(t *testing.T) {
		result := slicesutils.CountStrings([]string{"a", "b"}, "c")
		require.Equal(t, 0, result)
	})
}

func TestSlices_GetStringElementsNotInOtherSlice(t *testing.T) {
	t.Run("some not in other", func(t *testing.T) {
		result := slicesutils.GetStringElementsNotInOtherSlice([]string{"a", "b", "c"}, []string{"b", "c"})
		require.Equal(t, []string{"a"}, result)
	})

	t.Run("all in other", func(t *testing.T) {
		result := slicesutils.GetStringElementsNotInOtherSlice([]string{"a", "b"}, []string{"a", "b", "c"})
		require.Equal(t, []string{}, result)
	})
}

func TestSlices_RemoveEmptyStrings(t *testing.T) {
	t.Run("has empty", func(t *testing.T) {
		result := slicesutils.RemoveEmptyStrings([]string{"a", "", "b", ""})
		require.Equal(t, []string{"a", "b"}, result)
	})

	t.Run("no empty", func(t *testing.T) {
		result := slicesutils.RemoveEmptyStrings([]string{"a", "b"})
		require.Equal(t, []string{"a", "b"}, result)
	})
}

func TestSlices_RemoveString(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		result := slicesutils.RemoveString([]string{"a", "b", "c"}, "b")
		require.Equal(t, []string{"a", "c"}, result)
	})

	t.Run("remove non-existing", func(t *testing.T) {
		result := slicesutils.RemoveString([]string{"a", "b"}, "c")
		require.Equal(t, []string{"a", "b"}, result)
	})
}

func TestSlices_RemoveStringEntryAtIndex(t *testing.T) {
	t.Run("valid index", func(t *testing.T) {
		result := slicesutils.RemoveStringEntryAtIndex([]string{"a", "b", "c"}, 1)
		require.Equal(t, []string{"a", "c"}, result)
	})

	t.Run("index out of bounds", func(t *testing.T) {
		result := slicesutils.RemoveStringEntryAtIndex([]string{"a", "b"}, 5)
		require.Equal(t, []string{"a", "b"}, result)
	})
}

func TestSlices_SortStringSliceAndRemoveDuplicates(t *testing.T) {
	t.Run("with duplicates", func(t *testing.T) {
		result := slicesutils.SortStringSliceAndRemoveDuplicates([]string{"c", "a", "b", "a"})
		require.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("already sorted", func(t *testing.T) {
		result := slicesutils.SortStringSliceAndRemoveDuplicates([]string{"a", "b", "c"})
		require.Equal(t, []string{"a", "b", "c"}, result)
	})
}

func TestSlices_SortStringSliceAndRemoveEmpty(t *testing.T) {
	t.Run("with empty", func(t *testing.T) {
		result := slicesutils.SortStringSliceAndRemoveEmpty([]string{"c", "", "a", "b"})
		require.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("no empty", func(t *testing.T) {
		result := slicesutils.SortStringSliceAndRemoveEmpty([]string{"c", "a", "b"})
		require.Equal(t, []string{"a", "b", "c"}, result)
	})
}

func TestSlices_SplitStrings(t *testing.T) {
	t.Run("split", func(t *testing.T) {
		result := slicesutils.SplitStrings([]string{"a-b", "c-d"}, "-")
		require.Equal(t, []string{"a", "b", "c", "d"}, result)
	})

	t.Run("no separator", func(t *testing.T) {
		result := slicesutils.SplitStrings([]string{"ab", "cd"}, "-")
		require.Equal(t, []string{"ab", "cd"}, result)
	})
}

func TestSlices_SplitStringsAndRemoveEmpty(t *testing.T) {
	t.Run("split with empty", func(t *testing.T) {
		result := slicesutils.SplitStringsAndRemoveEmpty([]string{"a-", "c-d"}, "-")
		require.Equal(t, []string{"a", "c", "d"}, result)
	})
}

func TestSlices_ToLower(t *testing.T) {
	t.Run("mixed case", func(t *testing.T) {
		result := slicesutils.ToLower([]string{"Hello", "WORLD"})
		require.Equal(t, []string{"hello", "world"}, result)
	})
}

func TestSlices_TrimAllPrefix(t *testing.T) {
	t.Run("trim prefix", func(t *testing.T) {
		result := slicesutils.TrimAllPrefix([]string{"prea", "preb"}, "pre")
		require.Equal(t, []string{"a", "b"}, result)
	})
}

func TestSlices_TrimPrefix(t *testing.T) {
	t.Run("trim prefix", func(t *testing.T) {
		result := slicesutils.TrimPrefix([]string{"prea", "preb"}, "pre")
		require.Equal(t, []string{"a", "b"}, result)
	})
}
