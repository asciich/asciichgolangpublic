package kubeconfigutils

import "github.com/asciich/asciichgolangpublic/tracederrors"

type KubeConfigContext struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		User      string `yaml:"user"`
	} `yaml:"context"`
}

func (k KubeConfigContext) GetUserName() (userName string, err error) {
	userName = k.Context.User
	if userName == "" {
		return "", tracederrors.TracedError("user name not set")
	}

	return userName, nil
}
