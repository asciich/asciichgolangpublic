# dockerutils

Work with docker.

There are two implementations available:
* [nativedocker](nativedocker/) using the official docker golang implementation. This is the preferred way if docker is locally running or you have access to the docker socket directly.
* [commandexecutordocker](commandexecutordocker/) using CLI commands. This is the preferred way if only the CLI on a remote host is available.

## Examples

* [Run container and exec](./Example_RunContainerAndExec_test.go)
* [Run container and exec command while streaming stdout](./Example_RunCommandAndGetStdoutAsIoReadCloser_test.go)
* [Run container and exec command while streaming stdin](./Example_RunCommandAndGetStdinAsIoWriteCloser_test.go)
* [Pull container image](Example_PullContainerImage_test.go)
