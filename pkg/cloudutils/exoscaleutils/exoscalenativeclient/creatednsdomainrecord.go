package exoscalenativeclient

import (
	"context"
	"errors"
	"strings"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/publicips"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateDnsDomainRecord(ctx context.Context, client *v3.Client, domainName string, options *dnsoptions.CreateDnsDomainRecordOptions) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if domainName == "" {
		return tracederrors.TracedErrorEmptyString("domainName")
	}

	recordName, err := options.GetName()
	if err != nil {
		return err
	}

	domainName, recordName, err = SanitizeNames(domainName, recordName)
	if err != nil {
		return err
	}

	recordType, err := options.GetRecordType()
	if err != nil {
		return err
	}

	ttl, err := options.GetTTL()
	if err != nil {
		return err
	}

	ipV4Address, err := options.GetIPv4Address()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create Exoscale DNS record '%s' of type '%s' for domain '%s' started.", recordName, recordType, domainName)

	domainUUID, err := GetDomainUuid(ctx, client, domainName)
	if err != nil {
		return err
	}

	if !options.AllowMultipleEntries {
		existingRecords, err := ListDomainRecords(ctx, client, domainName, &dnsoptions.ListDnsDomainRecordOptions{
			Name: recordName,
			RecordType: recordType,
		})
		if err != nil {
			return err
		}

		if len(existingRecords) > 1 {
			logging.LogInfoByCtxf(ctx, "Multiple DNS records '%s' with type '%s' for domain '%s' found. Going to delete them since only 1 should be created.", recordName, recordType, domainName)

			for _, r := range existingRecords {
				_, err := client.DeleteDNSDomainRecord(ctx, domainUUID, r.ID)
				if err != nil {
					return tracederrors.TracedErrorf("Failed to delete DNS record '%s' for domain '%s': %w", recordName, domainName, err)
				}

				logging.LogChangedByCtxf(ctx, "Deleted DNS record '%s' (UUID='%s') of domain '%s'.", recordName, r.ID, domainName)
			}
		}
		if len(existingRecords) == 1 {
			var equal = true
			if existingRecords[0].Content != ipV4Address {
				equal = false
			}

			if existingRecords[0].Name != recordName {
				equal = false
			}

			if string(existingRecords[0].Type) != recordType {
				equal = false
			}

			if existingRecords[0].Ttl != int64(ttl) {
				equal = false
			}

			if !equal {
				logging.LogInfoByCtxf(ctx, "DNS record '%s' in domain '%s' already exists but with different settings. Going to delete before recreating it.", recordName, domainName)
				_, err := client.DeleteDNSDomainRecord(ctx, domainUUID, existingRecords[0].ID)
				if err != nil {
					return tracederrors.TracedErrorf("Failed to delete DNS record '%s' for domain '%s': %w", existingRecords[0].Name, domainName, err)
				}
			}
		}
	}

	var created = true
	_, err = client.CreateDNSDomainRecord(ctx, domainUUID, v3.CreateDNSDomainRecordRequest{
		Name:    recordName,
		Type:    v3.CreateDNSDomainRecordRequestType(recordType),
		Ttl:     int64(ttl),
		Content: ipV4Address,
	})
	if err != nil {
		if errors.Is(err, v3.ErrBadRequest) && strings.Contains(err.Error(), "already exists") {
			created = false
		} else {
			return tracederrors.TracedErrorf("Failed to create dns domain record '%s' for domain '%s': %w", recordName, domainName, err)
		}
	}

	if created {
		logging.LogChangedByCtxf(ctx, "Created Exoscale DNS record '%s' of type '%s' for domain '%s'.", recordName, recordType, domainName)
	}

	logging.LogInfoByCtxf(ctx, "Create Exoscale DNS record '%s' of type '%s' for domain '%s' finished.", recordName, recordType, domainName)

	return nil
}

// Creates (or updates) the record on the DNS domain with the current public IP address.
// This is usfull to act as dynamic DNS client.
func CreateDnsDomainRecordWithCurrentPublicAddress(ctx context.Context, client *v3.Client, domainName string, recordName string) error {
	if client == nil {
		return tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return tracederrors.TracedErrorEmptyString("domainName")
	}

	if recordName == "" {
		return tracederrors.TracedErrorEmptyString("recordName")
	}

	ipAddress, err := publicips.GetPublicIp(ctx)
	if err != nil {
		return err
	}

	return CreateDnsDomainRecord(ctx, client, domainName, &dnsoptions.CreateDnsDomainRecordOptions{
		RecordType:  "A",
		IPv4Address: ipAddress,
		Name:        recordName,
		TTL:         60,
	})
}
