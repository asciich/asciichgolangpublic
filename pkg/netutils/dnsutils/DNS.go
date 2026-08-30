package dnsutils

import (
	"context"
	"net"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func newResolverWithServerLogging(ctx context.Context) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(dialCtx context.Context, network, address string) (net.Conn, error) {
			logging.LogInfoByCtxf(ctx, "DNS server used: '%s'", address)
			var dialer net.Dialer
			return dialer.DialContext(dialCtx, network, address)
		},
	}
}

func DnsLookupIpV4(ctx context.Context, fqdn string) (ipV4Addresses []string, err error) {
	if fqdn == "" {
		return nil, tracederrors.TracedErrorEmptyString("fqdn")
	}

	logging.LogInfoByCtxf(ctx, "Going to perform DNS lookup for fqdn='%s'", fqdn)

	resolver := newResolverWithServerLogging(ctx)
	ips, err := resolver.LookupIP(ctx, "ip4", fqdn)
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
	resolver := newResolverWithServerLogging(ctx)

	fqdns, err = resolver.LookupAddr(ctx, ipAddress)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to reverse lookup ipAddress '%s': %w",
			ipAddress,
			err,
		)
	}

	logging.LogInfoByCtxf(ctx, "Resolved IP address '%s' to '%v'", ipAddress, fqdns)

	return fqdns, nil
}
