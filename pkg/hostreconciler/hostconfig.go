package hostreconciler

// HostConfigFileName is the name of the file inside a config directory
// which contains the host configuration.
const HostConfigFileName = "host.yaml"

// hostConfig represents the content of the host.yaml file.
type hostConfig struct {
	Hostname string `yaml:"hostname"`
}
