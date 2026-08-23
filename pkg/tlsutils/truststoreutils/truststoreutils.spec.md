# truststoreutils specifications

This are the specifications for the [`truststoreutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Overview

`truststoreutils` is a library providing utilities for managing trust stores (system certificate stores). It offers multiple implementation strategies to install and uninstall CA certificates into the operating system's native trust store.

---

## Architecture


### Interfaces

- The `truststoreinterfaces` package contains all shared interfaces and contracts.
- All implementations must satisfy the interfaces defined in `truststoreinterfaces`.

### Implementations

| Implementation | Description |
|---|---|
| `nativetruststoreoo` | Object-oriented native trust store implementation using [smallstep/truststore](https://github.com/smallstep/truststore). Interacts directly with the OS trust store via Go bindings. |
| `commandexecutortruststoreoo` | Object-oriented trust store implementation that manages certificates by executing shell commands (e.g., `update-ca-certificates`, `security add-trusted-cert`, etc.). |

---

## Interface Definitions (Reference)

### TrustStore Interface

The minimum required functions for any trust store implementation:

```go
package truststoreinterfaces

import "context"

// CaCertificate represents a CA certificate in the trust store.
type CaCertificate struct {
	CommonName   string
	Serial       string
	Issuer       string
	NotBefore    time.Time
	NotAfter     time.Time
	Fingerprint  string
	PemEncoded   string
}

// TrustStore defines the contract for any trust store implementation.
type TrustStore interface {
	// AddCaCertificateFromFile installs a CA certificate into the trust store from a file path.
	AddCaCertificateFromFile(ctx context.Context, path string) error

	// AddCaCertificateFromString installs a CA certificate into the trust store from a PEM-encoded string.
	AddCaCertificateFromString(ctx context.Context, caCertPEM string) error

	// DeleteCaCertificateBySerial removes a CA certificate from the trust store identified by its serial number.
	DeleteCaCertificateBySerial(ctx context.Context, serialNumber string) error

	// DeleteCaCertificatesByCommonName removes all CA certificates from the trust store matching the given common name.
	DeleteCaCertificatesByCommonName(ctx context.Context, commonName string) error

	// ListCaCertificates returns all CA certificates currently in the trust store.
	ListCaCertificates(ctx context.Context) ([]CaCertificate, error)
}
```

### Minimum Required Functions

| Function | Description |
|---|---|
| `AddCaCertificateFromFile` | Add a CA certificate to the trust store from a file path. |
| `AddCaCertificateFromString` | Add a CA certificate to the trust store from a PEM-encoded string. |
| `DeleteCaCertificateBySerial` | Delete a specific CA certificate identified by its unique serial number. |
| `DeleteCaCertificatesByCommonName` | Delete all CA certificates matching a given Common Name (CN). Useful for bulk removal when multiple certs share the same CN. |
| `ListCaCertificates` | List all CA certificates currently installed in the trust store. Returns structured certificate metadata. |

---

## Constructor Summary

| Implementation | Constructor Signature | Notes |
|---|---|---|
| `nativetruststoreoo` | `NewNativeTrustStore(trustStorePath string) (*NativeTrustStore, error)` | Pass `""` for default OS trust store, or a custom path. |
| `commandexecutortruststoreoo` | `NewCommandExecutorTrustStore(commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) (*CommandExecutorTrustStore, error)` | Accepts a `CommandExecutor` for flexible command routing. |

---

## Usage Examples

### Example 1: Native Trust Store (Localhost / Default OS Trust Store)

This example uses `nativetruststoreoo` to interact directly with the local system's trust store.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.asciich.ch/tools/truststoreutils.git/pkg/nativetruststoreoo"
)

func main() {
	ctx := context.Background()

	// Create the native trust store instance.
	// Pass "" for the default OS trust store path, or a custom path for a specific trust store location.
	trustStore, err := nativetruststoreoo.NewNativeTrustStore("")
	if err != nil {
		log.Fatalf("failed to create native trust store: %v", err)
	}

	// --- Add CA certificate from a file ---
	err = trustStore.AddCaCertificateFromFile(ctx, "/path/to/my-ca.crt")
	if err != nil {
		log.Fatalf("failed to install CA from file: %v", err)
	}
	fmt.Println("CA certificate installed from file")

	// --- Add CA certificate from a PEM string ---
	caPEM := `-----BEGIN CERTIFICATE-----
MIIBojCCAUmgAwIBAgIRAIx2MnIl...
-----END CERTIFICATE-----`

	err = trustStore.AddCaCertificateFromString(ctx, caPEM)
	if err != nil {
		log.Fatalf("failed to install CA from string: %v", err)
	}
	fmt.Println("CA certificate installed from string")

	// --- List all CA certificates ---
	certs, err := trustStore.ListCaCertificates(ctx)
	if err != nil {
		log.Fatalf("failed to list CA certificates: %v", err)
	}
	fmt.Printf("Found %d CA certificates in trust store:\n", len(certs))
	for _, cert := range certs {
		fmt.Printf("  - CN=%s Serial=%s NotAfter=%s\n", cert.CommonName, cert.Serial, cert.NotAfter)
	}

	// --- Delete CA certificate by serial number ---
	err = trustStore.DeleteCaCertificateBySerial(ctx, "1A:2B:3C:4D:5E:6F")
	if err != nil {
		log.Fatalf("failed to delete CA by serial: %v", err)
	}
	fmt.Println("CA certificate deleted by serial")

	// --- Delete all CA certificates by Common Name ---
	err = trustStore.DeleteCaCertificatesByCommonName(ctx, "My Test CA")
	if err != nil {
		log.Fatalf("failed to delete CAs by common name: %v", err)
	}
	fmt.Println("All CA certificates with CN='My Test CA' deleted")
}
```

---

### Example 2: Command Executor Trust Store with Docker Container Command Executor

This example uses `commandexecutortruststoreoo` with a Docker container-based `CommandExecutor`
to manage certificates inside a running Docker container.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/containerutils/dockerutils/nativedocker"
	"gitlab.asciich.ch/tools/truststoreutils.git/pkg/commandexecutortruststoreoo"
)

func main() {
	ctx := context.Background()

	// Create a Docker container as command executor.
	// The nativedocker.Container implements the commandexecutorinterfaces.CommandExecutor interface.
	container, err := nativedocker.NewContainer("my-test-container-abc123")
	if err != nil {
		log.Fatalf("failed to create docker container command executor: %v", err)
	}

	// Create the command executor trust store, injecting the Docker container as executor.
	trustStore, err := commandexecutortruststoreoo.NewCommandExecutorTrustStore(
		container, // injected command executor (implements commandexecutorinterfaces.CommandExecutor)
		true,      // useSudo
	)
	if err != nil {
		log.Fatalf("failed to create command executor trust store: %v", err)
	}

	// Add CA certificate from file via the Docker container command executor.
	err = trustStore.AddCaCertificateFromFile(ctx, "/path/to/my-ca.crt")
	if err != nil {
		log.Fatalf("failed to install CA: %v", err)
	}
	fmt.Println("CA certificate installed in Docker container from file")

	// Add CA certificate from PEM string via the Docker container command executor.
	caPEM := `-----BEGIN CERTIFICATE-----
MIIBojCCAUmgAwIBAgIRAIx2MnIl...
-----END CERTIFICATE-----`

	err = trustStore.AddCaCertificateFromString(ctx, caPEM)
	if err != nil {
		log.Fatalf("failed to install CA from string: %v", err)
	}
	fmt.Println("CA certificate installed in Docker container from string")

	// List all CA certificates in the container's trust store.
	certs, err := trustStore.ListCaCertificates(ctx)
	if err != nil {
		log.Fatalf("failed to list CA certificates: %v", err)
	}
	fmt.Printf("Found %d CA certificates in container trust store:\n", len(certs))
	for _, cert := range certs {
		fmt.Printf("  - CN=%s Serial=%s\n", cert.CommonName, cert.Serial)
	}

	// Delete CA certificate by serial number.
	err = trustStore.DeleteCaCertificateBySerial(ctx, "1A:2B:3C:4D:5E:6F")
	if err != nil {
		log.Fatalf("failed to delete CA by serial: %v", err)
	}
	fmt.Println("CA certificate deleted by serial in Docker container")

	// Delete all CA certificates matching a Common Name.
	err = trustStore.DeleteCaCertificatesByCommonName(ctx, "My Test CA")
	if err != nil {
		log.Fatalf("failed to delete CAs by common name: %v", err)
	}
	fmt.Println("All CA certificates with CN='My Test CA' deleted from Docker container")
}
```

---

## Testing

### General Rules

- **All exported functions and methods must have unit tests.**
- **Code coverage** should be maintained at a meaningful level for all packages.

### Environment Constraints

| Environment | Allowed? | Notes |
|---|---|---|
| Localhost (host machine) default trust store | **No** | Never test against the host's local trust store directly. This avoids unintended side effects on the developer's or CI runner's system. |
| Docker container | **Yes** | All trust store integration tests must run inside a Docker container to provide an isolated, reproducible environment. |

### Docker-based Integration Tests

Inside the Docker container, the following scenarios **must** be tested:

| Scenario | Description |
|---|---|
| **With `sudo` (running as root)** | Test installing/uninstalling certificates into the container's default trust store when the process has root privileges. |
| **Without `sudo` (running as non-root user)** | Test behavior and expected error handling when the process does **not** have root privileges. Verify that appropriate errors or permission-denied responses are returned. |

### Test Matrix

| Implementation | Unit Tests | Integration (root) | Integration (non-root) |
|---|---|---|---|
| `nativetruststoreoo` | Yes | Yes | Yes |
| `commandexecutortruststoreoo` | Yes | Yes | Yes |

---

## Guidelines

- **Never modify the host system's trust store** during any test run.
- All integration tests must be containerized and idempotent.
- Tests should use self-signed test CA certificates generated specifically for testing purposes.
- Tests must clean up any certificates they install (even on failure), ensuring no state leaks between test runs.
- Use dependency injection (custom `CommandExecutor`) to enable testing against different targets (Docker, SSH, local) without changing trust store logic.
- All methods accept `context.Context` as the first parameter to support cancellation and timeouts.
