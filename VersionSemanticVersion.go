package asciichgolangpublic

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionSemanticVersion struct {
	major int
	minor int
	patch int
}

func NewVersionSemanticVersion() (v *VersionSemanticVersion) {
	return new(VersionSemanticVersion)
}

func (v *VersionSemanticVersion) Equals(other Version) (isEqual bool) {
	if other == nil {
		return false
	}

	otherSemanticVersion, ok := other.(*VersionSemanticVersion)
	if !ok {
		return false
	}

	if v.major != otherSemanticVersion.major {
		return false
	}

	if v.minor != otherSemanticVersion.minor {
		return false
	}

	if v.patch != otherSemanticVersion.patch {
		return false
	}

	return true
}

func (v *VersionSemanticVersion) GetAsString() (versionString string, err error) {
	versionString, err = v.GetAsStringWithoutLeadingV()
	if err != nil {
		return "", err
	}

	return "v" + versionString, nil
}

func (v *VersionSemanticVersion) GetAsStringWithoutLeadingV() (versionString string, err error) {
	major, minor, patch, err := v.GetMajorMinorPatch()
	if err != nil {
		return "", err
	}

	versionString = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	return versionString, nil
}

func (v *VersionSemanticVersion) GetMajor() (major int, err error) {
	if v.major < 0 {
		return -1, TracedError("major not set")
	}

	return v.major, nil
}

func (v *VersionSemanticVersion) GetMajorMinorPatch() (major int, minor int, patch int, err error) {
	major, err = v.GetMajor()
	if err != nil {
		return -1, -1, -1, err
	}

	minor, err = v.GetMinor()
	if err != nil {
		return -1, -1, -1, err
	}

	patch, err = v.GetPatch()
	if err != nil {
		return -1, -1, -1, err
	}

	return major, minor, patch, nil
}

func (v *VersionSemanticVersion) GetMinor() (minor int, err error) {
	if v.minor < 0 {
		return -1, TracedError("minor not set")
	}

	return v.minor, nil
}

func (v *VersionSemanticVersion) GetNextVersion(versionType string) (nextVersion Version, err error) {
	if versionType == "" {
		return nil, TracedErrorEmptyString("versionType")
	}

	major, minor, patch, err := v.GetMajorMinorPatch()
	if err != nil {
		return nil, err
	}

	if versionType == "patch" {
		patch += 1
	} else if versionType == "minor" {
		minor += 1
		patch = 0
	} else if versionType == "major" {
		major += 1
		minor = 0
		patch = 0
	} else {
		return nil, TracedErrorf(
			"Unknown versionType='%s'",
			versionType,
		)
	}

	nextSemanticVersion := NewVersionSemanticVersion()
	err = nextSemanticVersion.Set(major, minor, patch)
	if err != nil {
		return nil, err
	}

	return nextSemanticVersion, nil
}

func (v *VersionSemanticVersion) GetPatch() (patch int, err error) {
	if v.patch < 0 {
		return -1, TracedError("patch not set")
	}

	return v.patch, nil
}

func (v *VersionSemanticVersion) IsNewerThan(other Version) (isNewerThan bool, err error) {
	if other == nil {
		return false, TracedErrorNil("other")
	}

	otherSemanticVersion, ok := other.(*VersionSemanticVersion)
	if !ok {
		return false, TracedErrorf(
			"Non compareable versions '%v' and '%v'",
			v,
			other,
		)
	}

	thisMajor, thisMinor, thisPatch, err := v.GetMajorMinorPatch()
	if err != nil {
		return false, err
	}

	otherMajor, otherMinor, otherPatch, err := otherSemanticVersion.GetMajorMinorPatch()
	if err != nil {
		return false, err
	}

	if thisMajor > otherMajor {
		return true, nil
	}
	if thisMajor < otherMajor {
		return false, nil
	}

	if thisMinor > otherMinor {
		return true, nil
	}
	if thisMinor < otherMinor {
		return false, nil
	}

	if thisPatch > otherPatch {
		return true, nil
	}

	return false, nil
}

func (v *VersionSemanticVersion) IsSemanticVersion() (isSemanticVersion bool) {
	return true
}

func (v *VersionSemanticVersion) MustGetAsString() (versionString string) {
	versionString, err := v.GetAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionString
}

func (v *VersionSemanticVersion) MustGetAsStringWithoutLeadingV() (versionString string) {
	versionString, err := v.GetAsStringWithoutLeadingV()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionString
}

func (v *VersionSemanticVersion) MustGetMajor() (major int) {
	major, err := v.GetMajor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return major
}

func (v *VersionSemanticVersion) MustGetMajorMinorPatch() (major int, minor int, patch int) {
	major, minor, patch, err := v.GetMajorMinorPatch()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return major, minor, patch
}

func (v *VersionSemanticVersion) MustGetMinor() (minor int) {
	minor, err := v.GetMinor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return minor
}

func (v *VersionSemanticVersion) MustGetNextVersion(versionType string) (nextVersion Version) {
	nextVersion, err := v.GetNextVersion(versionType)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nextVersion
}

func (v *VersionSemanticVersion) MustGetPatch() (patch int) {
	patch, err := v.GetPatch()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return patch
}

func (v *VersionSemanticVersion) MustIsNewerThan(other Version) (isNewerThan bool) {
	isNewerThan, err := v.IsNewerThan(other)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isNewerThan
}

func (v *VersionSemanticVersion) MustSet(major int, minor int, patch int) {
	err := v.Set(major, minor, patch)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) MustSetMajor(major int) {
	err := v.SetMajor(major)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) MustSetMajorMinorPatch(major int, minor int, patch int) {
	err := v.SetMajorMinorPatch(major, minor, patch)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) MustSetMinor(minor int) {
	err := v.SetMinor(minor)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) MustSetPatch(patch int) {
	err := v.SetPatch(patch)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) MustSetVersionByString(version string) {
	err := v.SetVersionByString(version)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (v *VersionSemanticVersion) Set(major int, minor int, patch int) (err error) {
	err = v.SetMajor(major)
	if err != nil {
		return err
	}

	err = v.SetMinor(minor)
	if err != nil {
		return err
	}

	err = v.SetPatch(patch)
	if err != nil {
		return err
	}

	return nil
}

func (v *VersionSemanticVersion) SetMajor(major int) (err error) {
	if major < 0 {
		return TracedErrorf("Invalid value '%d' for major", major)
	}

	v.major = major

	return nil
}

func (v *VersionSemanticVersion) SetMajorMinorPatch(major int, minor int, patch int) (err error) {
	err = v.SetMajor(major)
	if err != nil {
		return err
	}

	err = v.SetMinor(minor)
	if err != nil {
		return err
	}

	err = v.SetPatch(patch)
	if err != nil {
		return err
	}

	return nil
}

func (v *VersionSemanticVersion) SetMinor(minor int) (err error) {
	if minor < 0 {
		return TracedErrorf("Invalid value '%d' for minor", minor)
	}

	v.minor = minor

	return nil
}

func (v *VersionSemanticVersion) SetPatch(patch int) (err error) {
	if patch < 0 {
		return TracedErrorf("Invalid value '%d' for patch", patch)
	}

	v.patch = patch

	return nil
}

func (v *VersionSemanticVersion) SetVersionByString(version string) (err error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return TracedErrorEmptyString("version")
	}

	version = Strings().TrimPrefixIgnoreCase(version, "v")

	splitted := strings.Split(version, ".")
	if len(splitted) != 3 {
		return TracedErrorf(
			"Unexepected number of spitted elements '%v' for version string '%s'",
			splitted,
			version,
		)
	}

	major, err := strconv.Atoi(splitted[0])
	if err != nil {
		return TracedErrorf("Unable to parse major '%s': '%w'", splitted[0], err)
	}

	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		return TracedErrorf("Unable to parse minor '%s': '%w'", splitted[1], err)
	}

	patch, err := strconv.Atoi(splitted[2])
	if err != nil {
		return TracedErrorf("Unable to parse patch '%s': '%w'", splitted[2], err)
	}

	err = v.SetMajorMinorPatch(major, minor, patch)
	if err != nil {
		return err
	}

	return nil
}
