### `weather get-current-temperature`

Prints the current temperature in degree celsius at the given coordinates.

```bash
<binary> weather get-current-temperature --lat <latitude> --lon <longitude>
```

- `--lat`: Latitude in decimal degrees (-90 to 90). Required.
- `--lon`: Longitude in decimal degrees (-180 to 180). Required.

Only the temperature is written to stdout.

#### Example: Current temperature in Bern

```bash
<binary> weather get-current-temperature --lat 46.948 --lon 7.4474
```

Example output (the value changes with the weather):

```
5.8
```

Use it in a script:

```bash
temperature=$(<binary> weather get-current-temperature --lat 46.948 --lon 7.4474)
echo "It is currently ${temperature} °C in Bern."
```

The data is provided by [open-meteo.com](https://open-meteo.com/). It is free for non-commercial use and licensed under CC BY 4.0, so please credit Open-Meteo.
