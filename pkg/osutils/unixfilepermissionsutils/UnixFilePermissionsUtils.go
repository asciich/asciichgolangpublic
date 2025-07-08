package unixfilepermissionsutils

import (
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func MustGetPermissionString(permission int) (permissionString string) {
	permissionString, err := GetPermissionString(permission)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permissionString
}

func GetPermissionString(permission int) (permissionsString string, err error) {
	if permission < 0 || permission > 0o777 {
		return "", tracederrors.TracedErrorf("Invalid permissions value: '%d'", permission)
	}

	user, group, other, err := SplitPermissionValueInClasses(permission)
	if err != nil {
		return "", err
	}

	permissionsString, err = MergeClassValuesAsString(user, group, other)
	if err != nil {
		return "", err
	}

	return permissionsString, nil
}

func MustSplitPermissionValueInClassPermissionStrings(permission int) (user string, group string, other string) {
	user, group, other, err := SplitPermissionValueInClassPermissionStrings(permission)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return user, group, other
}

func SplitPermissionValueInClassPermissionStrings(permission int) (user string, group string, other string, err error) {
	userValue, groupValue, otherValue, err := SplitPermissionValueInClasses(permission)
	if err != nil {
		return "", "", "", err
	}

	user, err = GetPermissionStringForAccessClass(userValue)
	if err != nil {
		return "", "", "", err
	}

	group, err = GetPermissionStringForAccessClass(groupValue)
	if err != nil {
		return "", "", "", err
	}

	other, err = GetPermissionStringForAccessClass(otherValue)
	if err != nil {
		return "", "", "", err
	}

	return user, group, other, nil
}

func MustMergeClassValuesAsString(user int, group int, other int) (permission string) {
	permission, err := MergeClassValuesAsString(user, group, other)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permission
}

func MergeClassValuesAsString(user int, group int, other int) (permission string, err error) {
	userString, err := GetPermissionStringForAccessClass(user)
	if err != nil {
		return "", err
	}

	groupString, err := GetPermissionStringForAccessClass(group)
	if err != nil {
		return "", err
	}

	otherString, err := GetPermissionStringForAccessClass(other)
	if err != nil {
		return "", err
	}

	permission = fmt.Sprintf("u=%s,g=%s,o=%s", userString, groupString, otherString)

	return permission, nil
}

func MustMergeClassValues(user int, group int, other int) (permission int) {
	permission, err := MergeClassValues(user, group, other)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permission
}

func MergeClassValues(user int, group int, other int) (permission int, err error) {
	if user < 0 && user > 7 {
		return 0, tracederrors.TracedErrorf("Invalid user access value: '%d'", user)
	}

	if group < 0 && group > 7 {
		return 0, tracederrors.TracedErrorf("Invalid group access value: '%d'", group)
	}

	if other < 0 && other > 7 {
		return 0, tracederrors.TracedErrorf("Invalid user access value: '%d'", other)
	}

	permission = user<<6 | group<<3 | other

	return permission, nil
}

func MustSplitPermissionValueInClasses(permission int) (user int, group int, other int) {
	user, group, other, err := SplitPermissionValueInClasses(permission)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return user, group, other
}

func SplitPermissionValueInClasses(permission int) (user int, group int, other int, err error) {
	if permission < 0 || permission > 0o777 {
		return 0, 0, 0, tracederrors.TracedErrorf("Invalid permissions value: '%d'", permission)
	}

	user = permission >> 6 & 0o7
	if user < 0 && user > 7 {
		return 0, 0, 0, tracederrors.TracedErrorf("Invalid permissions value: '%d', user access '%d' is invalid.", permission, user)
	}

	group = permission >> 3 & 0o7
	if group < 0 && group > 7 {
		return 0, 0, 0, tracederrors.TracedErrorf("Invalid permissions value: '%d', group access '%d' is invalid.", permission, group)
	}

	other = permission & 0o7
	if other < 0 && other > 7 {
		return 0, 0, 0, tracederrors.TracedErrorf("Invalid permissions value: '%d', other access '%d' is invalid.", permission, other)
	}

	return user, group, other, nil
}

// Get the permission string for a single access class (user, group or other).
func MustGetPermissionStringForAccessClass(permission int) (permissionString string) {
	permissionString, err := GetPermissionStringForAccessClass(permission)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permissionString
}

// Get the permission string for a single access class (user, group or other).
func GetPermissionStringForAccessClass(permission int) (permissionsString string, err error) {
	lookup := map[int]string{
		0: "",
		1: "x",
		2: "w",
		3: "wx",
		4: "r",
		5: "rx",
		6: "rw",
		7: "rwx",
	}

	permissionsString, ok := lookup[permission]
	if !ok {
		return "", tracederrors.TracedErrorf("Unknown permission value '%d'", err)
	}

	return permissionsString, nil
}

func MustGetPermissionValueForAccessClassString(permissionString string) (permission int) {
	permission, err := GetPermissionValueForAccessClassString(permissionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permission
}

func GetPermissionValueForAccessClassString(permissionString string) (permission int, err error) {
	lookup := map[string]int{
		"":    0,
		"x":   1,
		"w":   2,
		"wx":  3,
		"xw":  3,
		"r":   4,
		"rx":  5,
		"xr":  5,
		"rw":  6,
		"wr":  6,
		"rwx": 7,
		"rxw": 7,
		"xrw": 7,
		"xwr": 7,
		"wrx": 7,
		"wxr": 7,
	}

	permission, ok := lookup[permissionString]
	if !ok {
		return 0, tracederrors.TracedErrorf("Unknown permsssion string for a single access class: '%s'", permissionString)
	}

	return permission, nil
}

func MustGetPermissionsValue(permissionsString string) (permissions int) {
	permissions, err := GetPermissionsValue(permissionsString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permissions
}

func GetPermissionsValue(permissionsString string) (permission int, err error) {
	var user, group, other int

	for _, part := range strings.Split(permissionsString, ",") {
		part = strings.TrimSpace(part)
		part = stringsutils.RepeatReplaceAll(part, " ", "")

		if part == "" {
			continue
		}

		if strings.HasPrefix(part, "u=") {
			user, err = GetPermissionValueForAccessClassString(strings.TrimPrefix(part, "u="))
			if err != nil {
				return 0, err
			}
			continue
		}

		if strings.HasPrefix(part, "g=") {
			group, err = GetPermissionValueForAccessClassString(strings.TrimPrefix(part, "g="))
			if err != nil {
				return 0, err
			}
			continue
		}

		if strings.HasPrefix(part, "o=") {
			other, err = GetPermissionValueForAccessClassString(strings.TrimPrefix(part, "o="))
			if err != nil {
				return 0, err
			}
			continue
		}

		return 0, tracederrors.TracedErrorf("Unexpected part '%s' in permission string '%s'", part, permissionsString)
	}

	permission, err = MergeClassValues(user, group, other)
	if err != nil {
		return 0, err
	}

	return permission, nil
}
