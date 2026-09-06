package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabCreateMergeRequestOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabCreateMergeRequestOptions{}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.SourceBranchName, copy.SourceBranchName)
		require.EqualValues(t, original.TargetBranchName, copy.TargetBranchName)
		require.EqualValues(t, original.Title, copy.Title)
		require.EqualValues(t, original.Description, copy.Description)
		require.EqualValues(t, original.Labels, copy.Labels)
		require.EqualValues(t, original.SquashEnabled, copy.SquashEnabled)
		require.EqualValues(t, original.DeleteSourceBranchOnMerge, copy.DeleteSourceBranchOnMerge)
		require.EqualValues(t, original.FailIfMergeRequestAlreadyExists, copy.FailIfMergeRequestAlreadyExists)
		require.EqualValues(t, original.AssignToSelf, copy.AssignToSelf)

		// Modify copy
		copy.SourceBranchName = "modified"
		copy.TargetBranchName = "modified"
		copy.Title = "modified"
		copy.Description = "modified"
		copy.Labels = []string{"modified"}
		copy.SquashEnabled = !original.SquashEnabled
		copy.DeleteSourceBranchOnMerge = !original.DeleteSourceBranchOnMerge
		copy.FailIfMergeRequestAlreadyExists = !original.FailIfMergeRequestAlreadyExists
		copy.AssignToSelf = !original.AssignToSelf

		// Original should be unchanged
		require.EqualValues(t, "", original.SourceBranchName)
		require.EqualValues(t, "", original.TargetBranchName)
		require.EqualValues(t, "", original.Title)
		require.EqualValues(t, "", original.Description)
		require.Nil(t, original.Labels)
	})

	t.Run("with Labels slice", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabCreateMergeRequestOptions{
			SourceBranchName: "feature",
			TargetBranchName: "main",
			Title:            "Add feature",
			Description:      "New feature implementation",
			Labels:           []string{"feature", "enhancement"},
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.SourceBranchName, copy.SourceBranchName)
		require.EqualValues(t, original.TargetBranchName, copy.TargetBranchName)
		require.EqualValues(t, original.Title, copy.Title)
		require.EqualValues(t, original.Description, copy.Description)
		require.EqualValues(t, original.Labels, copy.Labels)

		// Modify copy's labels slice
		copy.Labels[0] = "modified"
		copy.Labels = append(copy.Labels, "new-label")

		// Original should be unchanged (deep copy of slice)
		require.EqualValues(t, []string{"feature", "enhancement"}, original.Labels)
		require.EqualValues(t, []string{"modified", "enhancement", "new-label"}, copy.Labels)
	})

	t.Run("with all boolean flags set", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabCreateMergeRequestOptions{
			SourceBranchName:                "feature",
			TargetBranchName:                "main",
			Title:                           "Add feature",
			Description:                     "New feature",
			Labels:                          []string{"feature"},
			SquashEnabled:                   true,
			DeleteSourceBranchOnMerge:       true,
			FailIfMergeRequestAlreadyExists: true,
			AssignToSelf:                    true,
		}
		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.EqualValues(t, original.SquashEnabled, copy.SquashEnabled)
		require.EqualValues(t, original.DeleteSourceBranchOnMerge, copy.DeleteSourceBranchOnMerge)
		require.EqualValues(t, original.FailIfMergeRequestAlreadyExists, copy.FailIfMergeRequestAlreadyExists)
		require.EqualValues(t, original.AssignToSelf, copy.AssignToSelf)

		// Modify copy's booleans
		copy.SquashEnabled = false
		copy.DeleteSourceBranchOnMerge = false
		copy.FailIfMergeRequestAlreadyExists = false
		copy.AssignToSelf = false

		// Original should be unchanged
		require.True(t, original.SquashEnabled)
		require.True(t, original.DeleteSourceBranchOnMerge)
		require.True(t, original.FailIfMergeRequestAlreadyExists)
		require.True(t, original.AssignToSelf)
	})
}
