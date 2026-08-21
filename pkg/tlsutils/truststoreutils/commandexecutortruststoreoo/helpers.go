package commandexecutortruststoreoo

import (
	"fmt"
	"runtime"
	"time"
)

func (c *CommandExecutorTrustStore) getInstallCommand(certPath string) (string, error) {
	sudoPrefix := ""
	if c.useSudo {
		sudoPrefix = "sudo "
	}

	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%smkdir -p /usr/local/share/ca-certificates/ && %scp %s /usr/local/share/ca-certificates/ && %supdate-ca-certificates", sudoPrefix, sudoPrefix, certPath, sudoPrefix), nil
	case "darwin":
		return fmt.Sprintf("%ssecurity add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", sudoPrefix, certPath), nil
	case "windows":
		return fmt.Sprintf("certutil -addstore Root %s", certPath), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (c *CommandExecutorTrustStore) getUninstallCommand(certPath string) (string, error) {
	sudoPrefix := ""
	if c.useSudo {
		sudoPrefix = "sudo "
	}

	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%srm -f %s && %supdate-ca-certificates --fresh", sudoPrefix, certPath, sudoPrefix), nil
	case "darwin":
		return fmt.Sprintf("%ssecurity delete-certificate -k /Library/Keychains/System.keychain", sudoPrefix), nil
	case "windows":
		return "certutil -delstore Root", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (c *CommandExecutorTrustStore) getTempCertPath() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return "/tmp/truststore-cert-" + time.Now().Format("20060102150405") + ".crt"
	case "windows":
		return "%TEMP%\\truststore-cert-" + time.Now().Format("20060102150405") + ".crt"
	default:
		return "/tmp/truststore-cert.crt"
	}
}

func (c *CommandExecutorTrustStore) getCertPaths() []string {
	switch runtime.GOOS {
	case "linux":
		return []string{
			"/etc/ssl/certs",
			"/etc/pki/tls/certs",
			"/usr/local/share/ca-certificates",
			"/etc/ca-certificates/trust-source/anchors",
		}
	case "darwin":
		return []string{
			"/etc/ssl/certs",
			"/System/Library/Keychains",
			"/Library/Keychains",
		}
	case "windows":
		return []string{}
	default:
		return []string{}
	}
}
