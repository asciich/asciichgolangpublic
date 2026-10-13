# hostreconciler

> Bring the GitOps approach to your hosts.

`hostreconciler` periodically reads a declarative configuration and
reconciles the current host to match it: **desired state in, converged host out.**

## Config structure

```text
<config-dir>/
└── host.yaml
```

### Minimal config

**`host.yaml`**

```yaml
---
hostname: host.example.com
```

### Config with network interfaces

**`host.yaml`**

```yaml
---
hostname: host.example.com
network:
  interfaces:
    - name: enp0
      layer2_bridge: br0
      deactivate_offloading: true
```

⚠️ All settings of a network interface must be in the **same** list item (only one `-` per interface).

### Reference: `host.yaml`

| Field                                        | Type     | Required | Default | Description                                                                 |
|----------------------------------------------|----------|:--------:|:-------:|-----------------------------------------------------------------------------|
| `hostname`                                   | `string` |    ✅    |         | Fully qualified hostname of the host.                                       |
| `network`                                    | `object` |    ❌    |         | Network configuration of the host.                                          |
| `network.interfaces`                         | `list`   |    ❌    |         | Desired configuration of network interfaces (network cards).                |
| `network.interfaces[].name`                  | `string` |    ✅    |         | Name of the network card, e.g. `enp0` (max. 15 characters).                 |
| `network.interfaces[].layer2_bridge`         | `string` |    ❌    |         | Layer 2 bridge to create with this network card attached, e.g. `br0`.       |
| `network.interfaces[].deactivate_offloading`  | `bool`   |    ❌    | `false` | Deactivate TSO/GSO/GRO/LRO offloading for this network card.                |

### `network.interfaces[].deactivate_offloading`

Reconciled using `DeactivateOffloading` of [networkcardutils](../netutils/networkcardutils/README.md):

- Deactivates TCP segmentation offload (TSO/TSO6), generic segmentation offload (GSO), generic receive offload (GRO) and large receive offload (LRO).
- Applied instantly (ethtool ioctl) and persistent across reboots (systemd link file `/etc/systemd/network/10-<name>-no-offload.link`).
- Known workaround for network card hangs like the e1000e "Detected Hardware Unit Hang".
- Idempotent: Already deactivated features and an up to date link file are left untouched.

### `network.interfaces[].layer2_bridge`

Reconciled using `CreateLayer2Bridge` of [networkcardutils](../netutils/networkcardutils/README.md):

- The bridge takes over the MAC address of the network card and gets its IP configuration using DHCP.
- Applied instantly and persistent across reboots (systemd-networkd files in `/etc/systemd/network/`).
- Idempotent: An already existing and up to date bridge is left untouched.

If both are set, offloading is deactivated before the bridge is created.

### Requirements

- Linux.
- Root privileges (only if at least one network interface has `layer2_bridge` or `deactivate_offloading: true`).
- `layer2_bridge`: An active and enabled `systemd-networkd`.
- `deactivate_offloading`: systemd/udev (existing `/etc/systemd/network/` directory).

### Validation

Performed for the whole config before any change is applied:

- `name` is set and is a valid interface name.
- A network interface is defined only once.
- `layer2_bridge` (if set) is a valid interface name and differs from `name`.
- A layer 2 bridge is used by only one network interface.
- Nested bridges (a configured `layer2_bridge` also listed as `name`) are not supported.

ℹ️ Reconciliation is **additive**:

- Removing `layer2_bridge` (or the whole entry) does **not** remove the bridge from the host.
- Removing `deactivate_offloading` or setting it to `false` does **not** reactivate offloading.
