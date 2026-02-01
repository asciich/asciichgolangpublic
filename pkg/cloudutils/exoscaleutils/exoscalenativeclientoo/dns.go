package exoscalenativeclientoo

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ExoscaleDNS struct {
	client *ExoscaleClient
}

func (e *ExoscaleDNS) GetClient() (*ExoscaleClient, error) {
	if e.client == nil {
		return nil, tracederrors.TracedError("client not set")
	}

	return e.client, nil
}

func (e *ExoscaleDNS) GetNativeClient() (*v3.Client, error) {
	client, err := e.GetClient()
	if err != nil {
		return nil, err
	}

	return client.GetNativeClient()
}

func (e *ExoscaleDNS) ListDomainNames(ctx context.Context) ([]string, error) {
	nativeClient, err := e.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return exoscalenativeclient.ListDnsDomainNames(ctx, nativeClient)
}

func (e *ExoscaleDNS) CreateDomainRecord(ctx context.Context, domainName string, options *dnsoptions.CreateDnsDomainRecordOptions) error {
	if domainName == "" {
		return tracederrors.TracedErrorEmptyString("domainName")
	}

	nativeClient, err := e.GetNativeClient()
	if err != nil {
		return err
	}

	return exoscalenativeclient.CreateDnsDomainRecord(ctx, nativeClient, domainName, options)
}
