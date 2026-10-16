# hostreconciler

> Bring the GitOps approach to your hosts.

`hostreconciler` periodically reads a declarative configuration and
reconciles the current host to match it: **desired state in, converged host out.**

## Config structure

```text
<config-dir>/
├── host.yaml
└── docker_compose/              (optional)
    └── mydeployment/            ← docker compose project name
        ├── compose.yaml
        └── .env                 (optional)
```

### Minimal config

**`host.yaml`**

```yaml
---
hostname: host.example.com
```

### Full config

**`host.yaml`**

```yaml
---
before_script:
  - echo hello world
hostname: host.example.com
kvm:
  start: true
  enable: true
network:
  interfaces:
    - name: enp0
      layer2_bridge: br0
      deactivate_offloading: true
```

**`docker_compose/mydeployment/compose.yaml`**

```yaml
services:
  hello:
    image: nginx:alpine
```

⚠️ All settings of a network interface must be in the **same** list item (only one `-` per interface).

### Reference: `host.yaml`

| Field                                         | Type       | Required | Default | Description                                                               |
|-----------------------------------------------|------------|:--------:|:-------:|---------------------------------------------------------------------------|
| `before_script`                               | `[]string` |    ❌    |         | Shell commands executed before any change is applied to the host.         |
| `hostname`                                    | `string`   |    ✅    |         | Fully qualified hostname of the host.                                     |
| `kvm`                                         | `object`   |    ❌    |         | KVM configuration of the host.                                            |
| `kvm.start`                                   | `bool`     |    ❌    | `false` | Install the KVM packages and ensure `libvirtd` is running.                |
| `kvm.enable`                                  | `bool`     |    ❌    | `false` | Install the KVM packages and ensure `libvirtd` is started on boot.        |
| `network`                                     | `object`   |    ❌    |         | Network configuration of the host.                                        |
| `network.interfaces`                          | `list`     |    ❌    |         | Desired configuration of network interfaces (network cards).              |
| `network.interfaces[].name`                   | `string`   |    ✅    |         | Name of the network card, e.g. `enp0` (max. 15 characters).               |
| `network.interfaces[].layer2_bridge`          | `string`   |    ❌    |         | Layer 2 bridge to create with this network card attached, e.g. `br0`.    |
| `network.interfaces[].deactivate_offloading`  | `bool`     |    ❌    | `false` | Deactivate TSO/GSO/GRO/LRO offloading for this network card.              |

## Order of a reconcile run

1. Read and validate the whole config (`host.yaml` and `docker_compose/`). If it's invalid, nothing is executed or changed.
2. Check root privileges if `kvm.start` or `kvm.enable` is set. If they are missing, nothing is executed or changed.
3. Run `before_script`. If an entry fails, no change is applied.
4. Apply `hostname`.
5. Reconcile `network.interfaces`.
6. Reconcile `kvm`.
7. Reconcile the `docker_compose/` deployments (one after the other).

### `before_script`

- The entries are executed in the given order on the host.
- Every entry runs in its own `bash -e -o pipefail -c <entry>`:
    - Multi-line entries (YAML block scalar `|`) stop at the first failing command.
    - Shell state (e.g. `cd`, exported variables) is **not** shared between entries.
- Execution stops at the first failing entry and the host is left untouched.
- The entries run with the privileges of the hostreconciler process.

⚠️ `before_script` is executed on **every** reconcile run. Write the commands in an idempotent way (e.g. `mkdir -p` instead of `mkdir`).

Multi-line example:

```yaml
before_script:
  - echo hello world
  - |
    mkdir -p /opt/example
    touch /opt/example/marker
```

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

### `kvm`

Reconciled using [installkvm](../vmutils/kvmutils/installkvm/README.md):

- `kvm.enable`: Uses `Enable` to install missing KVM packages and enable the `libvirtd` systemd service (started on boot). `libvirtd` is **not** started by this option alone.
- `kvm.start`: Uses `Start` to install missing KVM packages and start the `libvirtd` systemd service. `libvirtd` is **not** enabled by this option alone.
- Set both options so `libvirtd` is running now **and** after a reboot.
- If both are set, `libvirtd` is enabled before it is started.
- The installed packages depend on the package manager of the host:
    - `apt-get` (Debian/Ubuntu): `qemu-system-x86`, `qemu-utils`, `libvirt-daemon-system`, `libvirt-clients`, `virtinst`, `ovmf`.
    - `pacman`/`yay` (Arch Linux): `qemu-system-x86`, `qemu-img`, `libvirt`, `virt-install`, `dnsmasq`, `edk2-ovmf`.
- KVM is reconciled **after** `network.interfaces`, so bridges like `br0` used by the virtual machines already exist when `libvirtd` starts.
- Idempotent: Already installed packages and an already running/enabled `libvirtd` are left untouched.

### `docker_compose/<deployment>/`

Every deployment is reconciled using `UpByComposeYamlPath` of [commandexecutordockercomposeutils](../containerutils/dockerutils/dockercomposeutils/commandexecutordockercomposeutils/README.md)
(equivalent to `cd docker_compose/<deployment> && docker compose up --detach --remove-orphans`):

- Every subdirectory of `docker_compose/` is one deployment. The **directory name is the docker compose project name**:
    - Only lowercase letters, digits, `-` and `_`, starting with a lowercase letter or digit.
    - Unique by design, so deployments never treat each other's services as orphans.
- Every deployment directory contains **exactly one** of `compose.yaml`, `compose.yml`, `docker-compose.yaml` or `docker-compose.yml`.
- A top level `name` in the compose file is optional. If set, it must equal the deployment directory name.
- `.env`, bind mounts, `env_file` and `build` contexts are resolved relative to the deployment directory.
- The deployments are applied one after the other in **lexical order** of their directory names. Use prefixes like `10-`, `20-` to control the order.
- `*.yaml` / `*.yml` files directly inside `docker_compose/` are rejected. Hidden entries and other files (e.g. `README.md`) are skipped.
- Services removed from a compose file are removed from the host (`--remove-orphans`).
- Idempotent: Already running services are left untouched.

## Requirements

- Linux.
- `bash` (only if `before_script` is configured).
- Root privileges (only if at least one network interface has `layer2_bridge` or `deactivate_offloading: true`, or if `kvm.start` or `kvm.enable` is `true`).
- `layer2_bridge`: An active and enabled `systemd-networkd`.
- `deactivate_offloading`: systemd/udev (existing `/etc/systemd/network/` directory).
- `kvm`:
    - systemd (to manage the `libvirtd` service).
    - A supported package manager (`apt-get`, `pacman` or `yay`) with access to the package servers.
    - A CPU with hardware virtualization support (Intel VT-x/AMD-V) enabled in the BIOS/UEFI to run virtual machines.
- `docker_compose/`: Docker Compose v2 plugin (`docker compose`) and access to the docker daemon.

## Validation

Performed for the whole config before anything is executed or changed:

- `before_script` contains no empty entries.
- `name` is set and is a valid interface name.
- A network interface is defined only once.
- `layer2_bridge` (if set) is a valid interface name and differs from `name`.
- A layer 2 bridge is used by only one network interface.
- Nested bridges (a configured `layer2_bridge` also listed as `name`) are not supported.
- `kvm.start` and `kvm.enable` (if set) are booleans.
- Root privileges are available if `kvm.start` or `kvm.enable` is `true`.
- Every deployment directory name is a valid docker compose project name.
- Every deployment directory contains exactly one compose file which is valid YAML.
- A top level `name` in a compose file (if set) equals its deployment directory name.
- No compose files are placed directly inside `docker_compose/`.

ℹ️ Reconciliation is **additive**:

- Removing `layer2_bridge` (or the whole entry) does **not** remove the bridge from the host.
- Removing `deactivate_offloading` or setting it to `false` does **not** reactivate offloading.
- Removing `kvm.start` or setting it to `false` does **not** stop `libvirtd`.
- Removing `kvm.enable` or setting it to `false` does **not** disable `libvirtd`.
- Removing the `kvm` section does **not** uninstall the KVM packages.
- Removing a whole deployment directory does **not** stop its project. Run `docker compose -p <deployment> down` manually.
