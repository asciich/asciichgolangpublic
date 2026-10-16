# installkvm specifications

## Implementation

- The `InstallPackages(ctx context.Context)` is used to install all needed KVM related system packages.
- The `Start(ctx context.Context)` Ensures all needed packages are installed and `libvirtd` is running:

# installkvm specifications

## Implementation

- The `InstallPackages(ctx context.Context)` is used to install all needed KVM related system packages.
- The `Start(ctx context.Context)` Ensures all needed packages are installed and `libvirtd` is running.
    - This does NOT create any libvirt config like creating the default pool.
- The `Enable(ctx context.Context)` Ensures all needed packages are installed and `libvirtd` is enabled (started on boot).
    - This does NOT create any libvirt config like creating the default pool.
- `Start` and `Enable` are independent: `Start` does not enable `libvirtd` and `Enable` does not start it.
