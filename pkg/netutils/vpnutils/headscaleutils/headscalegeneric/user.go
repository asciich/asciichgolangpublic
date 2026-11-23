package headscalegeneric

import (
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var _ headscaleinterfaces.User = (*User)(nil)

// User is a generic, purely in-memory representation of a headscale user.
type User struct {
	id       uint64
	name     string // Display name, can be empty.
	username string
	email    string // Can be empty.
	created  time.Time
}

func NewUser() *User {
	return &User{}
}

// GetDeepCopy returns a copy of the user which shares no memory with the original.
func (u *User) GetDeepCopy() (*User, error) {
	if u == nil {
		return nil, tracederrors.TracedErrorNil("user")
	}

	ret := *u

	return &ret, nil
}

func (u *User) GetID() (uint64, error) {
	if u == nil {
		return 0, tracederrors.TracedErrorNil("user")
	}

	if u.id == 0 {
		return 0, tracederrors.TracedError("id of the user is not set")
	}

	return u.id, nil
}

func (u *User) SetID(id uint64) error {
	if u == nil {
		return tracederrors.TracedErrorNil("user")
	}

	if id == 0 {
		return tracederrors.TracedError("id must not be 0")
	}

	u.id = id

	return nil
}

// GetName returns the display name. Returns an empty string if not set.
func (u *User) GetName() (string, error) {
	if u == nil {
		return "", tracederrors.TracedErrorNil("user")
	}

	return u.name, nil
}

func (u *User) SetName(name string) error {
	if u == nil {
		return tracederrors.TracedErrorNil("user")
	}

	u.name = name

	return nil
}

func (u *User) GetUsername() (string, error) {
	if u == nil {
		return "", tracederrors.TracedErrorNil("user")
	}

	if u.username == "" {
		return "", tracederrors.TracedError("username of the user is not set")
	}

	return u.username, nil
}

func (u *User) SetUsername(username string) error {
	if u == nil {
		return tracederrors.TracedErrorNil("user")
	}

	if username == "" {
		return tracederrors.TracedErrorEmptyString("username")
	}

	u.username = username

	return nil
}

// GetEmail returns the email address. Returns an empty string if not set.
func (u *User) GetEmail() (string, error) {
	if u == nil {
		return "", tracederrors.TracedErrorNil("user")
	}

	return u.email, nil
}

func (u *User) SetEmail(email string) error {
	if u == nil {
		return tracederrors.TracedErrorNil("user")
	}

	u.email = email

	return nil
}

func (u *User) GetCreated() (time.Time, error) {
	if u == nil {
		return time.Time{}, tracederrors.TracedErrorNil("user")
	}

	if u.created.IsZero() {
		return time.Time{}, tracederrors.TracedError("created time of the user is not set")
	}

	return u.created, nil
}

func (u *User) SetCreated(created time.Time) error {
	if u == nil {
		return tracederrors.TracedErrorNil("user")
	}

	if created.IsZero() {
		return tracederrors.TracedError("created must not be the zero time")
	}

	u.created = created

	return nil
}
