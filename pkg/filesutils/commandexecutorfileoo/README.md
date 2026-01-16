# commandexecutorfileoo

Object oriented file handling using a CommandExecutor.

A CommandExecutorFile implements the functionality of a `File` by executing commands (like: test, stat, cat...).
The benefit of this apporach is an easy way to access files on any
remote system like VMs, Containers, Hosts... while it easy to chain
like inside Container on VM behind Jumphost...

The downside of this is the poor performance and the possiblity to see
in the process table which operations where done.
