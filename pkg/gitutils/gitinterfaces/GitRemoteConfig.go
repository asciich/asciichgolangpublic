package gitinterfaces

type GitRemoteConfig interface {
	GetRemoteName() (string, error)
	GetUrlFetch() (string, error)
	GetUrlPush() (string, error)
	SetUrlFetch(url string) error
	SetUrlPush(url string) error
	Equals(other GitRemoteConfig) bool
}
