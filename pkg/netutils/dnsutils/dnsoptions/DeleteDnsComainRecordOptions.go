package dnsoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DeleteDnsDomainRecordOptions struct {
	RecordType string

	Name string

	Content string

	TTL int
}

func (c *DeleteDnsDomainRecordOptions) GetName() (string, error) {
	if c.Name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.Name, nil
}

func (c *DeleteDnsDomainRecordOptions) GetRecordType() (string, error) {
	if c.RecordType == "" {
		return "", tracederrors.TracedError("recordtype not set")
	}

	return c.RecordType, nil
}

func (c *DeleteDnsDomainRecordOptions) GetContent() (string, error) {
	if c.Content == "" {
		return "", tracederrors.TracedError("Content not set")
	}

	return c.Content, nil
}

func (c *DeleteDnsDomainRecordOptions) GetTTL() (int, error) {
	if c.TTL <= 0 {
		return 0, tracederrors.TracedError("TTl not set")
	}

	return c.TTL, nil
}

// Validate at least one value set.
// Usefull to check if no value was set by mistake which would delete all records of a domain.
func (c *DeleteDnsDomainRecordOptions) Validate() error {
	if c.Content != "" {
		return nil
	}

	if c.Name != "" {
		return nil
	}

	if c.RecordType != "" {
		return nil
	}

	if c.TTL > 0 {
		return nil
	}

	return tracederrors.TracedError("Options validation failed: No value is set.")
}
