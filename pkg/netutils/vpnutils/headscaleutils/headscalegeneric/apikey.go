package headscalegeneric

import (
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var _ headscaleinterfaces.ApiKey = (*ApiKey)(nil)

// ApiKey is a generic, purely in-memory representation of a headscale API key.
type ApiKey struct {
	id         uint64
	prefix     string
	expiration *time.Time
	created    time.Time
}

func NewApiKey() *ApiKey {
	return &ApiKey{}
}

// GetDeepCopy returns a copy of the API key which shares no memory with the original.
func (a *ApiKey) GetDeepCopy() (*ApiKey, error) {
	if a == nil {
		return nil, tracederrors.TracedErrorNil("apiKey")
	}

	ret := *a
	ret.expiration = normalizeTime(a.expiration)

	return &ret, nil
}

func (a *ApiKey) GetID() (uint64, error) {
	if a == nil {
		return 0, tracederrors.TracedErrorNil("apiKey")
	}

	if a.id == 0 {
		return 0, tracederrors.TracedError("id of the api key is not set")
	}

	return a.id, nil
}

func (a *ApiKey) SetID(id uint64) error {
	if a == nil {
		return tracederrors.TracedErrorNil("apiKey")
	}

	if id == 0 {
		return tracederrors.TracedError("id must not be 0")
	}

	a.id = id

	return nil
}

func (a *ApiKey) GetPrefix() (string, error) {
	if a == nil {
		return "", tracederrors.TracedErrorNil("apiKey")
	}

	if a.prefix == "" {
		return "", tracederrors.TracedError("prefix of the api key is not set")
	}

	return a.prefix, nil
}

func (a *ApiKey) SetPrefix(prefix string) error {
	if a == nil {
		return tracederrors.TracedErrorNil("apiKey")
	}

	if prefix == "" {
		return tracederrors.TracedErrorEmptyString("prefix")
	}

	a.prefix = prefix

	return nil
}

// GetExpiration returns a copy of the expiration. Returns nil if the key does not expire.
func (a *ApiKey) GetExpiration() (*time.Time, error) {
	if a == nil {
		return nil, tracederrors.TracedErrorNil("apiKey")
	}

	return normalizeTime(a.expiration), nil
}

// SetExpiration sets the expiration. nil and the zero time (0001-01-01 00:00:00) both mean "no expiration".
func (a *ApiKey) SetExpiration(expiration *time.Time) error {
	if a == nil {
		return tracederrors.TracedErrorNil("apiKey")
	}

	a.expiration = normalizeTime(expiration)

	return nil
}

func (a *ApiKey) GetCreated() (time.Time, error) {
	if a == nil {
		return time.Time{}, tracederrors.TracedErrorNil("apiKey")
	}

	if a.created.IsZero() {
		return time.Time{}, tracederrors.TracedError("created time of the api key is not set")
	}

	return a.created, nil
}

func (a *ApiKey) SetCreated(created time.Time) error {
	if a == nil {
		return tracederrors.TracedErrorNil("apiKey")
	}

	if created.IsZero() {
		return tracederrors.TracedError("created must not be the zero time")
	}

	a.created = created

	return nil
}
