package ansiblegalaxyutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
)

func Test_IsValidCollectionName(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName(""))
	})

	t.Run("single char", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("a"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("b"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("c"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("z"))
	})

	t.Run("double char", func(t *testing.T) {
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("ab"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("bc"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("cc"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("zc"))
	})

	t.Run("leading underscore is invalid", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("_b"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("_ab"))
	})

	t.Run("tailing underscore is invalid", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("b_"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("ab_"))
	})	
	
	t.Run("valid without underscore", func(t *testing.T) {
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("abcdef"))
	})


	t.Run("valid with underscore", func(t *testing.T) {
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("a_bcdef"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("ab_cdef"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("abc_def"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("abcd_ef"))
		require.True(t, ansiblegalaxyutils.IsValidCollectionName("abcde_f"))
	})

	t.Run("numbers are invalid", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("a1bcdef"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("ab0cdef"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abc2def"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abcd3ef"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abcde4f"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("a1bcdef5"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("ab0cdef6"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abc2def7"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abcd3ef8"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abcde4f9"))
	})

	t.Run("Uppercase is invalid", func(t *testing.T) {
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("Aa"))
		require.False(t, ansiblegalaxyutils.IsValidCollectionName("abcDefg"))
	})
}


func Test_CheckValidCollectionName(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName(""))
	})

	t.Run("single char", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("a"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("b"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("c"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("z"))
	})

	t.Run("double char", func(t *testing.T) {
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("ab"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("bc"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("cc"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("zc"))
	})

	t.Run("leading underscore is invalid", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("_b"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("_ab"))
	})

	t.Run("tailing underscore is invalid", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("b_"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("ab_"))
	})	
	
	t.Run("valid without underscore", func(t *testing.T) {
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("abcdef"))
	})


	t.Run("valid with underscore", func(t *testing.T) {
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("a_bcdef"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("ab_cdef"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("abc_def"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("abcd_ef"))
		require.NoError(t, ansiblegalaxyutils.CheckValidCollectionName("abcde_f"))
	})

	t.Run("numbers are invalid", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("a1bcdef"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("ab0cdef"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abc2def"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abcd3ef"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abcde4f"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("a1bcdef5"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("ab0cdef6"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abc2def7"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abcd3ef8"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abcde4f9"))
	})

	t.Run("Uppercase is invalid", func(t *testing.T) {
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("Aa"))
		require.Error(t, ansiblegalaxyutils.CheckValidCollectionName("abcDefg"))
	})
}
