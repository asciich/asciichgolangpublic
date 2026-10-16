package installkvm_test

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/installkvm"
)

// Ensures KVM is installed and libvirtd is running on localhost.
// Requires root privileges. Since there is no "// Output:" comment this example is compiled but not executed by "go test".
func ExampleStart() {
	ctx := context.Background()

	// Installs missing KVM packages and ensures the libvirtd service is enabled and running.
	mustutils.Must0(installkvm.Start(ctx))
}
