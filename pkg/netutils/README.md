# netutils

Network utilities for various network-related tasks.

## Features

* **TCP Port Utilities** (`TcpPorts.go`):
  - `IsTcpPortOpen()`: Check if a TCP port is open on a remote host
  - `WaitTcpPortOpen()`: Wait for a TCP port to become open
  - `IsTcpPortAvailableForListening()`: Check if a local port is available
  - `WaitPortAvailableForListening()`: Wait for a port to become available
  - `GetNextFreePort(ctx, startPort)`: Find the next free TCP port starting from a given port number

## Subpackages

* [dnsutils](./dnsutils/): Work with DNS.
* [macaddresses](./macaddresses/): Handle MAC addresses.
* [publicips](./publicips/): Get public IP addresses.
* [vpnutils](./vpnutils/): Work with VPN solutions.
    * [headscaleutils](./vpnutils/headscaleutils/): Work with Headscale.
    * [tailscaleutils](./vpnutils/tailscaleutils/): Work with Tailscale.
