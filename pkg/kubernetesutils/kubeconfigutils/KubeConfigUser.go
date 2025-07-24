package kubeconfigutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type KubeConfigUser struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
		Username              string `yaml:"username"`
		Password              string `yaml:"password"`
	} `yaml:"user"`
}

func (k *KubeConfigUser) GetClientKeyData() (string, error) {
	if k.User.ClientKeyData == "" {
		return "", tracederrors.TracedError("ClientKeyData not set")
	}

	return k.User.ClientKeyData, nil
}
