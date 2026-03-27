package dockeroptions

import (
	"slices"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DockerRunContainerOptions struct {
	ImageName            string
	Name                 string
	Command              []string
	Ports                []string
	Mounts               []string
	KeepStoppedContainer bool
	UseHostNet           bool
	AdditionalEnvVars    map[string]string

	// If Ports are specified this waits until a connect to all ports is accepted:
	WaitForPortsOpen bool
}

func NewDockerRunContainerOptions() (d *DockerRunContainerOptions) {
	return new(DockerRunContainerOptions)
}

func (d *DockerRunContainerOptions) GetDeepCopy() *DockerRunContainerOptions {
	copy := new(DockerRunContainerOptions)

	*copy = *d

	if d.Command != nil {
		copy.Command = slicesutils.GetDeepCopyOfStringsSlice(d.Command)
	}

	if d.Ports != nil {
		copy.Ports = slicesutils.GetDeepCopyOfStringsSlice(d.Ports)
	}

	if d.Mounts != nil {
		copy.Mounts = slicesutils.GetDeepCopyOfStringsSlice(d.Mounts)
	}

	return copy
}

func (d *DockerRunContainerOptions) GetCommand() (command []string, err error) {
	if d.Command == nil {
		return nil, tracederrors.TracedErrorf("Command not set")
	}

	if len(d.Command) <= 0 {
		return nil, tracederrors.TracedErrorf("Command has no elements")
	}

	return d.Command, nil
}

func (d *DockerRunContainerOptions) GetKeepStoppedContainer() (keepStoppedContainer bool, err error) {

	return d.KeepStoppedContainer, nil
}

func (d *DockerRunContainerOptions) GetMounts() (mounts []string, err error) {
	if d.Mounts == nil {
		return nil, tracederrors.TracedErrorf("Mounts not set")
	}

	if len(d.Mounts) <= 0 {
		return nil, tracederrors.TracedErrorf("Mounts has no elements")
	}

	return d.Mounts, nil
}

func (d *DockerRunContainerOptions) GetPorts() (ports []string, err error) {
	if d.Ports == nil {
		return nil, tracederrors.TracedErrorf("Ports not set")
	}

	if len(d.Ports) <= 0 {
		return nil, tracederrors.TracedErrorf("Ports has no elements")
	}

	return d.Ports, nil
}

func (d *DockerRunContainerOptions) GetUseHostNet() (useHostNet bool, err error) {

	return d.UseHostNet, nil
}

func (d *DockerRunContainerOptions) SetCommand(command []string) (err error) {
	if command == nil {
		return tracederrors.TracedErrorf("command is nil")
	}

	if len(command) <= 0 {
		return tracederrors.TracedErrorf("command has no elements")
	}

	d.Command = command

	return nil
}

func (d *DockerRunContainerOptions) SetImageName(imageName string) (err error) {
	if imageName == "" {
		return tracederrors.TracedErrorf("imageName is empty string")
	}

	d.ImageName = imageName

	return nil
}

func (d *DockerRunContainerOptions) SetKeepStoppedContainer(keepStoppedContainer bool) (err error) {
	d.KeepStoppedContainer = keepStoppedContainer

	return nil
}

func (d *DockerRunContainerOptions) SetMounts(mounts []string) (err error) {
	if mounts == nil {
		return tracederrors.TracedErrorf("mounts is nil")
	}

	if len(mounts) <= 0 {
		return tracederrors.TracedErrorf("mounts has no elements")
	}

	d.Mounts = mounts

	return nil
}

func (d *DockerRunContainerOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	d.Name = name

	return nil
}

func (d *DockerRunContainerOptions) SetPorts(ports []string) (err error) {
	if ports == nil {
		return tracederrors.TracedErrorf("ports is nil")
	}

	if len(ports) <= 0 {
		return tracederrors.TracedErrorf("ports has no elements")
	}

	d.Ports = ports

	return nil
}

func (d *DockerRunContainerOptions) SetUseHostNet(useHostNet bool) (err error) {
	d.UseHostNet = useHostNet

	return nil
}

func (o *DockerRunContainerOptions) GetImageName() (imageName string, err error) {
	if len(o.ImageName) <= 0 {
		return "", tracederrors.TracedError("ImageName not set")
	}

	return o.ImageName, nil
}

func (o *DockerRunContainerOptions) GetName() (name string, err error) {
	if len(o.Name) <= 0 {
		return "", tracederrors.TracedError("Name not set")
	}

	return o.Name, nil
}

// Returns all ports numbers on the host.
// So for a "123:456" the returned port would be 123.
func (d *DockerRunContainerOptions) GetPortsOnHost() ([]int, error) {
	ret := []int{}

	for _, port := range d.Ports {
		port = strings.TrimPrefix(port, "0.0.0.0:")

		splitted := strings.Split(port, ":")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Failed to extract host port from '%s'", port)
		}

		p, err := strconv.Atoi(splitted[0])
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse '%s' as port: %w", splitted[0], err)
		}

		ret = append(ret, p)
	}

	slices.Sort(ret)

	return ret, nil
}
