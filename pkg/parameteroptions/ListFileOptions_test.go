package parameteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestListFileOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := parameteroptions.NewListFileOptions()
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		copy.ReturnRelativePaths = true
		copy.OnlyFiles = true
		copy.NonRecursive = true
		copy.AllowEmptyListIfNoFileIsFound = true

		// Original should be unchanged
		require.EqualValues(t, false, original.ReturnRelativePaths)
		require.EqualValues(t, false, original.OnlyFiles)
		require.EqualValues(t, false, original.NonRecursive)
		require.EqualValues(t, false, original.AllowEmptyListIfNoFileIsFound)
	})

	t.Run("with MatchBasenamePattern slice", func(t *testing.T) {
		original := parameteroptions.NewListFileOptions()
		err := original.SetMatchBasenamePattern([]string{"*.txt", "*.md"})
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		origPattern, err := original.GetMatchBasenamePattern()
		require.NoError(t, err)
		copyPattern, err := copy.GetMatchBasenamePattern()
		require.NoError(t, err)

		require.EqualValues(t, origPattern, copyPattern)

		// Modify copy's slice
		copy.MatchBasenamePattern[0] = "modified"
		copy.MatchBasenamePattern = append(copy.MatchBasenamePattern, "*.go")

		// Original should be unchanged (deep copy of slice)
		origPattern, err = original.GetMatchBasenamePattern()
		require.NoError(t, err)
		require.EqualValues(t, []string{"*.txt", "*.md"}, origPattern)

		copyPattern, err = copy.GetMatchBasenamePattern()
		require.NoError(t, err)
		require.EqualValues(t, []string{"modified", "*.md", "*.go"}, copyPattern)
	})

	t.Run("with ExcludeBasenamePattern slice", func(t *testing.T) {
		original := parameteroptions.NewListFileOptions()
		err := original.SetExcludeBasenamePattern([]string{"*.tmp", "*.bak"})
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		origPattern, err := original.GetExcludeBasenamePattern()
		require.NoError(t, err)
		copyPattern, err := copy.GetExcludeBasenamePattern()
		require.NoError(t, err)

		require.EqualValues(t, origPattern, copyPattern)

		// Modify copy's slice
		copy.ExcludeBasenamePattern[0] = "modified"

		// Original should be unchanged
		origPattern, err = original.GetExcludeBasenamePattern()
		require.NoError(t, err)
		require.EqualValues(t, []string{"*.tmp", "*.bak"}, origPattern)
	})

	t.Run("with ExcludePatternWholepath slice", func(t *testing.T) {
		original := parameteroptions.NewListFileOptions()
		err := original.SetExcludePatternWholepath([]string{"/tmp/*", "/var/log/*"})
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		origPattern, err := original.GetExcludePatternWholepath()
		require.NoError(t, err)
		copyPattern, err := copy.GetExcludePatternWholepath()
		require.NoError(t, err)

		require.EqualValues(t, origPattern, copyPattern)

		// Modify copy's slice
		copy.ExcludePatternWholepath[0] = "modified"
		copy.ExcludePatternWholepath = append(copy.ExcludePatternWholepath, "/home/*")

		// Original should be unchanged
		origPattern, err = original.GetExcludePatternWholepath()
		require.NoError(t, err)
		require.EqualValues(t, []string{"/tmp/*", "/var/log/*"}, origPattern)
	})

	t.Run("with all slices set", func(t *testing.T) {
		original := parameteroptions.NewListFileOptions()
		err := original.SetMatchBasenamePattern([]string{"*.txt"})
		require.NoError(t, err)
		err = original.SetExcludeBasenamePattern([]string{"*.tmp"})
		require.NoError(t, err)
		err = original.SetExcludePatternWholepath([]string{"/tmp/*"})
		require.NoError(t, err)
		original.ReturnRelativePaths = true
		original.OnlyFiles = true
		original.NonRecursive = true
		original.AllowEmptyListIfNoFileIsFound = true

		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.EqualValues(t, original.MatchBasenamePattern, copy.MatchBasenamePattern)
		require.EqualValues(t, original.ExcludeBasenamePattern, copy.ExcludeBasenamePattern)
		require.EqualValues(t, original.ExcludePatternWholepath, copy.ExcludePatternWholepath)
		require.EqualValues(t, original.ReturnRelativePaths, copy.ReturnRelativePaths)
		require.EqualValues(t, original.OnlyFiles, copy.OnlyFiles)
		require.EqualValues(t, original.NonRecursive, copy.NonRecursive)
		require.EqualValues(t, original.AllowEmptyListIfNoFileIsFound, copy.AllowEmptyListIfNoFileIsFound)

		// Modify copy's slices
		copy.MatchBasenamePattern[0] = "modified"
		copy.ExcludeBasenamePattern[0] = "modified"
		copy.ExcludePatternWholepath[0] = "modified"

		// Original should be unchanged
		require.EqualValues(t, []string{"*.txt"}, original.MatchBasenamePattern)
		require.EqualValues(t, []string{"*.tmp"}, original.ExcludeBasenamePattern)
		require.EqualValues(t, []string{"/tmp/*"}, original.ExcludePatternWholepath)
	})
}
