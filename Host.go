package asciichgolangpublic

import (
	"fmt"
	"strings"
	"time"
)

// Host represents a classical host like a VM, Laptop, Desktop, Server.
type Host struct {
	CommandExecutorBase
	hostname    string
	sshUsername string
	Comment     string
}

func GetHostByHostname(hostname string) (host *Host, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, TracedError("hostname is empty string")
	}

	host = NewHost()
	err = host.SetHostname(hostname)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func MustGetHostByHostname(hostname string) (host *Host) {
	host, err := GetHostByHostname(hostname)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return host
}

func NewHost() (host *Host) {
	host = new(Host)

	host.MustSetParentCommandExecutorForBaseClass(host)

	return host
}

func (h *Host) AddSshHostKeyToKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: []string{
				fmt.Sprintf("ssh-keyscan -H '%s' >> ${HOME}/.ssh/known_hosts", hostname),
			},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Added host key of '%s' from known hosts", hostname)
	}

	return nil
}

func (h *Host) CheckFtpPortOpen(verbose bool) (err error) {
	isOpen, err := h.IsFtpPortOpen(verbose)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	if isOpen {
		if verbose {
			LogInfof("FTP port on host '%s' is open.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("FTP port on host '%s' is not open.", hostname)
		if verbose {
			LogError(errorMessage)
		}

		return TracedError(errorMessage)
	}

	return nil
}

func (h *Host) CheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool, err error) {
	kubernetesControlplane, err := h.GetAsKubernetesControlplaneHost()
	if err != nil {
		return false, err
	}

	isKubernetesControlplane, err = kubernetesControlplane.CheckIsKubernetesControlplane(verbose)
	if err != nil {
		return false, err
	}

	return isKubernetesControlplane, nil
}

func (h *Host) CheckIsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	kubernetesNode, err := h.GetAsKubernetesNodeHost()
	if err != nil {
		return false, err
	}

	isKubernetesNode, err = kubernetesNode.CheckIsKubernetesNode(verbose)
	if err != nil {
		return false, err
	}

	return isKubernetesNode, nil
}

func (h *Host) CheckReachableBySsh(verbose bool) (err error) {
	isReachable, err := h.IsReachableBySsh(verbose)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	if isReachable {
		if verbose {
			LogInfof("Host '%s' is reachable by SSH.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("Host '%s' is reachable by SSH.", hostname)
		if verbose {
			LogError(errorMessage)
		}

		return TracedError(errorMessage)
	}

	return nil
}

func (h *Host) GetAsKubernetesControlplaneHost() (kubernetesControlPlaneHost *KubernetesControlplaneHost, err error) {
	kubernetesControlPlaneHost = NewKubernetesControlplaneHost()

	hostname, err := h.GetHostname()
	if err != nil {
		return nil, err
	}

	err = kubernetesControlPlaneHost.SetHostname(hostname)
	if err != nil {
		return nil, err
	}

	if h.IsSshUserNameSet() {
		sshUserName, err := h.GetSshUserName()
		if err != nil {
			return nil, err
		}

		err = kubernetesControlPlaneHost.SetSshUserName(sshUserName)
		if err != nil {
			return nil, err
		}
	}

	return kubernetesControlPlaneHost, nil
}

func (h *Host) GetAsKubernetesNodeHost() (kubernetesNodeHost *KubernetesNodeHost, err error) {
	kubernetesNodeHost = NewKubernetesNodeHost()

	hostname, err := h.GetHostname()
	if err != nil {
		return nil, err
	}

	err = kubernetesNodeHost.SetHostname(hostname)
	if err != nil {
		return nil, err
	}

	if h.IsSshUserNameSet() {
		sshUserName, err := h.GetSshUserName()
		if err != nil {
			return nil, err
		}

		err = kubernetesNodeHost.SetSshUserName(sshUserName)
		if err != nil {
			return nil, err
		}
	}

	return kubernetesNodeHost, nil
}

func (h *Host) GetComment() (comment string, err error) {
	if h.Comment == "" {
		return "", TracedErrorf("Comment not set")
	}

	return h.Comment, nil
}

func (h *Host) GetDeepCopy() (deepCopy CommandExecutor) {
	d := NewHost()

	*d = *h

	deepCopy = d

	return deepCopy
}

func (h *Host) GetDirectoryByPath(path string) (directory Directory, err error) {
	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	commandExecutorDir, err := NewCommandExecutorDirectory(h)
	if err != nil {
		return nil, err
	}

	err = commandExecutorDir.SetDirPath(path)
	if err != nil {
		return nil, err
	}

	return commandExecutorDir, nil
}

func (h *Host) GetDockerContainerByName(containerName string) (dockerContainer *DockerContainer, err error) {
	if len(containerName) <= 0 {
		return nil, TracedError("containerName is empty string")
	}

	dockerService, err := h.GetDockerService()
	if err != nil {
		return nil, err
	}

	dockerContainer, err = dockerService.GetDockerContainerByName(containerName)
	if err != nil {
		return nil, err
	}

	return dockerContainer, nil
}

func (h *Host) GetDockerService() (dockerService *DockerService, err error) {
	dockerService = NewDockerService()

	err = dockerService.SetHost(h)
	if err != nil {
		return nil, err
	}

	return dockerService, nil
}

func (h *Host) GetHostDescription() (hostDescription string, err error) {
	hostDescription, err = h.GetHostname()
	if err != nil {
		return "", err
	}

	return hostDescription, nil
}

func (h *Host) GetHostname() (hostname string, err error) {
	if len(h.hostname) <= 0 {
		return "", TracedError("hostname not set")
	}

	return h.hostname, nil
}

func (h *Host) GetKvmHypervisor() (kvmHypervisor *KVMHypervisor, err error) {
	kvmHypervisor, err = GetKvmHypervisorByHost(h)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

func (h *Host) GetKvmStoragePoolNames(verbose bool) (storagePoolNames []string, err error) {
	kvmHypervisor, err := h.GetKvmHypervisor()
	if err != nil {
		return nil, err
	}

	storagePoolNames, err = kvmHypervisor.GetStoragePoolNames(verbose)
	if err != nil {
		return nil, err
	}

	return storagePoolNames, nil
}

func (h *Host) GetKvmVmsNames(verbose bool) (vmNames []string, err error) {
	kvmHypervisor, err := h.GetKvmHypervisor()
	if err != nil {
		return nil, err
	}

	vmNames, err = kvmHypervisor.GetVmNames(verbose)
	if err != nil {
		return nil, err
	}

	return vmNames, nil
}

func (h *Host) GetKvmVolumeNames(verbose bool) (volumeNames []string, err error) {
	hypervisor, err := h.GetKvmHypervisor()
	if err != nil {
		return nil, err
	}

	volumeNames, err = hypervisor.GetVolumeNames(verbose)
	if err != nil {
		return nil, err
	}

	return volumeNames, nil
}

func (h *Host) GetSSHClient() (sshClient *SSHClient, err error) {
	sshClient = NewSSHClient()
	err = sshClient.SetHost(h)
	if err != nil {
		return nil, err
	}

	if h.IsSshUserNameSet() {
		sshUserName, err := h.GetSshUserName()
		if err != nil {
			return nil, err
		}

		err = sshClient.SetSshUserName(sshUserName)
		if err != nil {
			return nil, err
		}
	}

	return sshClient, nil
}

func (h *Host) GetSshUserName() (sshUserName string, err error) {
	if !h.IsSshUserNameSet() {
		return "", TracedError("sshUserName not set")
	}

	return h.sshUsername, nil
}

func (h *Host) GetSshUsername() (sshUsername string, err error) {
	if h.sshUsername == "" {
		return "", TracedErrorf("sshUsername not set")
	}

	return h.sshUsername, nil
}

func (h *Host) InstallBinary(installOptions *InstallOptions) (err error) {
	if installOptions == nil {
		return TracedErrorNil("installOptions")
	}

	hostName, err := h.GetHostname()
	if err != nil {
		return err
	}

	sourceFile, err := installOptions.GetSourceFile()
	if err != nil {
		return err
	}

	sourceFilePath, err := sourceFile.GetLocalPath()
	if err != nil {
		return err
	}

	binaryName, err := installOptions.GetBinaryName()
	if err != nil {
		return err
	}

	if installOptions.Verbose {
		LogInfof(
			"Install '%s' as '%s' on host '%s' started.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	tempCopy, err := TemporaryFiles().CreateTemporaryFileFromFile(sourceFile, installOptions.Verbose)
	if err != nil {
		return err
	}

	destPath, err := installOptions.GetInstallationPathOrDefaultIfUnset()
	if err != nil {
		return err
	}

	installedFile, err := tempCopy.MoveToPath(destPath, installOptions.UseSudoToInstall, installOptions.Verbose)
	if err != nil {
		return err
	}

	err = installedFile.Chmod(
		&ChmodOptions{
			PermissionsString: "u=rwx,g=rx,o=rx",
			UseSudo: installOptions.UseSudoToInstall,
			Verbose: installOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	err = installedFile.Chown(
		&ChownOptions{
			UserName:  "root",
			GroupName: "root",
			UseSudo:   installOptions.UseSudoToInstall,
			Verbose:   installOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if installOptions.Verbose {
		LogChangedf(
			"Install '%s' as '%s' on host '%s' finished.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	return nil
}

func (h *Host) IsFtpPortOpen(verbose bool) (isOpen bool, err error) {
	isOpen, err = h.IsTcpPortOpen(FTP().GetDefaultPort(), verbose)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (h *Host) IsPingable(verbose bool) (isPingable bool, err error) {
	hostname, err := h.GetHostname()
	if err != nil {
		return false, err
	}

	stdout, err := Bash().RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{"bash", "-c", fmt.Sprintf("ping -c 1 '%s' &>/dev/null && echo yes || echo no", hostname)},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)
	if stdout == "yes" {
		return true, nil
	}
	if stdout == "no" {
		return false, nil
	}

	return false, TracedErrorf("Unexpected stdout: '%v'", stdout)
}

func (h *Host) IsReachableBySsh(verbose bool) (isReachable bool, err error) {
	sshClient, err := h.GetSSHClient()
	if err != nil {
		return false, err
	}

	isReachable, err = sshClient.IsReachable(verbose)
	if err != nil {
		return false, err
	}

	return isReachable, nil
}

func (h *Host) IsSshUserNameSet() (isSet bool) {
	return len(h.sshUsername) > 0
}

func (h *Host) IsTcpPortOpen(portNumber int, verbose bool) (isOpen bool, err error) {
	if portNumber <= 0 {
		return false, TracedErrorf("Invalid portNumber: '%d'", portNumber)
	}

	hostname, err := h.GetHostname()
	if err != nil {
		return false, err
	}

	isOpen, err = TcpPorts().IsPortOpen(hostname, portNumber, verbose)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (h *Host) MustAddSshHostKeyToKnownHosts(verbose bool) {
	err := h.AddSshHostKeyToKnownHosts(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustCheckFtpPortOpen(verbose bool) {
	err := h.CheckFtpPortOpen(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustCheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool) {
	isKubernetesControlplane, err := h.CheckIsKubernetesControlplane(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesControlplane
}

func (h *Host) MustCheckIsKubernetesNode(verbose bool) (isKubernetesNode bool) {
	isKubernetesNode, err := h.CheckIsKubernetesNode(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesNode
}

func (h *Host) MustCheckReachableBySsh(verbose bool) {
	err := h.CheckReachableBySsh(verbose)
	if err != nil {
		LogFatalf("host.CheckReachableBySsh failed: '%v'", err)
	}
}

func (h *Host) MustGetAsKubernetesControlplaneHost() (kubernetesControlPlaneHost *KubernetesControlplaneHost) {
	kubernetesControlPlaneHost, err := h.GetAsKubernetesControlplaneHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return kubernetesControlPlaneHost
}

func (h *Host) MustGetAsKubernetesNodeHost() (kubernetesNodeHost *KubernetesNodeHost) {
	kubernetesNodeHost, err := h.GetAsKubernetesNodeHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return kubernetesNodeHost
}

func (h *Host) MustGetComment() (comment string) {
	comment, err := h.GetComment()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return comment
}

func (h *Host) MustGetDirectoryByPath(path string) (directory Directory) {
	directory, err := h.GetDirectoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return directory
}

func (h *Host) MustGetDockerContainerByName(containerName string) (dockerContainer *DockerContainer) {
	dockerContainer, err := h.GetDockerContainerByName(containerName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dockerContainer
}

func (h *Host) MustGetDockerService() (dockerService *DockerService) {
	dockerService, err := h.GetDockerService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dockerService
}

func (h *Host) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := h.GetHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostDescription
}

func (h *Host) MustGetHostname() (hostname string) {
	hostname, err := h.GetHostname()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostname
}

func (h *Host) MustGetKvmHypervisor() (kvmHypervisor *KVMHypervisor) {
	kvmHypervisor, err := h.GetKvmHypervisor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return kvmHypervisor
}

func (h *Host) MustGetKvmStoraPoolNames(verbose bool) (storagePoolNames []string) {
	storagePoolNames, err := h.GetKvmStoragePoolNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return storagePoolNames
}

func (h *Host) MustGetKvmStoragePoolNames(verbose bool) (storagePoolNames []string) {
	storagePoolNames, err := h.GetKvmStoragePoolNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return storagePoolNames
}

func (h *Host) MustGetKvmVmsNames(verbose bool) (vmNames []string) {
	vmNames, err := h.GetKvmVmsNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return vmNames
}

func (h *Host) MustGetKvmVolumeNames(verbose bool) (volumeNames []string) {
	volumeNames, err := h.GetKvmVolumeNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return volumeNames
}

func (h *Host) MustGetSSHClient() (sshClient *SSHClient) {
	sshClient, err := h.GetSSHClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshClient
}

func (h *Host) MustGetSshUserName() (sshUserName string) {
	sshUserName, err := h.GetSshUserName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshUserName
}

func (h *Host) MustGetSshUsername() (sshUsername string) {
	sshUsername, err := h.GetSshUsername()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshUsername
}

func (h *Host) MustInstallBinary(installOptions *InstallOptions) {
	err := h.InstallBinary(installOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustIsFtpPortOpen(verbose bool) (isOpen bool) {
	isOpen, err := h.IsFtpPortOpen(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *Host) MustIsPingable(verbose bool) (isPingable bool) {
	isPingable, err := h.IsPingable(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isPingable
}

func (h *Host) MustIsReachableBySsh(verbose bool) (isReachable bool) {
	isReachable, err := h.IsReachableBySsh(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isReachable
}

func (h *Host) MustIsTcpPortOpen(portNumber int, verbose bool) (isOpen bool) {
	isOpen, err := h.IsTcpPortOpen(portNumber, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *Host) MustRemoveKvmVm(removeOptions *KvmRemoveVmOptions) {
	err := h.RemoveKvmVm(removeOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustRemoveSshHostKeyFromKnownHosts(verbose bool) {
	err := h.RemoveSshHostKeyFromKnownHosts(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustRenewSshHostKey(verbose bool) {
	err := h.RenewSshHostKey(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := h.RunCommand(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (h *Host) MustSetComment(comment string) {
	err := h.SetComment(comment)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustSetHostname(hostname string) {
	err := h.SetHostname(hostname)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustSetSshUserName(username string) {
	err := h.SetSshUserName(username)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustSetSshUsername(sshUsername string) {
	err := h.SetSshUsername(sshUsername)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustWaitUntilPingable(verbose bool) {
	err := h.WaitUntilPingable(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) MustWaitUntilReachableBySsh(renewHostKey bool, verbose bool) {
	err := h.WaitUntilReachableBySsh(renewHostKey, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (h *Host) RemoveKvmVm(removeOptions *KvmRemoveVmOptions) (err error) {
	if removeOptions == nil {
		return TracedError("removeOptions is nil")
	}

	if len(removeOptions.VmName) <= 0 {
		return TracedError("removeOptions.VmName is empty string")
	}

	hypervisor, err := h.GetKvmHypervisor()
	if err != nil {
		return err
	}

	err = hypervisor.RemoveVm(removeOptions)
	if err != nil {
		return err
	}

	return nil
}

func (h *Host) RemoveSshHostKeyFromKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: []string{"ssh-keygen", "-R", hostname},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Removed host key of '%s' from known hosts", hostname)
	}

	return nil
}

func (h *Host) RenewSshHostKey(verbose bool) (err error) {
	err = h.RemoveSshHostKeyFromKnownHosts(verbose)
	if err != nil {
		return err
	}

	err = h.AddSshHostKeyToKnownHosts(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (h *Host) RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	sshClient, err := h.GetSSHClient()
	if err != nil {
		return nil, err
	}

	commandOutput, err = sshClient.RunCommand(options)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (h *Host) SetComment(comment string) (err error) {
	if comment == "" {
		return TracedErrorf("comment is empty string")
	}

	h.Comment = comment

	return nil
}

func (h *Host) SetHostname(hostname string) (err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return TracedError("hostname is empty string")
	}

	h.hostname = hostname

	return nil
}

func (h *Host) SetSshUserName(username string) (err error) {
	if len(username) <= 0 {
		return TracedError("username is empty string")
	}

	h.sshUsername = username

	return nil
}

func (h *Host) SetSshUsername(sshUsername string) (err error) {
	if sshUsername == "" {
		return TracedErrorf("sshUsername is empty string")
	}

	h.sshUsername = sshUsername

	return nil
}

func (h *Host) WaitUntilPingable(verbose bool) (err error) {
	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		isPingable, err := h.IsPingable(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isPingable {
			if verbose {
				LogGoodf("Host '%s' is pingable after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not pingable after '%v'", hostname, elapsedTime)
			if verbose {
				LogError(errorMessage)
			}
			return TracedError(errorMessage)
		}

		if verbose {
			LogInfof(
				"Wait '%v' for host '%s' to get reachable by ping. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}

func (h *Host) WaitUntilReachableBySsh(renewHostKey bool, verbose bool) (err error) {
	hostname, err := h.GetHostname()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		if renewHostKey {
			err = h.RenewSshHostKey(verbose)
			if err != nil {
				LogWarn("Renewing host key failed, but error is ignored in WaitUntilReachableBySsh since running in a retry loop.")
			}
		}

		isReachableBySsh, err := h.IsReachableBySsh(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isReachableBySsh {
			if verbose {
				LogGoodf("Host '%s' is reachable by SSH after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable by SSH after '%v'", hostname, elapsedTime)
			if verbose {
				LogError(errorMessage)
			}
			return TracedError(errorMessage)
		}

		if verbose {
			LogInfof(
				"Wait '%v' for host '%s' to get reachable by SSH. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}
