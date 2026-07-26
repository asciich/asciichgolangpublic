# headscaleutils package

Work with [Headscale](https://github.com/juanfont/headscale) - an open source implementation of the Tailcontrol server for Tailscale.

## Examples

* [Connect two nodes with same user](./Example_connecttwonodeswithsameuser_test.go): Sets up a headscale development server with one user, starts two containers as tailscale clients/nodes, and performs a `tailscale ping` to check connectivity.

## Subpackages

* [commandexecutorheadscale](./commandexecutorheadscale/): Non object-oriented implementation to execute headscale commands.
* [commandexecutorheadscaleoo](./commandexecutorheadscaleoo/): Object-oriented implementation to execute headscale commands.
* [headscalegeneric](./headscalegeneric/): Generic headscale functionality.
* [headscaleinterfaces](./headscaleinterfaces/): Interfaces for headscale implementations.
* [headscalelocaldevserver](./headscalelocaldevserver/): Local development server for headscale testing.
