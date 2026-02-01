package exoscalenativeclientoo

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ExoscaleIAMUsers struct {
	iam *ExoscaleIAM
}

func (e *ExoscaleIAMUsers) GetIam() (*ExoscaleIAM, error) {
	if e.iam == nil {
		return nil, tracederrors.TracedError("iam not set")
	}

	return e.iam, nil
}

func (e *ExoscaleIAMUsers) GetNativeClient() (*v3.Client, error) {
	iam, err := e.GetIam()
	if err != nil {
		return nil, err
	}

	return iam.GetNativeClient()
}

func (e *ExoscaleIAMUsers) ListUserNames(ctx context.Context) ([]string, error) {
	nativeClient, err := e.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return exoscalenativeclient.ListIamUsers(ctx, nativeClient)
}
