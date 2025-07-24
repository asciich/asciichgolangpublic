package versionutils

import (
	"fmt"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadDateVersionFromString(versionString string) (*DateVersion, error) {
	version := &DateVersion{}
	err := version.SetVersion(versionString)
	if err != nil {
		return nil, err
	}

	return version, nil
}

type DateVersion struct {
	version string
}

// Creates a new DateVersion set to the current time and date.
func NewCurrentDateVersion() (version Version) {
	version, err := ReadFromString(GetNewDateVersionString())
	if err != nil {
		panic(fmt.Sprintf("internal error: %v", err))
	}

	return version
}
func (v *DateVersion) GetNextVersion(versionType string) (nextVersion Version, err error) {
	return &DateVersion{}, nil
}

func (v *DateVersion) GetVersion() (version string, err error) {
	if v.version == "" {
		return "", tracederrors.TracedErrorf("version not set")
	}

	return v.version, nil
}

func (v *DateVersion) SetVersion(version string) (err error) {
	if version == "" {
		return tracederrors.TracedErrorf("version is empty string")
	}

	v.version = version

	return nil
}

func (v DateVersion) Equals(other Version) (isEqual bool) {
	if other == nil {
		return false
	}

	otherDateVersionPtr, ok := other.(*DateVersion)
	if ok {
		return v.version == otherDateVersionPtr.version
	}

	otherDateVersion, ok := other.(*DateVersion)
	if ok {
		return v.version == otherDateVersion.version
	}

	return false
}

func (v DateVersion) GetAsString() (version string, err error) {
	return v.GetVersion()
}

func (v DateVersion) IsNewerThan(other Version) (isNewerThan bool, err error) {
	if other == nil {
		return false, tracederrors.TracedErrorNil("other")
	}

	otherDateVersion, ok := other.(*DateVersion)
	if !ok {
		return false, tracederrors.TracedErrorf(
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

func (v DateVersion) IsSemanticVersion() (isSemanticVersion bool) {
	return false
}

func (v DateVersion) String() string {
	data, err := v.GetAsString()
	if err != nil {
		return fmt.Sprintf("<Unknown DateVersion '%s'>", v.version)
	}

	return data
}
