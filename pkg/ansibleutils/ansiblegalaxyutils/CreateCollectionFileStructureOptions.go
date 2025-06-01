package ansiblegalaxyutils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Options to create a collection file structure
//
// Bases on:
//
//	https://docs.ansible.com/ansible/latest/dev_guide/collections_galaxy_meta.html
type CreateCollectionFileStructureOptions struct {
	// The namespace of the collection.
	//
	// This can be a company/brand/organization or product namespace under which all content lives.
	//
	// May only contain alphanumeric lowercase characters and underscores. Namespaces cannot start with underscores or numbers and cannot contain consecutive underscores.
	Namespace string

	// The name of the collection.
	//
	// Has the same character restrictions as namespace.
	Name string

	// The version of the collection.
	//
	// Must be compatible with semantic versioning.
	Version string

	// A list of the collection’s content authors.
	//
	// Can be just the name or in the format ‘Full Name <email> (url) @nicks:irc/im.site#channel’.
	Authors []string
}

func (c *CreateCollectionFileStructureOptions) GetNamespace() (string, error) {
	if c.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	err := CheckValidCollectionName(c.Namespace)
	if err != nil {
		return "", err
	}

	return c.Namespace, nil
}

func (c *CreateCollectionFileStructureOptions) GetName() (string, error) {
	if c.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	err := CheckValidCollectionName(c.Name)
	if err != nil {
		return "", err
	}

	return c.Name, nil
}

func (c *CreateCollectionFileStructureOptions) GetVersion() (versionutils.Version, error) {
	if c.Version == "" {
		return nil, tracederrors.TracedError("Version not set")
	}

	return versionutils.ReadFromString(c.Version)
}

func (c *CreateCollectionFileStructureOptions) GetVersionAsString() (string, error) {
	version, err := c.GetVersion()
	if err != nil {
		return "", err
	}

	versionString, err := version.GetAsString()

	return strings.TrimPrefix(versionString, "v"), nil
}

func (c *CreateCollectionFileStructureOptions) GetAuthors() ([]string, error) {
	if len(c.Authors) <= 0 {
		return nil, tracederrors.TracedError("Authors not set")
	}

	return slicesutils.GetDeepCopyOfStringsSlice(c.Authors), nil
}
