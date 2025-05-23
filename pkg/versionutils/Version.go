package versionutils

type Version interface {
	Equals(other Version) (isEqual bool)
	IsSemanticVersion() (isSemanticVersion bool)
	IsNewerThan(other Version) (isNewerThan bool, err error)
	GetAsString() (version string, err error)
	GetNextVersion(versionType string) (version Version, err error)
	MustGetAsString() (version string)
	MustGetNextVersion(versionType string) (version Version)
}
