package openmeteoutils_test

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/weatherutils/openmeteoutils"
)

func ExampleGetCurrentTemperature() {
	ctx := context.Background()

	// Request the current temperature near Trimbach (Olten, Switzerland).
	// The value is returned in degree celsius.
	temperature, err := openmeteoutils.GetCurrentTemperature(ctx, 47.36583, 7.903761)
	if err != nil {
		fmt.Printf("Unable to get the temperature: %v\n", err)
		return
	}

	fmt.Printf("Current temperature: %.1f °C\n", temperature)

	// No "Output:" comment on purpose: The result changes all the time
	// and the example must not perform a real request when running the unit tests.
}
