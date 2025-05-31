package versionutils

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func CheckIsDateVersionString(versionString string) error {
	isVersionString := IsVersionString(versionString)

	if isVersionString {
		return nil
	} else {
		return tracederrors.TracedErrorf("'%s' is not a version string", versionString)
	}
}

func GetLatestVersionFromSlice(versions []Version) (latestVersion Version, err error) {
	for _, toCheck := range versions {
		if toCheck == nil {
			return nil, tracederrors.TracedErrorNilf(
				"toCheck is nil in '%v'",
				versions,
			)
		}

		if latestVersion == nil {
			latestVersion = toCheck
		}

		isNewer, err := toCheck.IsNewerThan(latestVersion)
		if err != nil {
			return nil, err
		}

		if isNewer {
			latestVersion = toCheck
		}
	}

	if latestVersion == nil {
		return nil, tracederrors.TracedErrorf("Unable to find latest version in '%v'", versions)
	}

	return latestVersion, nil
}

// Get the current date and time formated as version string 'YYYYmmdd_HHMMSS'.
//
// To get the current date and time as `versionutils.Version` use `versionutils.NewCurrentDateVersion()`.
func GetNewDateVersionString() (versionString string) {
	versionString = time.Now().Format("20060102_150405")
	return versionString
}

func ReadFromString(versionString string) (version Version, err error) {
	if !IsVersionString(versionString) {
		return nil, tracederrors.TracedErrorf("versionString '%s' is not a valid version string", versionString)
	}

	if IsDateVersionString(versionString) {
		return ReadDateVersionFromString(versionString)
	}

	if IsSemanticVersionString(versionString) {
		return ReadSemanticVersionFormString(versionString)
	}

	return nil, tracederrors.TracedErrorf("Not implemented for versionString='%s'", versionString)
}

func GetSoftwareVersionEnvVarName() (envVarName string) {
	return "SOFTWARE_VERSION"
}

func GetSoftwareVersionFromEnvVarOrEmptyStringIfUnset(verbose bool) (softwareVersion string) {
	envVarName := GetSoftwareVersionEnvVarName()

	softwareVersion = os.Getenv(envVarName)
	softwareVersion = strings.TrimSpace(softwareVersion)

	if softwareVersion == "" {
		if verbose {
			logging.LogInfof("Software version is not set in environment variable '%s'.", envVarName)
		}
		return ""
	} else {
		if verbose {
			logging.LogInfof(
				"Get software version from environment variable '%s' as '%s'.",
				envVarName,
				softwareVersion,
			)
		}
		return softwareVersion
	}
}

func GetVersionStringsFromStringSlice(input []string) (versionStrings []string) {
	versionStrings = []string{}

	for _, toCheck := range input {
		if IsVersionString(toCheck) {
			versionStrings = append(versionStrings, toCheck)
		}
	}

	return versionStrings
}

func GetVersionStringsFromVersionSlice(versions []Version) (versionStrings []string, err error) {
	if versions == nil {
		return nil, tracederrors.TracedErrorNil("versions")
	}

	versionStrings = []string{}
	for _, v := range versions {
		toAdd, err := v.GetAsString()
		if err != nil {
			return nil, err
		}

		versionStrings = append(versionStrings, toAdd)
	}

	return versionStrings, nil
}

func GetVersionsFromStringSlice(stringSlice []string) (versions []Version, err error) {
	if stringSlice == nil {
		return nil, tracederrors.TracedErrorNil("stringSlice")
	}

	versions = []Version{}
	for _, stringVersion := range stringSlice {
		toAdd, err := ReadFromString(stringVersion)
		if err != nil {
			return nil, err
		}

		versions = append(versions, toAdd)
	}

	return versions, nil
}

func IsDateVersionString(versionString string) (isVersionString bool) {
	regex := regexp.MustCompile("^[0-9]{8}_[0-9]{6}$")
	if regex.Match([]byte(versionString)) {
		return true
	}

	regex = regexp.MustCompile("^v[0-9]{8}_[0-9]{6}$")
	return regex.Match([]byte(versionString))
}

func IsSemanticVersionString(versionString string) (isSemanticVersionString bool) {
	regex := regexp.MustCompile("^[vV]{0,1}[0-9]{1,}\\.[0-9]{1,}\\.[0-9]{1,}$")
	isSemanticVersionString = regex.Match([]byte(versionString))
	return isSemanticVersionString
}

func IsVersionString(versionString string) (isVersionString bool) {
	if IsDateVersionString(versionString) {
		return true
	}

	return IsSemanticVersionString(versionString)
}

func ReturnNewerVersion(v1 Version, v2 Version) (newerVersion Version, err error) {
	if v1 == nil {
		return nil, tracederrors.TracedErrorNil("v1")
	}

	if v2 == nil {
		return nil, tracederrors.TracedErrorNil("v2")
	}

	isNewer, err := v1.IsNewerThan(v2)
	if err != nil {
		return nil, err
	}

	if isNewer {
		return v1, nil
	}

	return v2, nil
}

func SortStringSlice(versionStrings []string) (sorted []string, err error) {
	if versionStrings == nil {
		return nil, tracederrors.TracedErrorNil("versionStrings")
	}

	versions, err := GetVersionsFromStringSlice(versionStrings)
	if err != nil {
		return nil, err
	}

	versions, err = SortVersionSlice(versions)
	if err != nil {
		return nil, err
	}

	return GetVersionStringsFromVersionSlice(versions)
}

func SortVersionSlice(versions []Version) (sorted []Version, err error) {
	if versions == nil {
		return nil, tracederrors.TracedErrorNil("versions")
	}

	var errDuringSort error
	sort.Slice(
		versions,
		func(i int, j int) bool {
			isNewer, err := versions[i].IsNewerThan(versions[j])
			if err != nil {
				errDuringSort = err
				return false
			}

			return !isNewer
		},
	)

	if errDuringSort != nil {
		return nil, errDuringSort
	}

	return versions, nil
}
