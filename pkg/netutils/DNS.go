package netutils

import (
	"context"
	"net"
	"sort"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func DnsLookupIpV4(ctx context.Context, fqdn string) (ipV4Addresses []string, err error) {
	if fqdn == "" {
		return nil, tracederrors.TracedErrorEmptyString("fqdn")
	}

	logging.LogInfoByCtxf(ctx, "Going to perform DNS lookup for fqdn='%s'", fqdn)

	ips, err := net.LookupIP(fqdn)
	if err != nil {
		return nil, tracederrors.TracedErrorf("LookupIp failed for hostname '%s': %w", fqdn, err)
	}

	for _, ip := range ips {
		v4Addr := ip.To4()
		if v4Addr != nil {
			ipV4Addresses = append(ipV4Addresses, v4Addr.String())
		}
	}

	sort.Strings(ipV4Addresses)

	if len(ipV4Addresses) <= 0 {
		return nil, tracederrors.TracedErrorf("No IPv4 address for host '%s' found.", fqdn)
	}

	logging.LogInfoByCtxf(ctx, "Resolved '%s' to IPv4 addresses '%v'", fqdn, ipV4Addresses)

	return ipV4Addresses, nil
}

func DnsReverseLookup(ctx context.Context, ipAddress string) (fqdns []string, err error) {
	fqdns, err = net.LookupAddr(ipAddress)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to reverse lookup ipAddress '%s': %w",
			ipAddress,
			err,
		)
	}

	logging.LogInfoByCtxf(ctx, "Resolved IP address '%s' to  '%v'", ipAddress, fqdns)

	return fqdns, nil
}

func MustDnsLookupIpV4(ctx context.Context, fqdn string) (ipV4Addresses []string) {
	ipV4Addresses, err := DnsLookupIpV4(ctx, fqdn)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ipV4Addresses
}

func MustDnsReverseLookup(ctx context.Context, ipAddress string) (fqdn []string) {
	fqdn, err := DnsReverseLookup(ctx, ipAddress)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}
