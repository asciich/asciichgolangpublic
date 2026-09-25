package gitgeneric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
)

// ---------- NewBranch ----------

func TestNewBranch_ReturnsUnsetBranch(t *testing.T) {
	branch := gitgeneric.NewBranch()

	require.NotNil(t, branch)

	name, err := branch.GetName()
	require.Error(t, err)
	require.EqualValues(t, "", name)
}

func TestNewBranch_ReturnsIndependentInstances(t *testing.T) {
	branch1 := gitgeneric.NewBranch()
	branch2 := gitgeneric.NewBranch()

	require.NotSame(t, branch1, branch2)

	require.NoError(t, branch1.SetName("main"))

	name, err := branch1.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", name)

	_, err = branch2.GetName()
	require.Error(t, err)
}

// ---------- GetName ----------

func TestBranch_GetName_UnsetReturnsError(t *testing.T) {
	branch := gitgeneric.Branch{}

	name, err := branch.GetName()

	require.Error(t, err)
	require.EqualValues(t, "", name)
}

func TestBranch_GetName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"main"},
		{"master"},
		{"feature/abc"},
		{"a"},
		{"release/v1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch := gitgeneric.Branch{}

			err := branch.SetName(tt.name)
			require.NoError(t, err)

			gotName, err := branch.GetName()
			require.NoError(t, err)
			require.EqualValues(t, tt.name, gotName)
		})
	}
}

// ---------- SetName ----------

func TestBranch_SetName_BlankStringReturnsError(t *testing.T) {
	tests := []struct {
		testCaseName string
		name         string
	}{
		{"empty", ""},
		{"single space", " "},
		{"multiple spaces", "   "},
		{"tab", "\t"},
		{"newline", "\n"},
		{"carriage return", "\r"},
		{"mixed whitespace", " \t\n\r "},
	}

	for _, tt := range tests {
		t.Run(tt.testCaseName, func(t *testing.T) {
			branch := gitgeneric.Branch{}

			err := branch.SetName(tt.name)
			require.Error(t, err)

			// Name must stay unset, so GetName still errors out.
			name, err := branch.GetName()
			require.Error(t, err)
			require.EqualValues(t, "", name)
		})
	}
}

func TestBranch_SetName_BlankStringDoesNotOverwriteExistingName(t *testing.T) {
	tests := []struct {
		testCaseName string
		name         string
	}{
		{"empty", ""},
		{"single space", " "},
		{"tab", "\t"},
		{"newline", "\n"},
		{"mixed whitespace", " \t\n\r "},
	}

	for _, tt := range tests {
		t.Run(tt.testCaseName, func(t *testing.T) {
			branch := gitgeneric.Branch{}

			err := branch.SetName("main")
			require.NoError(t, err)

			err = branch.SetName(tt.name)
			require.Error(t, err)

			name, err := branch.GetName()
			require.NoError(t, err)
			require.EqualValues(t, "main", name)
		})
	}
}

func TestBranch_SetName_KeepsSurroundingWhitespace(t *testing.T) {
	tests := []struct {
		testCaseName string
		name         string
	}{
		{"leading space", " main"},
		{"trailing space", "main "},
		{"surrounding spaces", " main "},
	}

	for _, tt := range tests {
		t.Run(tt.testCaseName, func(t *testing.T) {
			branch := gitgeneric.Branch{}

			err := branch.SetName(tt.name)
			require.NoError(t, err)

			name, err := branch.GetName()
			require.NoError(t, err)
			require.EqualValues(t, tt.name, name)
		})
	}
}

func TestBranch_SetName_Overwrite(t *testing.T) {
	tests := []struct {
		firstName  string
		secondName string
	}{
		{"main", "master"},
		{"feature/a", "feature/b"},
		{"x", "y"},
	}

	for _, tt := range tests {
		t.Run(tt.firstName+"->"+tt.secondName, func(t *testing.T) {
			branch := gitgeneric.Branch{}

			require.NoError(t, branch.SetName(tt.firstName))

			name, err := branch.GetName()
			require.NoError(t, err)
			require.EqualValues(t, tt.firstName, name)

			require.NoError(t, branch.SetName(tt.secondName))

			name, err = branch.GetName()
			require.NoError(t, err)
			require.EqualValues(t, tt.secondName, name)
		})
	}
}

// ---------- GetDeepCopy ----------

func TestBranch_GetDeepCopy_NotNil(t *testing.T) {
	branch := gitgeneric.NewBranch()
	require.NoError(t, branch.SetName("main"))

	deepCopy := branch.GetDeepCopy()

	require.NotNil(t, deepCopy)
}

func TestBranch_GetDeepCopy_CopiesName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"main"},
		{"master"},
		{"feature/abc"},
		{"a"},
		{"release/v1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch := gitgeneric.NewBranch()
			require.NoError(t, branch.SetName(tt.name))

			deepCopy := branch.GetDeepCopy()

			copiedName, err := deepCopy.GetName()
			require.NoError(t, err)
			require.EqualValues(t, tt.name, copiedName)
		})
	}
}

func TestBranch_GetDeepCopy_UnsetBranchStaysUnset(t *testing.T) {
	branch := gitgeneric.NewBranch()

	deepCopy := branch.GetDeepCopy()

	require.NotNil(t, deepCopy)

	name, err := deepCopy.GetName()
	require.Error(t, err)
	require.EqualValues(t, "", name)
}

func TestBranch_GetDeepCopy_ReturnsDifferentInstance(t *testing.T) {
	branch := gitgeneric.NewBranch()
	require.NoError(t, branch.SetName("main"))

	deepCopy := branch.GetDeepCopy()

	// The copy must not be the very same object as the original.
	require.NotSame(t, branch, deepCopy)
}

func TestBranch_GetDeepCopy_ChangingCopyDoesNotAffectOriginal(t *testing.T) {
	branch := gitgeneric.NewBranch()
	require.NoError(t, branch.SetName("main"))

	deepCopy := branch.GetDeepCopy()
	require.NoError(t, deepCopy.SetName("master"))

	originalName, err := branch.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", originalName)

	copiedName, err := deepCopy.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "master", copiedName)
}

func TestBranch_GetDeepCopy_ChangingOriginalDoesNotAffectCopy(t *testing.T) {
	branch := gitgeneric.NewBranch()
	require.NoError(t, branch.SetName("main"))

	deepCopy := branch.GetDeepCopy()
	require.NoError(t, branch.SetName("master"))

	copiedName, err := deepCopy.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", copiedName)

	originalName, err := branch.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "master", originalName)
}

func TestBranch_GetDeepCopy_OfDeepCopy(t *testing.T) {
	branch := gitgeneric.NewBranch()
	require.NoError(t, branch.SetName("main"))

	deepCopy := branch.GetDeepCopy().GetDeepCopy()

	copiedName, err := deepCopy.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", copiedName)

	require.NoError(t, deepCopy.SetName("master"))

	originalName, err := branch.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", originalName)
}

// ---------- Pointer usage ----------

func TestBranch_PointerReceiverWorksOnHeapAllocatedBranch(t *testing.T) {
	branch := new(gitgeneric.Branch)

	_, err := branch.GetName()
	require.Error(t, err)

	require.NoError(t, branch.SetName("main"))

	name, err := branch.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "main", name)
}
