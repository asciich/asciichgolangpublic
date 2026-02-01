package exoscalenativeclient

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListDnsDomainNames(ctx context.Context, client *v3.Client) ([]string, error) {
	logging.LogInfoByCtxf(ctx, "List Exoscale DNS domain names started.")

	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	domains, err := client.ListDNSDomains(ctx)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list exoscale DNS domains: %w", err)
	}

	names := []string{}
	for _, d := range domains.DNSDomains {
		names = append(names, d.UnicodeName)
	}

	logging.LogInfoByCtxf(ctx, "List Exoscale DNS domain names finished. Found %d domains.", len(names))

	return names, nil
}
