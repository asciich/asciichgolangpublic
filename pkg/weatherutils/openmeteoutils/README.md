# openmeteoutils

Get weather data from [open-meteo.com](https://open-meteo.com/). The API is free for non-commercial use and needs no API key. Its data is licensed under CC BY 4.0, so please credit Open-Meteo.

To help with the attribution, every successful request logs a credit message.

## Usage

```golang
import (
    "github.com/asciich/asciichgolangpublic/pkg/weatherutils/openmeteoutils"
)

temperature, err := openmeteoutils.GetCurrentTemperature(ctx, 47.36583, 7.903761)
```

The returned value is modeled data for the given coordinates in degree celsius. It is not the reading of a local sensor.

## Examples

- [ExampleGetCurrentTemperature](Example_GetCurrentTemperature_test.go)

## Specifications

For specifications see [openmeteoutils.spec.md](openmeteoutils.spec.md)
