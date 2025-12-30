package dockeroptions

import (
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
