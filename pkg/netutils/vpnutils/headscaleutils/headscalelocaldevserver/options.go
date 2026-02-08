package headscalelocaldevserver

const DEFAULT_PORT = 8088
const DEFAULT_CONTAINER_NAME = "headscale-localdev"

type RunOptions struct {
	Port          int
	ContainerName string

	RestartAlreadyRunningDevServer bool
}

func (r *RunOptions) GetPort() int {
	if r.Port <= 0 {
		return DEFAULT_PORT
	}

	return r.Port
}

func (r *RunOptions) GetContainerName() string {
	if r.ContainerName == "" {
		return DEFAULT_CONTAINER_NAME
	}

	return r.ContainerName
}
