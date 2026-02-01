package exoscalenativeclient

import (
	"context"
	"errors"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetDomainUuid(ctx context.Context, client *v3.Client, domainName string) (v3.UUID, error) {
	if client == nil {
		return "", tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return "", tracederrors.TracedErrorEmptyString("domainName")
	}

	domains, err := client.ListDNSDomains(ctx)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to list exoscale domains: %w", err)
	}

	domain, err := domains.FindDNSDomain(domainName)
	if err != nil {
		if errors.Is(err, v3.ErrNotFound) {
			return "", tracederrors.TracedErrorf("Failed to get uuid of exoscale dns domain '%s'. %w: %w", domainName, dnsgeneric.ErrDnsDomainNotFound, err)
		}
		return "", tracederrors.TracedErrorf("Failed to get domain '%s': %w", domainName, err)
	}

	uuid := domain.ID

	logging.LogInfoByCtxf(ctx, "Exoscale DNS domain '%s' has UUID='%s'", domainName, uuid)

	return uuid, nil
}
