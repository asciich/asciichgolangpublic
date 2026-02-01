package exoscalenativeclientoo

import (
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ExoscaleClient struct {
	client *v3.Client
}

func (e *ExoscaleClient) GetNativeClient() (*v3.Client, error) {
	if e.client == nil {
		return nil, tracederrors.TracedError("native client not set")
	}

	return e.client, nil
}

func (e *ExoscaleClient) DNS() (*ExoscaleDNS, error) {
	dns := &ExoscaleDNS{
		client: e,
	}

	return dns, nil
}

func (e *ExoscaleClient) IAM() (*ExoscaleIAM, error) {
	iam := &ExoscaleIAM{
		client: e,
	}

	return iam, nil
}
