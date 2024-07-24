package asciichgolangpublic

import (
	"os"
	"regexp"
	"strings"
	"time"
)

type VersionsService struct {
}

func GetVersionByString(versionString string) (version Version, err error) {
	if versionString == "" {
		return nil, TracedErrorEmptyString("version")
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
		LogGoErrorFatal(err)
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
		return false, TracedErrorf("'%s' is not a version string", versionString)
	}
}

func (v *VersionsService) GetLatestVersionFromSlice(versions []Version) (latestVersion Version, err error) {
	for _, toCheck := range versions {
		if toCheck == nil {
			return nil, TracedErrorNilf(
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
		return nil, TracedErrorf("Unable to find latest version in '%v'", versions)
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
		return nil, TracedErrorf("versionString '%s' is not a valid version string", versionString)
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

	return nil, TracedErrorf("Not implemented for versionString='%s'", versionString)
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
			LogInfof("Software version is not set in environment variable '%s'.", envVarName)
		}
		return ""
	} else {
		if verbose {
			LogInfof(
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

func (v *VersionsService) GetVersionsFromStringSlice(stringSlice []string) (versions []Version, err error) {
	if stringSlice == nil {
		return nil, TracedErrorNil("stringSlice")
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
	isVersionString = regex.Match([]byte(versionString))
	return isVersionString
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
		LogGoErrorFatal(err)
	}

	return isVersionString
}

func (v *VersionsService) MustGetLatestVersionFromSlice(versions []Version) (latestVersion Version) {
	latestVersion, err := v.GetLatestVersionFromSlice(versions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestVersion
}

func (v *VersionsService) MustGetNewDateVersion() (version Version) {
	version, err := v.GetNewDateVersion()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return version
}

func (v *VersionsService) MustGetNewDateVersionString() (versionString string) {
	versionString, err := v.GetNewDateVersionString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionString
}

func (v *VersionsService) MustGetNewVersionByString(versionString string) (version Version) {
	version, err := v.GetNewVersionByString(versionString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return version
}

func (v *VersionsService) MustGetVersionsFromStringSlice(stringSlice []string) (versions []Version) {
	versions, err := v.GetVersionsFromStringSlice(stringSlice)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versions
}
