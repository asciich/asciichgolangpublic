package unixfilepermissionsutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeClassValues_InvalidUserValue(t *testing.T) {
	tests := []struct {
		User   int
		Group  int
		Others int
		name   string
	}{
		{-1, 0, 0, "negative_user"},
		{-5, 0, 0, "negative_user_large"},
		{8, 0, 0, "user_too_large"},
		{10, 0, 0, "user_much_too_large"},
		{100, 0, 0, "user_way_too_large"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				permission, err := MergeClassValues(tt.User, tt.Group, tt.Others)
				require.Error(t, err)
				require.EqualValues(t, 0, permission)
				require.Contains(t, err.Error(), "Invalid user access value")
			},
		)
	}
}

func TestMergeClassValues_InvalidGroupValue(t *testing.T) {
	tests := []struct {
		User   int
		Group  int
		Others int
		name   string
	}{
		{0, -1, 0, "negative_group"},
		{0, -5, 0, "negative_group_large"},
		{0, 8, 0, "group_too_large"},
		{0, 10, 0, "group_much_too_large"},
		{0, 100, 0, "group_way_too_large"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				permission, err := MergeClassValues(tt.User, tt.Group, tt.Others)
				require.Error(t, err)
				require.EqualValues(t, 0, permission)
				require.Contains(t, err.Error(), "Invalid group access value")
			},
		)
	}
}

func TestMergeClassValues_InvalidOtherValue(t *testing.T) {
	tests := []struct {
		User   int
		Group  int
		Others int
		name   string
	}{
		{0, 0, -1, "negative_other"},
		{0, 0, -5, "negative_other_large"},
		{0, 0, 8, "other_too_large"},
		{0, 0, 10, "other_much_too_large"},
		{0, 0, 100, "other_way_too_large"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				permission, err := MergeClassValues(tt.User, tt.Group, tt.Others)
				require.Error(t, err)
				require.EqualValues(t, 0, permission)
				require.Contains(t, err.Error(), "Invalid user access value")
			},
		)
	}
}

func TestMergeClassValues_ValidValues(t *testing.T) {
	tests := []struct {
		User               int
		Group              int
		Others             int
		ExpectedPermission int
		name               string
	}{
		{0, 0, 0, 0o000, "all_zeros"},
		{1, 2, 3, 0o123, "small_values"},
		{7, 4, 1, 0o741, "max_user"},
		{0, 7, 0, 0o070, "max_group"},
		{0, 0, 7, 0o007, "max_other"},
		{7, 7, 7, 0o777, "all_max"},
		{4, 5, 6, 0o456, "mixed_values"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				permission, err := MergeClassValues(tt.User, tt.Group, tt.Others)
				require.NoError(t, err)
				require.EqualValues(t, tt.ExpectedPermission, permission)
			},
		)
	}
}

func TestSplitPermissionValueInClasses_InvalidPermissionValue(t *testing.T) {
	tests := []struct {
		permission int
		name       string
	}{
		{-1, "negative_permission"},
		{-5, "negative_permission_large"},
		{-100, "negative_permission_way_large"},
		{0o1000, "permission_too_large"},
		{0o7777, "permission_much_too_large"},
		{1000, "permission_decimal_too_large"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				user, group, other, err := SplitPermissionValueInClasses(tt.permission)
				require.Error(t, err)
				require.EqualValues(t, 0, user)
				require.EqualValues(t, 0, group)
				require.EqualValues(t, 0, other)
				require.Contains(t, err.Error(), "Invalid permissions value")
			},
		)
	}
}

func TestSplitPermissionValueInClasses_ValidPermissionValues(t *testing.T) {
	tests := []struct {
		permission     int
		expectedUser   int
		expectedGroup  int
		expectedOthers int
		name           string
	}{
		{0o000, 0, 0, 0, "all_zeros"},
		{0o123, 1, 2, 3, "small_values"},
		{0o741, 7, 4, 1, "max_user"},
		{0o070, 0, 7, 0, "max_group"},
		{0o007, 0, 0, 7, "max_other"},
		{0o777, 7, 7, 7, "all_max"},
		{0o456, 4, 5, 6, "mixed_values"},
		{0o755, 7, 5, 5, "common_755"},
		{0o644, 6, 4, 4, "common_644"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				user, group, other, err := SplitPermissionValueInClasses(tt.permission)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedUser, user)
				require.EqualValues(t, tt.expectedGroup, group)
				require.EqualValues(t, tt.expectedOthers, other)
			},
		)
	}
}

func TestMergeClassValues_BoundaryValues(t *testing.T) {
	t.Run("user_at_lower_boundary", func(t *testing.T) {
		_, err := MergeClassValues(-1, 0, 0)
		require.Error(t, err, "Should reject user=-1")

		permission, err := MergeClassValues(0, 0, 0)
		require.NoError(t, err, "Should accept user=0")
		require.EqualValues(t, 0o000, permission)
	})

	t.Run("user_at_upper_boundary", func(t *testing.T) {
		permission, err := MergeClassValues(7, 0, 0)
		require.NoError(t, err, "Should accept user=7")
		require.EqualValues(t, 0o700, permission)

		_, err = MergeClassValues(8, 0, 0)
		require.Error(t, err, "Should reject user=8")
	})

	t.Run("group_at_lower_boundary", func(t *testing.T) {
		_, err := MergeClassValues(0, -1, 0)
		require.Error(t, err, "Should reject group=-1")

		permission, err := MergeClassValues(0, 0, 0)
		require.NoError(t, err, "Should accept group=0")
		require.EqualValues(t, 0o000, permission)
	})

	t.Run("group_at_upper_boundary", func(t *testing.T) {
		permission, err := MergeClassValues(0, 7, 0)
		require.NoError(t, err, "Should accept group=7")
		require.EqualValues(t, 0o070, permission)

		_, err = MergeClassValues(0, 8, 0)
		require.Error(t, err, "Should reject group=8")
	})

	t.Run("other_at_lower_boundary", func(t *testing.T) {
		_, err := MergeClassValues(0, 0, -1)
		require.Error(t, err, "Should reject other=-1")

		permission, err := MergeClassValues(0, 0, 0)
		require.NoError(t, err, "Should accept other=0")
		require.EqualValues(t, 0o000, permission)
	})

	t.Run("other_at_upper_boundary", func(t *testing.T) {
		permission, err := MergeClassValues(0, 0, 7)
		require.NoError(t, err, "Should accept other=7")
		require.EqualValues(t, 0o007, permission)

		_, err = MergeClassValues(0, 0, 8)
		require.Error(t, err, "Should reject other=8")
	})
}

func TestSplitPermissionValueInClasses_BoundaryValues(t *testing.T) {
	t.Run("permission_at_lower_boundary", func(t *testing.T) {
		_, _, _, err := SplitPermissionValueInClasses(-1)
		require.Error(t, err, "Should reject permission=-1")

		user, group, other, err := SplitPermissionValueInClasses(0)
		require.NoError(t, err, "Should accept permission=0")
		require.EqualValues(t, 0, user)
		require.EqualValues(t, 0, group)
		require.EqualValues(t, 0, other)
	})

	t.Run("permission_at_upper_boundary", func(t *testing.T) {
		user, group, other, err := SplitPermissionValueInClasses(0o777)
		require.NoError(t, err, "Should accept permission=0o777")
		require.EqualValues(t, 7, user)
		require.EqualValues(t, 7, group)
		require.EqualValues(t, 7, other)

		_, _, _, err = SplitPermissionValueInClasses(0o1000)
		require.Error(t, err, "Should reject permission=0o1000")
	})
}
