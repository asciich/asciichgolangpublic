package ansiblegalaxyutils

type ListInstalledCollectionsOptions struct {
	// Local path to the root directory of the python virtualenv containing ansible.
	AnsibleVirtualenvPath string
}

func (l *ListInstalledCollectionsOptions) GetAnsiblePath() (string, error) {
	return GetAnsiblePath(l)
}

func (l *ListInstalledCollectionsOptions) GetAnsibleGalaxyPath() (string, error) {
	return GetAnsibleGalaxyPath(l)
}
