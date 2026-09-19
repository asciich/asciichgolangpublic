package kubeconfigutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

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

func (k KubeConfigContext) GetClusterName() (clusterName string, err error) {
	clusterName = k.Context.Cluster
	if clusterName == "" {
		return "", tracederrors.TracedError("cluster name not set")
	}

	return clusterName, nil
}

func (k *KubeConfig) GetContextNames() (contextNames []string, err error) {
	for _, entry := range k.Contexts {
		toAdd := entry.Name
		if toAdd == "" {
			return nil, tracederrors.TracedErrorf("Got empty context name toAdd")
		}

		contextNames = append(contextNames, toAdd)
	}

	return contextNames, nil
}

func (k *KubeConfig) GetUserNames() (userNames []string, err error) {
	for _, entry := range k.Users {
		toAdd := entry.Name
		if toAdd == "" {
			return nil, tracederrors.TracedErrorf("Got empty user name toAdd")
		}

		userNames = append(userNames, toAdd)
	}

	return userNames, nil
}
