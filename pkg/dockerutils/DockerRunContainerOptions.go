package dockerutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DockerRunContainerOptions struct {
	ImageName               string
	Name                    string
	Command                 []string
	Ports                   []string
	Mounts                  []string
	KeepStoppedContainer    bool
	Verbose                 bool
	VerboseDockerRunCommand bool
	UseHostNet              bool
}

func NewDockerRunContainerOptions() (d *DockerRunContainerOptions) {
	return new(DockerRunContainerOptions)
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

func (d *DockerRunContainerOptions) GetVerbose() (verbose bool, err error) {

	return d.Verbose, nil
}

func (d *DockerRunContainerOptions) GetVerboseDockerRunCommand() (verboseDockerRunCommand bool, err error) {

	return d.VerboseDockerRunCommand, nil
}

func (d *DockerRunContainerOptions) MustGetCommand() (command []string) {
	command, err := d.GetCommand()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return command
}

func (d *DockerRunContainerOptions) MustGetImageName() (imageName string) {
	imageName, err := d.GetImageName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return imageName
}

func (d *DockerRunContainerOptions) MustGetKeepStoppedContainer() (keepStoppedContainer bool) {
	keepStoppedContainer, err := d.GetKeepStoppedContainer()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keepStoppedContainer
}

func (d *DockerRunContainerOptions) MustGetMounts() (mounts []string) {
	mounts, err := d.GetMounts()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mounts
}

func (d *DockerRunContainerOptions) MustGetName() (name string) {
	name, err := d.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (d *DockerRunContainerOptions) MustGetPorts() (ports []string) {
	ports, err := d.GetPorts()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ports
}

func (d *DockerRunContainerOptions) MustGetUseHostNet() (useHostNet bool) {
	useHostNet, err := d.GetUseHostNet()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return useHostNet
}

func (d *DockerRunContainerOptions) MustGetVerbose() (verbose bool) {
	verbose, err := d.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (d *DockerRunContainerOptions) MustGetVerboseDockerRunCommand() (verboseDockerRunCommand bool) {
	verboseDockerRunCommand, err := d.GetVerboseDockerRunCommand()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verboseDockerRunCommand
}

func (d *DockerRunContainerOptions) MustSetCommand(command []string) {
	err := d.SetCommand(command)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetImageName(imageName string) {
	err := d.SetImageName(imageName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetKeepStoppedContainer(keepStoppedContainer bool) {
	err := d.SetKeepStoppedContainer(keepStoppedContainer)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetMounts(mounts []string) {
	err := d.SetMounts(mounts)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetName(name string) {
	err := d.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetPorts(ports []string) {
	err := d.SetPorts(ports)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetUseHostNet(useHostNet bool) {
	err := d.SetUseHostNet(useHostNet)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetVerbose(verbose bool) {
	err := d.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DockerRunContainerOptions) MustSetVerboseDockerRunCommand(verboseDockerRunCommand bool) {
	err := d.SetVerboseDockerRunCommand(verboseDockerRunCommand)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (d *DockerRunContainerOptions) SetVerbose(verbose bool) (err error) {
	d.Verbose = verbose

	return nil
}

func (d *DockerRunContainerOptions) SetVerboseDockerRunCommand(verboseDockerRunCommand bool) (err error) {
	d.VerboseDockerRunCommand = verboseDockerRunCommand

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
