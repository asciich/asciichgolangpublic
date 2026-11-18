# networkcardutils package

Utilities to work with network cards: configure layer 2 bridges and deactivate offloading.

## Subpackages

* [genericnetworkcardutils](./genericnetworkcardutils/README.md): Generic, pure in-memory helpers (no I/O).
* [nativenetworkcardutils](./nativenetworkcardutils/README.md): Native Go implementation (local machine only).
* [commandexecutornetowkrcardutils](./commandexecutornetowkrcardutils/README.md): Implementation using the command executor (works on remote machines).

## Examples

* [Create a layer 2 bridge](./Example_CreateLayer2Bridge_test.go)
