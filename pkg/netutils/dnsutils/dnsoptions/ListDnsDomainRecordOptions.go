package dnsoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ListDnsDomainRecordOptions struct {
	RecordType string

	Name string

	Content string

	TTL int
}

func (c *ListDnsDomainRecordOptions) GetName() (string, error) {
	if c.Name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.Name, nil
}

func (c *ListDnsDomainRecordOptions) GetRecordType() (string, error) {
	if c.RecordType == "" {
		return "", tracederrors.TracedError("recordtype not set")
	}

	return c.RecordType, nil
}

func (c *ListDnsDomainRecordOptions) GetContent() (string, error) {
	if c.Content == "" {
		return "", tracederrors.TracedError("Content not set")
	}

	return c.Content, nil
}

func (c *ListDnsDomainRecordOptions) GetTTL() (int, error) {
	if c.TTL <= 0 {
		return 0, tracederrors.TracedError("TTl not set")
	}

	return c.TTL, nil
}

func (c *ListDnsDomainRecordOptions) GetRecordNameOrEmptyStringIfUnset() string {
	return c.Name
}
