package truststoreutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/commandexecutortruststoreoo"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/truststoreinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

func getTestContainer(ctx context.Context, t *testing.T, containerName string, implementationName string, imageName string) (truststoreinterfaces.TrustStore, containerinterfaces.Container) {
	t.Helper()
	require.NotEmpty(t, containerName)
	require.NotEmpty(t, implementationName)
	require.NotEmpty(t, imageName)

	if implementationName != "commandexecutortruststoreoo" {
		t.Fatalf("Unknown implementationName = %s", implementationName)
	}

	container, err := dockerutils.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: imageName,
		Name:      containerName,
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	_, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"bash", "-c", "apt-get update -qq && apt-get install -y -qq ca-certificates"},
	})
	require.NoError(t, err)

	trustStore, err := commandexecutortruststoreoo.NewCommandExecutorTrustStore(container, false)
	require.NoError(t, err)

	return trustStore, container
}

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_AddAndRemoveCert(t *testing.T) {
	var tests = []struct {
		implementationName string
		imageName          string
	}{
		{"commandexecutortruststoreoo", "ubuntu:latest"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				truststore, container := getTestContainer(ctx, t, "test-add-and-remove-cert", tt.implementationName, tt.imageName)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				commonName := "Integration Test CA"

				// 1. Generate a temporary certificate to test with
				options := &x509options.X509CreateCertificateOptions{
					CommonName:     commonName,
					Organization:   "Test Org",
					PrivateKeySize: 2048,
				}
				certKeyPair, err := x509utils.CreateRootCaCertificate(ctx, options)
				require.NoError(t, err)

				certPEMString, err := certKeyPair.GetCertificateAsPEMString()
				require.NoError(t, err)

				// 2. Get baseline of installed certificates
				initialCerts, err := truststore.ListCaCertificates(ctx)
				require.NoError(t, err)
				initialCount := len(initialCerts)

				// 3. Add the generated certificate
				err = truststore.AddCaCertificateFromString(ctx, certPEMString)
				require.NoError(t, err)

				// 4. Verify the certificate was added
				certsAfterAdd, err := truststore.ListCaCertificates(ctx)
				require.NoError(t, err)
				require.Len(t, certsAfterAdd, initialCount+1, "Expected one more certificate in the trust store")

				var found bool
				for _, cert := range certsAfterAdd {
					if cert.Subject.CommonName == commonName {
						found = true
						break
					}
				}
				require.True(t, found, "Added certificate was not found in the trust store list by CommonName")

				// 5. Remove the certificate
				err = truststore.DeleteCaCertificatesByCommonName(ctx, commonName)
				require.NoError(t, err)

				// 6. Verify the certificate was removed
				certsAfterRemove, err := truststore.ListCaCertificates(ctx)
				require.NoError(t, err)
				require.Len(t, certsAfterRemove, initialCount, "Expected certificate count to return to initial baseline")

				found = false
				for _, cert := range certsAfterRemove {
					if cert.Subject.CommonName == commonName {
						found = true
						break
					}
				}
				require.False(t, found, "Certificate should no longer be present in the trust store")
			},
		)
	}
}
