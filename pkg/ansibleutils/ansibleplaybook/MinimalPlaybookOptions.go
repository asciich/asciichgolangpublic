package ansibleplaybook

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type MinimalPlaybookOptions struct {
	Hostname   string
	Roles      []string
	RemoteUser string
}

func (m *MinimalPlaybookOptions) GetHostName() (string, error) {
	if m.Hostname == "" {
		return "", tracederrors.TracedError("Hostname not set")
	}

	return m.Hostname, nil
}
