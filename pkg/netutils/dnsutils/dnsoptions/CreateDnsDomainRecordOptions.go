package dnsoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateDnsDomainRecordOptions struct {
	RecordType string

	Name string

	IPv4Address string

	TTL int

	// If set to true multiple entries with the same "name" and "recordType" is allowed.
	// Set this to true if you want multiple entries with different IPs.
	AllowMultipleEntries bool
}

func (c *CreateDnsDomainRecordOptions) GetName() (string, error) {
	if c.Name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.Name, nil
}

func (c *CreateDnsDomainRecordOptions) GetRecordType() (string, error) {
	if c.RecordType == "" {
		return "", tracederrors.TracedError("recordtype not set")
	}

	return c.RecordType, nil
}

func (c *CreateDnsDomainRecordOptions) GetIPv4Address() (string, error) {
	if c.IPv4Address == "" {
		return "", tracederrors.TracedError("IPv4Address not set")
	}

	return c.IPv4Address, nil
}

func (c *CreateDnsDomainRecordOptions) GetTTL() (int, error) {
	if c.TTL <= 0 {
		return 0, tracederrors.TracedError("TTl not set")
	}

	return c.TTL, nil
}
