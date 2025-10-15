package ansibleutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type RunOptions struct {
	// Path of the playbook to execute
	PlaybookPath string

	// The Limit to a host or group used by ansible
	Limit string

	// The tags to run. If empty the --tags argument is omitted.
	Tags []string

	// Local path to the root directory of the python virtualenv containing ansible.
	AnsibleVirtualenvPath string

	// Local path to the ansible inventory. If not set the Limit will be used as adhoc inventoroy
	AnsibleInventoryPath string

	// Do not delete the temprary playbook (if one is in use) to allow easier debugging.
	KeepTemporaryPlaybook bool
}

func (r *RunOptions) DeepCopy() *RunOptions {
	if r == nil {
		return nil
	}

	ret := new(RunOptions)

	*ret = *r

	return ret
}

func (r *RunOptions) GetPlaybookPath() (string, error) {
	if r.PlaybookPath == "" {
		return "", tracederrors.TracedError("PlaybookPath not set")
	}

	return r.PlaybookPath, nil
}

func (r *RunOptions) GetLimit() (string, error) {
	if r.Limit == "" {
		return "", tracederrors.TracedError("Limit not set")
	}

	return r.Limit, nil
}
