package ansiblegalaxyutils

import "github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"

type ListInstalledCollectionsOptions struct {
	// Local path to the root directory of the python virtualenv containing ansible.
	AnsibleVirtualenvPath string
}

func (l *ListInstalledCollectionsOptions) GetAnsiblePath() (string, error) {
	return ansibleparemeteroptions.GetAnsiblePath(l)
}

func (l *ListInstalledCollectionsOptions) GetAnsibleGalaxyPath() (string, error) {
	return ansibleparemeteroptions.GetAnsibleGalaxyPath(l)
}
