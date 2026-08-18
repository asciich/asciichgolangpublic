package commandexecutorx509utils_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/commandexecutorx509utils"
)

func getCommandExecutorImplementationByName(t *testing.T, name string) commandexecutorinterfaces.CommandExecutor {
	t.Helper()

	if name == "bash" {
		return commandexecutorbashoo.Bash()
	}

	if name == "exec" {
		return commandexecutorexecoo.Exec()
	}

	t.Fatalf("Unknown implementation name '%s'", name)
	return nil
}

func generateSelfSignedCert(t *testing.T, dir string) string {
	t.Helper()

	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", "/C=CH/L=Zurich/O=TestOrg/CN=TestCert",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	return certPath
}

func Test_ReadCertificateFromFile(t *testing.T) {
	implementations := []string{"exec", "bash"}

	for _, implName := range implementations {
		implName := implName

		t.Run(implName+"_nil path", func(t *testing.T) {
			executor := getCommandExecutorImplementationByName(t, implName)

			cert, err := commandexecutorx509utils.ReadCertificateFromFile(contextutils.ContextVerbose(), executor, "")
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(implName+"_non-existent file", func(t *testing.T) {
			executor := getCommandExecutorImplementationByName(t, implName)

			cert, err := commandexecutorx509utils.ReadCertificateFromFile(contextutils.ContextVerbose(), executor, "/non/existent/path.pem")
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(implName+"_valid certificate file", func(t *testing.T) {
			executor := getCommandExecutorImplementationByName(t, implName)
			tmpDir := t.TempDir()
			certPath := generateSelfSignedCert(t, tmpDir)

			cert, err := commandexecutorx509utils.ReadCertificateFromFile(contextutils.ContextVerbose(), executor, certPath)
			require.NoError(t, err)
			require.NotNil(t, cert)
			require.Equal(t, "TestOrg", cert.Subject.Organization[0])
			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Zurich", cert.Subject.Locality[0])
			require.Equal(t, "TestCert", cert.Subject.CommonName)
		})

		t.Run(implName+"_invalid certificate file", func(t *testing.T) {
			executor := getCommandExecutorImplementationByName(t, implName)
			tmpDir := t.TempDir()
			invalidPath := filepath.Join(tmpDir, "invalid.pem")
			err := os.WriteFile(invalidPath, []byte("not a valid certificate"), 0644)
			require.NoError(t, err)

			cert, err := commandexecutorx509utils.ReadCertificateFromFile(contextutils.ContextVerbose(), executor, invalidPath)
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(implName+"_certificate is parseable and consistent", func(t *testing.T) {
			executor := getCommandExecutorImplementationByName(t, implName)
			tmpDir := t.TempDir()
			certPath := generateSelfSignedCert(t, tmpDir)

			ctx := contextutils.ContextVerbose()

			cert1, err := commandexecutorx509utils.ReadCertificateFromFile(ctx, executor, certPath)
			require.NoError(t, err)

			cert2, err := commandexecutorx509utils.ReadCertificateFromFile(ctx, executor, certPath)
			require.NoError(t, err)

			require.True(t, cert1.Equal(cert2), "Reading the same certificate file twice should yield equal certificates")
		})
	}
}
