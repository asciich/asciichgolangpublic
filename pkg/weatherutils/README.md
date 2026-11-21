# weatherutils

Convenience functions to get weather data.

The functions delegate to the implementations in the subpackages:

- [openmeteoutils](openmeteoutils/README.md): Weather data from [open-meteo.com](https://open-meteo.com/).

## Usage

```golang
import (
    "github.com/asciich/asciichgolangpublic/pkg/weatherutils"
)

temperature, err := weatherutils.GetCurrentTemperature(ctx, 47.36583, 7.903761)
```

The returned value is the current temperature in degree celsius. It is modeled data for the given coordinates and not the reading of a local sensor.

The data is provided by Open-Meteo. It is licensed under CC BY 4.0, so please credit Open-Meteo. A credit message is logged after every successful request.

## Examples

- [ExampleGetCurrentTemperature](Example_GetCurrentTemperature_test.go)

## Specifications

For specifications see [weatherutils.spec.md](weatherutils.spec.md)
