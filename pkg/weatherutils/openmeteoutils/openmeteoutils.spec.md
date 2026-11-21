# openmeteoutils specifications

This are the specifications for the [`openmeteoutils` package](README.md).

This document extends the parent specifications [weatherutils.spec.md](../weatherutils.spec.md).


## Implementation

- The `GetCurrentTemperature(ctx context.Context, lat float64, lon float64) (float64, error)` requests the current temperature at the given coordinates and returns them back in degree celsius.
- The API is free for non-commercial use, needs no API key and its data is licensed under CC BY 4.0, so please credit Open-Meteo. Therefore add a credit log.


## Testing

- Do at most one real request to `open-meteo.com` per unittests to avoid rate limiting.
