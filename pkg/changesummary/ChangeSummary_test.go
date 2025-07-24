package changesummary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeSummryIsChanged(t *testing.T) {
	changeSummary := NewChangeSummary()
	require.False(t, changeSummary.IsChanged())

	for i := 0; i < 2; i++ {
		changeSummary.SetChanged(false)
		require.False(t, changeSummary.IsChanged())
	}

	for i := 0; i < 2; i++ {
		changeSummary.SetChanged(true)
		require.True(t, changeSummary.IsChanged())
	}

	for i := 0; i < 2; i++ {
		changeSummary.SetChanged(false)
		require.False(t, changeSummary.IsChanged())
	}
}

func TestChangeSummryNumberOfChnages(t *testing.T) {
	changeSummary := NewChangeSummary()
	require.False(t, changeSummary.IsChanged())
	require.EqualValues(t, 0, changeSummary.GetNumberOfChanges())

	nLoops := 3

	for i := 0; i < nLoops; i++ {
		changeSummary.IncrementNumberOfChanges()
		require.True(t, changeSummary.IsChanged())
		require.EqualValues(t, i+1, changeSummary.GetNumberOfChanges())
	}

	for i := 0; i < 2; i++ {
		changeSummary.SetChanged(true)
		require.True(t, changeSummary.IsChanged())
		require.EqualValues(t, nLoops, changeSummary.GetNumberOfChanges())
	}

	for i := 0; i < 2; i++ {
		changeSummary.SetChanged(false)
		require.False(t, changeSummary.IsChanged())
		require.EqualValues(t, 0, changeSummary.GetNumberOfChanges())
	}
}
