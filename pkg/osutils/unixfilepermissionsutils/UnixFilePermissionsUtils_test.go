package unixfilepermissionsutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGetPermissionStringForAccessClass(t *testing.T) {
	tests := []struct {
		permission               int
		expectedPermissionString string
	}{
		{0, ""},
		{1, "x"},
		{2, "w"},
		{3, "wx"},
		{4, "r"},
		{5, "rx"},
		{6, "rw"},
		{7, "rwx"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.expectedPermissionString,
					MustGetPermissionStringForAccessClass(tt.permission),
				)
			},
		)
	}
}

func TestGetPermissionValueForAccessClassString(t *testing.T) {
	tests := []struct {
		permission         string
		expectedPermission int
	}{
		{"", 0},
		{"x", 1},
		{"w", 2},
		{"wx", 3},
		{"xw", 3},
		{"r", 4},
		{"rx", 5},
		{"xr", 5},
		{"rw", 6},
		{"wr", 6},
		{"rwx", 7},
		{"rxw", 7},
		{"xrw", 7},
		{"xwr", 7},
		{"wxr", 7},
		{"wrx", 7},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.expectedPermission,
					MustGetPermissionValueForAccessClassString(tt.permission),
				)
			},
		)
	}
}

func TestSplitPermissionValueInClasses(t *testing.T) {
	tests := []struct {
		permission     int
		expectedUser   int
		expectedGroup  int
		expectedOthers int
	}{
		{0o000, 0, 0, 0},
		{0o123, 1, 2, 3},
		{0o741, 7, 4, 1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("0o%o", tt.permission),
			func(t *testing.T) {
				user, group, others := MustSplitPermissionValueInClasses(tt.permission)

				require.EqualValues(t, tt.expectedUser, user)
				require.EqualValues(t, tt.expectedGroup, group)
				require.EqualValues(t, tt.expectedOthers, others)
			},
		)
	}
}

func TestMergeClassValues(t *testing.T) {
	tests := []struct {
		User               int
		Group              int
		Others             int
		ExpectedPermission int
	}{
		{0, 0, 0, 0o000},
		{1, 2, 3, 0o123},
		{7, 4, 1, 0o741},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("0o%o", tt.ExpectedPermission),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.ExpectedPermission,
					MustMergeClassValues(tt.User, tt.Group, tt.Others),
				)
			},
		)
	}
}

func TestMergeClassValuesAsString(t *testing.T) {
	tests := []struct {
		User               int
		Group              int
		Others             int
		ExpectedPermission string
	}{
		{0, 0, 0, "u=,g=,o="},
		{1, 2, 3, "u=x,g=w,o=wx"},
		{7, 4, 1, "u=rwx,g=r,o=x"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.ExpectedPermission,
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.ExpectedPermission,
					MustMergeClassValuesAsString(tt.User, tt.Group, tt.Others),
				)
			},
		)
	}
}

func TestSplitPermissionValueInClassPermissionStrings(t *testing.T) {
	tests := []struct {
		permission     int
		expectedUser   string
		expectedGroup  string
		expectedOthers string
	}{
		{0o000, "", "", ""},
		{0o123, "x", "w", "wx"},
		{0o741, "rwx", "r", "x"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			fmt.Sprintf("0o%o", tt.permission),
			func(t *testing.T) {
				user, group, others := MustSplitPermissionValueInClassPermissionStrings(tt.permission)

				require.EqualValues(t, tt.expectedUser, user)
				require.EqualValues(t, tt.expectedGroup, group)
				require.EqualValues(t, tt.expectedOthers, others)
			},
		)
	}
}

func TestGetPermissionString(t *testing.T) {
	tests := []struct {
		permission               int
		expectedPermissionString string
	}{
		{0o000, "u=,g=,o="},
		{0o100, "u=x,g=,o="},
		{0o010, "u=,g=x,o="},
		{0o001, "u=,g=,o=x"},
		{0o200, "u=w,g=,o="},
		{0o020, "u=,g=w,o="},
		{0o002, "u=,g=,o=w"},
		{0o300, "u=wx,g=,o="},
		{0o030, "u=,g=wx,o="},
		{0o003, "u=,g=,o=wx"},
		{0o400, "u=r,g=,o="},
		{0o040, "u=,g=r,o="},
		{0o004, "u=,g=,o=r"},
		{0o500, "u=rx,g=,o="},
		{0o050, "u=,g=rx,o="},
		{0o005, "u=,g=,o=rx"},
		{0o600, "u=rw,g=,o="},
		{0o060, "u=,g=rw,o="},
		{0o006, "u=,g=,o=rw"},
		{0o700, "u=rwx,g=,o="},
		{0o070, "u=,g=rwx,o="},
		{0o007, "u=,g=,o=rwx"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.expectedPermissionString,
					MustGetPermissionString(tt.permission),
				)
			},
		)
	}
}

func TestGetPermissionValue(t *testing.T) {
	tests := []struct {
		permission         string
		expectedPermission int
	}{
		{"", 0o000},
		{"u=,g=,o=", 0o000},
		{"u=,g=", 0o000},
		{"u=,g=,", 0o000},
		{"u=,g=,,,,,,,,,", 0o000},
		{"u=", 0o000},
		{"u=x,g=,o=", 0o100},
		{"u=x,g=", 0o100},
		{"u=x,o=", 0o100},
		{"u=x,", 0o100},
		{"u=x", 0o100},
		{"u=,g=x,o=", 0o010},
		{"g=x,o=", 0o010},
		{"g=x", 0o010},
		{"u=,g=,o=x", 0o001},
		{"o=x", 0o001},
		{"u=w,g=,o=", 0o200},
		{"u=,g=w,o=", 0o020},
		{"u=,g=,o=w", 0o002},
		{"u=wx,g=,o=", 0o300},
		{"u=,g=wx,o=", 0o030},
		{"u=,g=,o=wx", 0o003},
		{"u=r,g=,o=", 0o400},
		{"u=,g=r,o=", 0o040},
		{"u=,g=,o=r", 0o004},
		{"u=rx,g=,o=", 0o500},
		{"u=,g=rx,o=", 0o050},
		{"u=,g=,o=rx", 0o005},
		{"u=rw,g=,o=", 0o600},
		{"u=,g=rw,o=", 0o060},
		{"u=,g=,o=rw", 0o006},
		{"u=rwx,g=,o=", 0o700},
		{"u=,g=rwx,o=", 0o070},
		{"u=,g=,o=rwx", 0o007},
		{"u=x,g=r,o=rwx", 0o147},
		{"o=rwx,u=x,g=r", 0o147},
		{"g=r,o=rwx,u=x", 0o147},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.expectedPermission,
					MustGetPermissionsValue(tt.permission),
				)
			},
		)
	}
}
