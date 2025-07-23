package versionutils

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/stringsutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type SemanticVersion struct {
	major int
	minor int
	patch int
}

func ReadSemanticVersionFormString(versionString string) (*SemanticVersion, error) {
	semanticVersion := NewVersionSemanticVersion()
	err := semanticVersion.SetVersionByString(versionString)
	if err != nil {
		return nil, err
	}

	return semanticVersion, nil
}

func NewVersionSemanticVersion() (v *SemanticVersion) {
	return new(SemanticVersion)
}

func (v *SemanticVersion) Equals(other Version) (isEqual bool) {
	if other == nil {
		return false
	}

	otherSemanticVersion, ok := other.(*SemanticVersion)
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

func (v *SemanticVersion) GetAsString() (versionString string, err error) {
	versionString, err = v.GetAsStringWithoutLeadingV()
	if err != nil {
		return "", err
	}

	return "v" + versionString, nil
}

func (v *SemanticVersion) GetAsStringWithoutLeadingV() (versionString string, err error) {
	major, minor, patch, err := v.GetMajorMinorPatch()
	if err != nil {
		return "", err
	}

	versionString = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	return versionString, nil
}

func (v *SemanticVersion) GetMajor() (major int, err error) {
	if v.major < 0 {
		return -1, tracederrors.TracedError("major not set")
	}

	return v.major, nil
}

func (v *SemanticVersion) GetMajorMinorPatch() (major int, minor int, patch int, err error) {
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

func (v *SemanticVersion) GetMinor() (minor int, err error) {
	if v.minor < 0 {
		return -1, tracederrors.TracedError("minor not set")
	}

	return v.minor, nil
}

func (v *SemanticVersion) GetNextVersion(versionType string) (nextVersion Version, err error) {
	if versionType == "" {
		return nil, tracederrors.TracedErrorEmptyString("versionType")
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
		return nil, tracederrors.TracedErrorf(
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

func (v *SemanticVersion) GetPatch() (patch int, err error) {
	if v.patch < 0 {
		return -1, tracederrors.TracedError("patch not set")
	}

	return v.patch, nil
}

func (v *SemanticVersion) IsNewerThan(other Version) (isNewerThan bool, err error) {
	if other == nil {
		return false, tracederrors.TracedErrorNil("other")
	}

	otherSemanticVersion, ok := other.(*SemanticVersion)
	if !ok {
		return false, tracederrors.TracedErrorf(
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

func (v *SemanticVersion) IsSemanticVersion() (isSemanticVersion bool) {
	return true
}

func (v *SemanticVersion) Set(major int, minor int, patch int) (err error) {
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

func (v *SemanticVersion) SetMajor(major int) (err error) {
	if major < 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for major", major)
	}

	v.major = major

	return nil
}

func (v *SemanticVersion) SetMajorMinorPatch(major int, minor int, patch int) (err error) {
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

func (v *SemanticVersion) SetMinor(minor int) (err error) {
	if minor < 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for minor", minor)
	}

	v.minor = minor

	return nil
}

func (v *SemanticVersion) SetPatch(patch int) (err error) {
	if patch < 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for patch", patch)
	}

	v.patch = patch

	return nil
}

func (v *SemanticVersion) SetVersionByString(version string) (err error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return tracederrors.TracedErrorEmptyString("version")
	}

	version = stringsutils.TrimPrefixIgnoreCase(version, "v")

	splitted := strings.Split(version, ".")
	if len(splitted) != 3 {
		return tracederrors.TracedErrorf(
			"Unexepected number of spitted elements '%v' for version string '%s'",
			splitted,
			version,
		)
	}

	major, err := strconv.Atoi(splitted[0])
	if err != nil {
		return tracederrors.TracedErrorf("Unable to parse major '%s': '%w'", splitted[0], err)
	}

	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		return tracederrors.TracedErrorf("Unable to parse minor '%s': '%w'", splitted[1], err)
	}

	patch, err := strconv.Atoi(splitted[2])
	if err != nil {
		return tracederrors.TracedErrorf("Unable to parse patch '%s': '%w'", splitted[2], err)
	}

	err = v.SetMajorMinorPatch(major, minor, patch)
	if err != nil {
		return err
	}

	return nil
}

func (v SemanticVersion) String() string {
	data, err := v.GetAsString()
	if err != nil {
		return fmt.Sprintf("<Unknown SemanticVersion major='%d', minor='%d', patch='%d'>", v.major, v.minor, v.patch)
	}

	return data
}