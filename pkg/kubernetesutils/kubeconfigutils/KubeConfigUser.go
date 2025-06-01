package kubeconfigutils

type KubeConfigUser struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
		Username              string `yaml:"username"`
		Password              string `yaml:"password"`
	} `yaml:"user"`
}
