# dockerutils specifications

This are the specifications for the [`dockerutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- For building containers use the `BuildContainerOptions` with at least these options:
    - `ImageNameAndTag`
    - For specifying the Dockerfile exclusively one must be set:
        - `DockerfilePath`: The path to an already existing Dockerfile.
        - `DockerfileContent`: Option to directly set the Dockerfile content.
- The `Container` interface must provide a `GetLogs(ctx context.Context) ([]byte, []byte, err)` function to return stdout and stderr.
    - Only return `nil` when there is an error. If fetching the log is successful return a empty `[]byte{}` if there is no stderr or stdout data.
- The `Container` interface must provide a `WaitUntilFinished(context.Context, time.Duration) error` function.
- The `Image` interface must provide a `Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error` function.
- The `Docker` interface must provide a `RunCommandInTemporaryContainer` function similar to the `RunCommandInTemporaryPod` in the `kubernetesutils` package.
    - This requires an additional test `Example_RunCommandInTemporaryContainer`.

## Testing

- Use small `alpine` based containers for testing wherever possible.
- All specifications have to be explicitly validated by unittests.
- All implementations must be tested in the this package using a test for all implementation types:
    ```golang
    func Test_TestName(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)
    ```
- Always implement all available implementation types. Avoid implementing if/elses depending on the implementation type:
    ```golang
    // Note: nativeDocker returns NotImplemented for all operations  # Avoid this special handling
    // so we skip validation tests for it                            # Avoid this special handling
    if tt.implementationName == "nativeDocker" {                    // Avoid this special handling
    ```
- Use `Wait...` functions instead of a hardcoded delay.:
    - Example:
        ```golang
        // Wait a bit for the container to finish # Avoid this, use the Wait... function instead.
        time.Sleep(time.Second * 2)              // Avoid this, use the Wait... function instead.
        ```
    - If it does not exist you have to implement it.
    - In any case ensure the `Wait...` functions are tested in all available implementations.
