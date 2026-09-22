package ollamautils

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// ollamaVolumeMount is the shared volume mount used by all run commands.
	ollamaVolumeMount = "ollama:/root/.ollama"

	// ollamaContainerName is the docker container name used by all run commands.
	ollamaContainerName = "ollama"

	// ollamaImageCpuNvidia is the default Ollama image (CPU and NVIDIA GPU).
	ollamaImageCpuNvidia = "ollama/ollama"

	// ollamaImageAmd is the ROCm enabled Ollama image for AMD GPUs.
	ollamaImageAmd = "ollama/ollama:rocm"

	// ollamaDefaultContextLength is the default context window (num_ctx) used
	// for all run commands when no explicit value is provided via RunOptions.
	// It is passed to the container as the OLLAMA_CONTEXT_LENGTH environment
	// variable.
	ollamaDefaultContextLength = 32768
)

// ErrNoGpuDetected is returned by RunGPU when no supported GPU could be
// detected. Falling back to CPU only mode is intentionally not done.
var ErrNoGpuDetected = errors.New("no supported GPU detected")

// RunOptions holds the adjustable settings for the various Run* commands.
//
// It is intentionally passed as an optional variadic argument so existing
// callers using e.g. RunCpuOnly(ctx) keep working while new callers can
// override defaults using RunCpuOnly(ctx, &RunOptions{ContextLength: 8192}).
type RunOptions struct {
	// ContextLength sets the OLLAMA_CONTEXT_LENGTH environment variable inside
	// the container. When it is <= 0 (or no options are given at all) the
	// ollamaDefaultContextLength is used instead.
	ContextLength int
}

// resolveContextLength returns the effective context length to use based on the
// (optional) provided options, falling back to ollamaDefaultContextLength.
func resolveContextLength(options ...*RunOptions) int {
	for _, o := range options {
		if o != nil && o.ContextLength > 0 {
			return o.ContextLength
		}
	}

	return ollamaDefaultContextLength
}

// pathExists reports whether the given path exists on the target machine.
//
// It intentionally does not rely on the exit code of a single command:
// the shell always exits 0 and prints "yes" or "no". This guarantees the
// command execution itself was performed and lets us distinguish a real
// execution error from a "does not exist" result.
func pathExists(ctx context.Context, path string) (ret bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	output, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", "test -e '" + path + "' && echo yes || echo no"},
		},
	)
	if err != nil {
		return false, err
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	switch stdout {
	case "yes":
		ret = true
	case "no":
		ret = false
	default:
		return false, tracederrors.TracedErrorf("Unexpected output when checking if path '%s' exists: '%s'", path, stdout)
	}

	return ret, nil
}

// isOllamaContainerRunning reports whether the ollama container is already running.
func isOllamaContainerRunning(ctx context.Context) (ret bool, err error) {
	output, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"docker ps --filter 'name=^" + ollamaContainerName + "$' --format '{{.Names}}' | grep -qx '" + ollamaContainerName + "' && echo yes || echo no",
			},
		},
	)
	if err != nil {
		return false, err
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	switch stdout {
	case "yes":
		ret = true
	case "no":
		ret = false
	default:
		return false, tracederrors.TracedErrorf("Unexpected output when checking if ollama container is running: '%s'", stdout)
	}

	return ret, nil
}

// runOllamaContainer runs the ollama docker container with the given extra
// arguments (e.g. GPU related flags), image and context length in an
// idempotent way.
func runOllamaContainer(ctx context.Context, extraArgs []string, image string, contextLength int) error {
	if image == "" {
		return tracederrors.TracedErrorEmptyString("image")
	}

	if contextLength <= 0 {
		return tracederrors.TracedErrorf("contextLength must be greater than 0 but got '%d'", contextLength)
	}

	isRunning, err := isOllamaContainerRunning(ctx)
	if err != nil {
		return err
	}

	if isRunning {
		logging.LogInfoByCtxf(ctx, "Ollama container '%s' is already running. Nothing to do.", ollamaContainerName)
		return nil
	}

	command := []string{"docker", "run", "-d"}
	command = append(command, extraArgs...)
	command = append(command,
		"-e", "OLLAMA_CONTEXT_LENGTH="+strconv.Itoa(contextLength),
		"-v", ollamaVolumeMount,
		"-p", "11434:11434",
		"--name", ollamaContainerName,
		image,
	)

	_, err = commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Started ollama container '%s' using image '%s' with context length '%d'.", ollamaContainerName, image, contextLength)

	return nil
}

// RunCpuOnly starts ollama in a docker container in CPU only / no GPU mode.
func RunCpuOnly(ctx context.Context, options ...*RunOptions) error {
	logging.LogInfoByCtxf(ctx, "Run ollama in cpu only mode started.")

	contextLength := resolveContextLength(options...)

	err := runOllamaContainer(ctx, nil, ollamaImageCpuNvidia, contextLength)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ollama in cpu only mode finished.")

	return nil
}

// RunGPUNvidia starts ollama in a docker container with nvidia GPU support.
func RunGPUNvidia(ctx context.Context, options ...*RunOptions) error {
	logging.LogInfoByCtxf(ctx, "Run ollama with nvidia GPU support started.")

	contextLength := resolveContextLength(options...)

	extraArgs := []string{"--gpus", "all"}

	err := runOllamaContainer(ctx, extraArgs, ollamaImageCpuNvidia, contextLength)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ollama with nvidia GPU support finished.")

	return nil
}

// RunGPUAmd starts ollama in a docker container with amd (ROCm) GPU support.
func RunGPUAmd(ctx context.Context, options ...*RunOptions) error {
	logging.LogInfoByCtxf(ctx, "Run ollama with amd GPU support started.")

	contextLength := resolveContextLength(options...)

	// The numeric GIDs of 'video' and 'render' from the host are required to
	// access /dev/kfd and /dev/dri. docker's --group-add resolves *names*
	// against the container image (where these groups do not exist), so the
	// numeric id from the host must be passed instead.
	videoGid, err := resolveGroupId(ctx, "video")
	if err != nil {
		return err
	}

	renderGid, err := resolveGroupId(ctx, "render")
	if err != nil {
		return err
	}

	// AMD GPUs are accessed by passing through the kernel fusion driver
	// (/dev/kfd) and the direct rendering devices (/dev/dri) instead of
	// using a dedicated container runtime like nvidia does.
	extraArgs := []string{
		"--device", "/dev/kfd",
		"--device", "/dev/dri",
		"--group-add", videoGid,
		"--group-add", renderGid,
		"--security-opt", "seccomp=unconfined",
	}

	err = runOllamaContainer(ctx, extraArgs, ollamaImageAmd, contextLength)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Run ollama with amd GPU support finished.")

	return nil
}

// RunGPU starts ollama in a docker container with autodetected GPU support.
// It prefers AMD (ROCm) if an AMD GPU is detected, then nvidia.
// If no supported GPU is detected ErrNoGpuDetected is returned; falling back
// to CPU only mode is intentionally not an option.
func RunGPU(ctx context.Context, options ...*RunOptions) error {
	logging.LogInfoByCtxf(ctx, "Run ollama with autodetected GPU support started.")

	amdAvailable, err := IsAmdGpuAvailable(ctx)
	if err != nil {
		return err
	}

	if amdAvailable {
		logging.LogInfoByCtxf(ctx, "AMD GPU detected. Using amd GPU support.")

		err = RunGPUAmd(ctx, options...)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Run ollama with autodetected GPU support finished.")

		return nil
	}

	nvidiaAvailable, err := IsNvidiaGpuAvailable(ctx)
	if err != nil {
		return err
	}

	if nvidiaAvailable {
		logging.LogInfoByCtxf(ctx, "Nvidia GPU detected. Using nvidia GPU support.")

		err = RunGPUNvidia(ctx, options...)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Run ollama with autodetected GPU support finished.")

		return nil
	}

	return tracederrors.TracedError(ErrNoGpuDetected)
}

// IsAmdGpuAvailable reports whether an AMD GPU usable by ROCm is present.
// The presence of the kernel fusion driver device node (/dev/kfd) together
// with the direct rendering infrastructure (/dev/dri) is a reliable indicator.
func IsAmdGpuAvailable(ctx context.Context) (ret bool, err error) {
	kfdExists, err := pathExists(ctx, "/dev/kfd")
	if err != nil {
		return false, err
	}

	driExists, err := pathExists(ctx, "/dev/dri")
	if err != nil {
		return false, err
	}

	ret = kfdExists && driExists

	logging.LogInfoByCtxf(ctx, "AMD GPU available: '%v'.", ret)

	return ret, nil
}

// IsNvidiaGpuAvailable reports whether an nvidia GPU is present by checking
// for the device nodes created by the nvidia kernel driver.
func IsNvidiaGpuAvailable(ctx context.Context) (ret bool, err error) {
	nvidia0Exists, err := pathExists(ctx, "/dev/nvidia0")
	if err != nil {
		return false, err
	}

	nvidiaCtlExists, err := pathExists(ctx, "/dev/nvidiactl")
	if err != nil {
		return false, err
	}

	ret = nvidia0Exists || nvidiaCtlExists

	logging.LogInfoByCtxf(ctx, "Nvidia GPU available: '%v'.", ret)

	return ret, nil
}

// resolveGroupId resolves the numeric group id (GID) of a group on the target
// machine by its name. The numeric id must be used with docker's
// '--group-add' because docker resolves group *names* against the container
// image's group file, where 'render'/'video' usually do not exist.
func resolveGroupId(ctx context.Context, groupName string) (ret string, err error) {
	if groupName == "" {
		return "", tracederrors.TracedErrorEmptyString("groupName")
	}

	// 'getent group <name>' prints "name:x:GID:members" and exits 0 when the
	// group exists. We wrap it so the execution itself always succeeds and we
	// can distinguish "not found" from an execution error.
	output, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"getent group '" + groupName + "' | cut -d: -f3 || true",
			},
		},
	)
	if err != nil {
		return "", err
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	ret = strings.TrimSpace(stdout)
	if ret == "" {
		return "", tracederrors.TracedErrorf("Group '%s' not found on the target machine.", groupName)
	}

	logging.LogInfoByCtxf(ctx, "Resolved group '%s' to gid '%s'.", groupName, ret)

	return ret, nil
}
