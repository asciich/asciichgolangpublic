package exoscalenativeclientoo

import (
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ExoscaleIAM struct {
	client *ExoscaleClient
}

func (e *ExoscaleIAM) GetClient() (*ExoscaleClient, error) {
	if e.client == nil {
		return nil, tracederrors.TracedError("Client not set")
	}

	return e.client, nil
}

func (e *ExoscaleIAM) GetNativeClient() (*v3.Client, error) {
	client, err := e.GetClient()
	if err != nil {
		return nil, err
	}

	return client.GetNativeClient()
}

func (e *ExoscaleIAM) Users() (*ExoscaleIAMUsers, error) {
	users := &ExoscaleIAMUsers{
		iam: e,
	}

	return users, nil
}
