# dockerutils

Work with docker.

There are two implementations available:
* [nativedocker](nativedocker/) using the official docker golang implementaition. This is the prefered way if docker is locally running or you have access to the docker socket directly.
* [commandexecutordocker](commandexecutordocker/) using CLI commands. This is the prefered way if only the CLI on a remote host is available.
