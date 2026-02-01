package exoscalenativeclient

import (
	"context"
	"errors"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetDnsDomainAndRecordUUID(ctx context.Context, client *v3.Client, domainName string, recordName string) (string, string, error) {
	if client == nil {
		return "", "", tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return "", "", tracederrors.TracedErrorEmptyString("domainName")
	}

	if recordName == "" {
		return "", "", tracederrors.TracedErrorEmptyString("recordName")
	}

	domainName, recordName, err := SanitizeNames(domainName, recordName)
	if err != nil {
		return "", "", err
	}

	domainUuid, err := GetDomainUuid(ctx, client, domainName)
	if err != nil {
		return "", "", err
	}

	records, err := client.ListDNSDomainRecords(ctx, domainUuid)
	if err != nil {
		return "", "", tracederrors.TracedErrorf("Failed to list Exoscale DNS domain records for domainUUID='%s': %w", domainUuid, err)
	}

	record, err := records.FindDNSDomainRecord(recordName)
	if err != nil {
		if errors.Is(err, v3.ErrNotFound) {
			return "", "", tracederrors.TracedErrorf("%w: %w", dnsgeneric.ErrDnsDomainRecordNotFound, err)
		}
		return "", "", tracederrors.TracedErrorf("Failed to find DNS domain record named '%s': %w", recordName, err)
	}

	recordUUID := record.ID

	logging.LogInfoByCtxf(ctx, "Exoscale DNS domain record '%s' on domain '%s' (domainId='%s') has id '%s'", recordName, domainName, domainUuid, recordUUID)

	return string(domainUuid), string(recordUUID), nil
}

func ListDomainRecords(ctx context.Context, client *v3.Client, domainName string, options *dnsoptions.ListDnsDomainRecordOptions) ([]*v3.DNSDomainRecord, error) {
	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return nil, tracederrors.TracedErrorEmptyString("domainName")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtxf(ctx, "List Exoscale DNS domain records of domain '%s' started.", domainName)

	domainUuid, err := GetDomainUuid(ctx, client, domainName)
	if err != nil {
		return nil, err
	}

	records, err := client.ListDNSDomainRecords(ctx, domainUuid)
	if err != nil {
		return nil, err
	}

	ret := []*v3.DNSDomainRecord{}
	for _, r := range records.DNSDomainRecords {
		if options.Name != "" {
			_, recordName, err := SanitizeNames(domainName, options.Name)
			if err != nil {
				return nil, err
			}

			if r.Name != recordName {
				continue
			}
		}

		if options.Content != "" {
			if r.Content != options.Content {
				continue
			}
		}

		if options.RecordType != "" {
			if string(r.Type) != options.RecordType {
				continue
			}
		}

		if options.TTL > 0 {
			if r.Ttl != int64(options.TTL) {
				continue
			}
		}

		ret = append(ret, &r)
	}

	logging.LogInfoByCtxf(ctx, "List Exoscale DNS domain records of domain '%s' finished. Found %d matching domains.", domainName, len(ret))

	return ret, nil
}

func IsDomainRecordEqual(ctx context.Context, client *v3.Client, domainName string, options *dnsoptions.IsEqualOptions) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return false, tracederrors.TracedErrorEmptyString("domainName")
	}

	if options == nil {
		return false, tracederrors.TracedErrorNil("options")
	}

	records, err := ListDomainRecords(ctx, client, domainName, &dnsoptions.ListDnsDomainRecordOptions{
		Name:       options.Name,
		RecordType: options.RecordType,
	})
	if err != nil {
		return false, err
	}

	if len(records) != 1 {
		return false, nil
	}

	if records[0].Content != options.Content {
		return false, nil
	}

	if records[0].Name != options.Name {
		return false, nil
	}

	if string(records[0].Type) != options.RecordType {
		return false, nil
	}

	if records[0].Ttl != int64(options.TTL) {
		return false, nil
	}

	return true, nil
}

func DnsDomainRecordExists(ctx context.Context, client *v3.Client, domainName string, options *dnsoptions.DnsDomainRecordExistsOptions) (bool, error) {
	if client == nil {
		return false, tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return false, tracederrors.TracedErrorEmptyString("domainName")
	}

	if options == nil {
		return false, tracederrors.TracedErrorNil("options")
	}

	domainUuid, err := GetDomainUuid(ctx, client, domainName)
	if err != nil {
		return false, err
	}

	records, err := client.ListDNSDomainRecords(ctx, domainUuid)
	if err != nil {
		return false, err
	}

	var exists bool
	var recordName = options.Name
	for _, r := range records.DNSDomainRecords {
		if options.Name != "" {
			_, recordName, err = SanitizeNames(domainName, options.Name)
			if err != nil {
				return false, err
			}
			if r.Name != recordName {
				continue
			}
		}

		if options.RecordType != "" {
			if r.Type != v3.DNSDomainRecordType(options.RecordType) {
				continue
			}
		}

		recordName = r.Name
		exists = true
		break
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Exoscale DNS domain record '%s' of domain '%s' exists.", recordName, domainName)
	} else {
		logging.LogInfoByCtxf(ctx, "Exoscale DNS domain record '%s' of domain '%s' does not exists.", recordName, domainName)
	}

	return exists, nil
}

func GetDnsDomainRecordContent(ctx context.Context, client *v3.Client, domainName string, recordName string) (string, error) {
	if client == nil {
		return "", tracederrors.TracedErrorNil("client")
	}

	if domainName == "" {
		return "", tracederrors.TracedErrorEmptyString("domainName")
	}

	if recordName == "" {
		return "", tracederrors.TracedErrorEmptyString("recordName")
	}

	logging.LogInfoByCtxf(ctx, "Get Exoscale DNS domain '%s' record '%s' content started.", domainName, recordName)

	domainUUID, recordUUID, err := GetDnsDomainAndRecordUUID(ctx, client, domainName, recordName)
	if err != nil {
		return "", err
	}

	record, err := client.GetDNSDomainRecord(ctx, v3.UUID(domainUUID), v3.UUID(recordUUID))
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get DNS domain '%s' reacord '%s': %w", domainName, recordName, err)
	}

	var content = record.Content

	logging.LogInfoByCtxf(ctx, "Get Exoscale DNS domain '%s' record '%s' content finished. Content is '%s'", domainName, recordName, content)

	return content, nil
}
