package ansiblegalaxyutils

import "github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"

type InstallCollectionOptions struct {
	// Local path to the root directory of the ansible collection to install.
	LocalCollectionPath string

	// Local path to the root directory of the python virtualenv containing ansible.
	AnsibleVirtualenvPath string
}

func (i *InstallCollectionOptions) GetAnsiblePath() (string, error) {
	return ansibleparemeteroptions.GetAnsiblePath(i)
}

func (i *InstallCollectionOptions) GetAnsibleGalaxyPath() (string, error) {
	return ansibleparemeteroptions.GetAnsibleGalaxyPath(i)
}
