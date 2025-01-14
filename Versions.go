package asciichgolangpublic

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type VersionsService struct {
}

func GetVersionByString(versionString string) (version Version, err error) {
	if versionString == "" {
		return nil, errors.TracedErrorEmptyString("version")
	}

	version, err = Versions().GetNewVersionByString(versionString)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func MustGetVersionByString(versionString string) (version Version) {
	version, err := GetVersionByString(versionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func NewVersionsService() (v *VersionsService) {
	return new(VersionsService)
}

func Versions() (v *VersionsService) {
	return NewVersionsService()
}

func (v *VersionsService) CheckDateVersionString(versionString string) (isVersionString bool, err error) {
	isVersionString = v.IsVersionString(versionString)

	if isVersionString {
		return true, nil
	} else {
		return false, errors.TracedErrorf("'%s' is not a version string", versionString)
	}
}

func (v *VersionsService) GetLatestVersionFromSlice(versions []Version) (latestVersion Version, err error) {
	for _, toCheck := range versions {
		if toCheck == nil {
			return nil, errors.TracedErrorNilf(
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
		return nil, errors.TracedErrorf("Unable to find latest version in '%v'", versions)
	}

	return latestVersion, nil
}

func (v *VersionsService) GetNewDateVersion() (version Version, err error) {
	versionString, err := v.GetNewDateVersionString()
	if err != nil {
		return nil, err
	}

	version, err = v.GetNewVersionByString(versionString)
	if err != nil {
		return nil, err
	}

	return version, err
}

func (v *VersionsService) GetNewDateVersionString() (versionString string, err error) {
	versionString = time.Now().Format("20060102_150405") // TODO use Time module from
	return versionString, nil
}

func (v *VersionsService) GetNewVersionByString(versionString string) (version Version, err error) {
	if !v.IsVersionString(versionString) {
		return nil, errors.TracedErrorf("versionString '%s' is not a valid version string", versionString)
	}

	if v.IsDateVersionString(versionString) {
		dateVersion := NewVersionDateVersion()
		err = dateVersion.SetVersion(versionString)
		if err != nil {
			return nil, err
		}

		return dateVersion, nil
	}

	if v.IsSemanticVersionString(versionString) {
		semanticVersion := NewVersionSemanticVersion()
		err = semanticVersion.SetVersionByString(versionString)
		if err != nil {
			return nil, err
		}

		return semanticVersion, nil
	}

	return nil, errors.TracedErrorf("Not implemented for versionString='%s'", versionString)
}

func (v *VersionsService) GetSoftwareVersionEnvVarName() (envVarName string) {
	return "SOFTWARE_VERSION"
}

func (v *VersionsService) GetSoftwareVersionFromEnvVarOrEmptyStringIfUnset(verbose bool) (softwareVersion string) {
	envVarName := Versions().GetSoftwareVersionEnvVarName()

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

func (v *VersionsService) GetVersionStringsFromStringSlice(input []string) (versionStrings []string) {
	versionStrings = []string{}

	for _, toCheck := range input {
		if v.IsVersionString(toCheck) {
			versionStrings = append(versionStrings, toCheck)
		}
	}

	return versionStrings
}

func (v *VersionsService) GetVersionStringsFromVersionSlice(versions []Version) (versionStrings []string, err error) {
	if versions == nil {
		return nil, errors.TracedErrorNil("versions")
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

func (v *VersionsService) GetVersionsFromStringSlice(stringSlice []string) (versions []Version, err error) {
	if stringSlice == nil {
		return nil, errors.TracedErrorNil("stringSlice")
	}

	versions = []Version{}
	for _, stringVersion := range stringSlice {
		toAdd, err := GetVersionByString(stringVersion)
		if err != nil {
			return nil, err
		}

		versions = append(versions, toAdd)
	}

	return versions, nil
}

func (v *VersionsService) IsDateVersionString(versionString string) (isVersionString bool) {
	regex := regexp.MustCompile("^[0-9]{8}_[0-9]{6}$")
	if regex.Match([]byte(versionString)) {
		return true
	}

	regex = regexp.MustCompile("^v[0-9]{8}_[0-9]{6}$")
	return regex.Match([]byte(versionString))
}

func (v *VersionsService) IsSemanticVersionString(versionString string) (isSemanticVersionString bool) {
	regex := regexp.MustCompile("^[vV]{0,1}[0-9]{1,}\\.[0-9]{1,}\\.[0-9]{1,}$")
	isSemanticVersionString = regex.Match([]byte(versionString))
	return isSemanticVersionString
}

func (v *VersionsService) IsVersionString(versionString string) (isVersionString bool) {
	if v.IsDateVersionString(versionString) {
		return true
	}

	return v.IsSemanticVersionString(versionString)
}

func (v *VersionsService) MustCheckDateVersionString(versionString string) (isVersionString bool) {
	isVersionString, err := v.CheckDateVersionString(versionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isVersionString
}

func (v *VersionsService) MustGetLatestVersionFromSlice(versions []Version) (latestVersion Version) {
	latestVersion, err := v.GetLatestVersionFromSlice(versions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestVersion
}

func (v *VersionsService) MustGetNewDateVersion() (version Version) {
	version, err := v.GetNewDateVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (v *VersionsService) MustGetNewDateVersionString() (versionString string) {
	versionString, err := v.GetNewDateVersionString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionString
}

func (v *VersionsService) MustGetNewVersionByString(versionString string) (version Version) {
	version, err := v.GetNewVersionByString(versionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (v *VersionsService) MustGetVersionStringsFromVersionSlice(versions []Version) (versionStrings []string) {
	versionStrings, err := v.GetVersionStringsFromVersionSlice(versions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionStrings
}

func (v *VersionsService) MustGetVersionsFromStringSlice(stringSlice []string) (versions []Version) {
	versions, err := v.GetVersionsFromStringSlice(stringSlice)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versions
}

func (v *VersionsService) MustReturnNewerVersion(v1 Version, v2 Version) (newerVersion Version) {
	newerVersion, err := v.ReturnNewerVersion(v1, v2)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newerVersion
}

func (v *VersionsService) MustSortStringSlice(versionStrings []string) (sorted []string) {
	sorted, err := v.SortStringSlice(versionStrings)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sorted
}

func (v *VersionsService) MustSortVersionSlice(versions []Version) (sorted []Version) {
	sorted, err := v.SortVersionSlice(versions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sorted
}

func (v *VersionsService) ReturnNewerVersion(v1 Version, v2 Version) (newerVersion Version, err error) {
	if v1 == nil {
		return nil, errors.TracedErrorNil("v1")
	}

	if v2 == nil {
		return nil, errors.TracedErrorNil("v2")
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

func (v *VersionsService) SortStringSlice(versionStrings []string) (sorted []string, err error) {
	if versionStrings == nil {
		return nil, errors.TracedErrorNil("versionStrings")
	}

	versions, err := v.GetVersionsFromStringSlice(versionStrings)
	if err != nil {
		return nil, err
	}

	versions, err = v.SortVersionSlice(versions)
	if err != nil {
		return nil, err
	}

	return v.GetVersionStringsFromVersionSlice(versions)
}

func (v *VersionsService) SortVersionSlice(versions []Version) (sorted []Version, err error) {
	if versions == nil {
		return nil, errors.TracedErrorNil("versions")
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
