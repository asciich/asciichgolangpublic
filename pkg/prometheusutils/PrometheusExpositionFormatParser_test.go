package prometheusutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestPrometheusExpositionFormatParserParseExample(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				gitRepo := asciichgolangpublic.MustGetLocalGitRepositoryByPath(".")
				metricsTxt := gitRepo.MustReadFileInDirectoryAsString("testdata", "PrometheusExpositionFormatParser", "metrics.txt")

				parsedMetrics := PrometheusExpositionFormatParser().MustParseString(metricsTxt)

				require.EqualValues(
					12.47,
					parsedMetrics.MustGetMetricValueAsFloat64("metric_without_timestamp_and_labels"),
				)

				require.EqualValues(
					1.458255915e9,
					parsedMetrics.MustGetMetricValueAsFloat64("msdos_file_access_time_seconds"),
				)
			},
		)
	}
}

func TestPrometheusExpositionFormatParserParseGauge(t *testing.T) {
	tests := []struct {
		expectedValue float64
	}{
		{1.0},
		{-1.0},
		{-123},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				metricsTxt := ""
				metricsTxt += "# HELP abc_value help text\n"
				metricsTxt += "# TYPE abc_value gauge\n"
				metricsTxt += fmt.Sprintf("abc_value %f\n", tt.expectedValue)

				parsedMetrics := PrometheusExpositionFormatParser().MustParseString(metricsTxt)

				require.EqualValues(
					tt.expectedValue,
					parsedMetrics.MustGetMetricValueAsFloat64("abc_value"),
				)
			},
		)
	}
}

func TestPrometheusExpositionFormatParserParseCounter(t *testing.T) {
	tests := []struct {
		expectedValue float64
	}{
		{11},
		{-1.0},
		{-123},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				metricsTxt := ""
				metricsTxt += "# HELP abc_value help text\n"
				metricsTxt += "# TYPE abc_value counter\n"
				metricsTxt += fmt.Sprintf("abc_value %f\n", tt.expectedValue)

				parsedMetrics := PrometheusExpositionFormatParser().MustParseString(metricsTxt)

				require.EqualValues(
					tt.expectedValue,
					parsedMetrics.MustGetMetricValueAsFloat64("abc_value"),
				)
			},
		)
	}
}
