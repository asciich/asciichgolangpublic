# weatherutilscmd specifications

This are the specifications for the [`weatherutilscmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../defaultclicommands.spec.md).

## Implementation

- Provide the command `get-current-temperature` which accepts `--lat` and `--lon`. It must wire up the `weatherutils.GetCurrentTemperature` function.
