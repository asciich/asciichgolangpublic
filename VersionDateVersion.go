package asciichgolangpublic

import "sort"

type VersionDateVersion struct {
	version string
}

func NewVersionDateVersion() (v *VersionDateVersion) {
	return new(VersionDateVersion)
}

func (v *VersionDateVersion) GetNextVersion(versionType string) (nextVersion Version, err error) {
	return Versions().GetNewDateVersion()
}

func (v *VersionDateVersion) GetVersion() (version string, err error) {
	if v.version == "" {
		return "", TracedErrorf("version not set")
	}

	return v.version, nil
}

func (v *VersionDateVersion) MustGetNextVersion(versionType string) (nextVersion Version) {
	nextVersion, err := v.GetNextVersion(versionType)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nextVersion
}

func (v *VersionDateVersion) MustGetVersion() (version string) {
	version, err := v.GetVersion()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return version
}

func (v *VersionDateVersion) MustIsNewerThan(other Version) (isNewerThan bool) {
	isNewerThan, err := v.IsNewerThan(other)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isNewerThan
}

func (v *VersionDateVersion) MustSetVersion(version string) {
	err := v.SetVersion(version)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionDateVersion) SetVersion(version string) (err error) {
	if version == "" {
		return TracedErrorf("version is empty string")
	}

	v.version = version

	return nil
}

func (v VersionDateVersion) Equals(other Version) (isEqual bool) {
	if other == nil {
		return false
	}

	otherDateVersionPtr, ok := other.(*VersionDateVersion)
	if ok {
		return v.version == otherDateVersionPtr.version
	}

	otherDateVersion, ok := other.(*VersionDateVersion)
	if ok {
		return v.version == otherDateVersion.version
	}

	return false
}

func (v VersionDateVersion) GetAsString() (version string, err error) {
	return v.GetVersion()
}

func (v VersionDateVersion) IsNewerThan(other Version) (isNewerThan bool, err error) {
	if other == nil {
		return false, TracedErrorNil("other")
	}

	otherDateVersion, ok := other.(*VersionDateVersion)
	if !ok {
		return false, TracedErrorf(
			"Incompatible versions to compare: '%s' and other '%s'",
			v,
			other,
		)
	}

	thisVersionString, err := v.GetAsString()
	if err != nil {
		return false, err
	}

	otherVersionString, err := otherDateVersion.GetAsString()
	if err != nil {
		return false, err
	}

	if thisVersionString == otherVersionString {
		return false, nil
	}

	sorted := []string{thisVersionString, otherVersionString}
	sort.Strings(sorted)
	isNewerThan = sorted[1] == thisVersionString
	return isNewerThan, nil
}

func (v VersionDateVersion) IsSemanticVersion() (isSemanticVersion bool) {
	return false
}

func (v VersionDateVersion) MustGetAsString() (version string) {
	version, err := v.GetAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return version
}
