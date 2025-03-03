package ansibleutils

import "github.com/asciich/asciichgolangpublic/tracederrors"

type AnsibleHost struct {
	hostName string
}

func NewAnsibleHost() (a *AnsibleHost) {
	return new(AnsibleHost)
}

func NewAnsibleHostByName(hostName string) (a *AnsibleHost, err error) {
	if hostName == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostName")
	}

	a = NewAnsibleHost()

	err = a.SetHostName(hostName)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AnsibleHost) GetHostName() (hostName string, err error) {
	if a.hostName == "" {
		return "", tracederrors.TracedError("hostName not set")
	}

	return a.hostName, nil
}

func (a *AnsibleHost) SetHostName(hostName string) (err error) {
	if hostName == "" {
		return tracederrors.TracedErrorEmptyString("hostName")
	}

	a.hostName = hostName

	return nil
}
