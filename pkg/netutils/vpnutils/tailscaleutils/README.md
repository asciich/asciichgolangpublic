# tailscaleutils package

Work with [Tailscale](https://tailscale.com/) - a zero config SDN based on WireGuard.

## Examples

* [Serve and request HTTP over Tailscale](./Example_ServeAndRequestHttpOverTailscale_test.go): Shows how to connect two tailscale nodes where one serves as HTTP server and the other acts as HTTP client. Uses a local headscale instance for testing independence from public control plane servers.

## Subpackages

* [nativetailscaleclient](./nativetailscaleclient/): Non object-oriented Tailscale client implementation.
* [nativetailscaleclientoo](./nativetailscaleclientoo/): Object-oriented Tailscale client implementation.
* [nativetailscalehttpclient](./nativetailscalehttpclient/): HTTP client over Tailscale.
* [nativetailscalehttpserver](./nativetailscalehttpserver/): HTTP server over Tailscale.
* [tailscalegeneric](./tailscalegeneric/): Generic Tailscale functionality.
* [tailscaleoptions](./tailscaleoptions/): Configuration options for Tailscale operations.
