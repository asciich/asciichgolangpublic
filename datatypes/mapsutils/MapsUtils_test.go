package mapsutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/datatypes/mapsutils"
)

func Test_GetKeysOfStringMap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		keys := mapsutils.GetKeysOfStringMap(nil)
		require.EqualValues(t, []string{}, keys)
	})

	t.Run("one entry", func(t *testing.T) {
		keys := mapsutils.GetKeysOfStringMap(map[string]string{
			"a": "abc",
		})
		require.EqualValues(t, []string{"a"}, keys)
	})

	t.Run("two entry sorted", func(t *testing.T) {
		keys := mapsutils.GetKeysOfStringMap(map[string]string{
			"a": "abc",
			"b": "def",
		})
		require.EqualValues(t, []string{"a", "b"}, keys)
	})

	t.Run("two entry unsorted", func(t *testing.T) {
		keys := mapsutils.GetKeysOfStringMap(map[string]string{
			"b": "abc",
			"a": "def",
		})
		require.EqualValues(t, []string{"a", "b"}, keys)
	})
}

func Test_DeepCopyBytesMap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		copy := mapsutils.DeepCopyBytesMap(nil)
		require.Nil(t, copy)
	})

	t.Run("one entry", func(t *testing.T) {
		originalMap := map[string][]byte{
			"abc": []byte("entry"),
		}

		copy := mapsutils.DeepCopyBytesMap(originalMap)

		require.Len(t, copy, 1)
		require.EqualValues(t, originalMap, copy)

		copy["abc"] = []byte{}
		require.NotEqualValues(t, originalMap, copy)
		require.EqualValues(t, originalMap["abc"], []byte("entry"))
	})
}
