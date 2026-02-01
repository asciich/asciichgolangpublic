package exoscalenativeclient

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Deletes all dns domain records matching the set options.
func DeleteDNSDomainRecords(ctx context.Context, client *v3.Client, domainName string, options *dnsoptions.DeleteDnsDomainRecordOptions) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return tracederrors.TracedErrorEmptyString("domainName")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	err := options.Validate()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete Exoscale DNS domain records of domain '%s' started.", domainName)

	records, err := ListDomainRecords(ctx, client, domainName, &dnsoptions.ListDnsDomainRecordOptions{
		Name:       options.Name,
		RecordType: options.RecordType,
		Content:    options.Content,
		TTL:        options.TTL,
	})
	if err != nil {
		return err
	}

	if len(records) <= 0 {
		logging.LogInfoByCtxf(ctx, "No Exoscale DNS domain records of domain '%s' found to delete. Skip delete.", domainName)
	} else {
		logging.LogInfoByCtxf(ctx, "Going to delete %d records in domain '%s'.", len(records), domainName)

		domainUuid, err := GetDomainUuid(ctx, client, domainName)
		if err != nil {
			return err
		}

		for _, r := range records {
			recordUUID := r.ID
			_, err := client.DeleteDNSDomainRecord(ctx, domainUuid, recordUUID)
			if err != nil {
				return tracederrors.TracedErrorf("Failed to delete Exoscale domain record '%s' UUID='%s' of domain '%s' (UUID='%s')", r.Name, recordUUID, domainName, domainUuid)
			}

			logging.LogChangedByCtxf(ctx, "Deleted Exoscale DNS domain record '%s' (UUID='%s') of domain '%s' (UUID='%s').", r.Name, r.ID, domainName, domainUuid)
		}
	}
	logging.LogInfoByCtxf(ctx, "Delete Exoscale DNS domain records of domain '%s' finished.", domainName)

	return nil
}
