package netutils

import (
	"context"
	"net"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func LookupIpV4(ctx context.Context, fqdn string) (ipV4Addresses []string, err error) {
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

	if len(ipV4Addresses) <= 0 {
		return nil, tracederrors.TracedErrorf("No IPv4 address for host '%s' found.", fqdn)
	}

	logging.LogInfoByCtxf(ctx, "Resolved '%s' to IPv4 addresses '%s'", fqdn, ipV4Addresses)

	return ipV4Addresses, nil
}

func MustLookupIpV4(ctx context.Context, fqdn string) (ipV4Addresses []string) {
	ipV4Addresses, err := LookupIpV4(ctx, fqdn)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ipV4Addresses
}
