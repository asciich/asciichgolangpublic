package installkvm_test

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/installkvm"
)

// Installs all KVM related system packages on localhost.
// Requires root privileges. Since there is no "// Output:" comment this example is compiled but not executed by "go test".
func ExampleInstallPackages() {
	ctx := context.Background()

	// Idempotent: Packages already installed are left untouched.
	mustutils.Must0(installkvm.InstallPackages(ctx))
}
