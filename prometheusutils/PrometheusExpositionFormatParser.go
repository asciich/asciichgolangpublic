package prometheusutils

import (
	"strings"

	"github.com/prometheus/common/expfmt"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type PrometheusExpositionFormatParserService struct{}

// Can be used to parse Prometheus exposition format as documented here:
//   - https://prometheus.io/docs/instrumenting/exposition_formats/
func PrometheusExpositionFormatParser() (p *PrometheusExpositionFormatParserService) {
	return NewPrometheusExpositionFormatParserService()
}

func NewPrometheusExpositionFormatParserService() (p *PrometheusExpositionFormatParserService) {
	return new(PrometheusExpositionFormatParserService)
}

func (p *PrometheusExpositionFormatParserService) MustParseString(toParse string) (parsed *PrometheusParsedMetrics) {
	parsed, err := p.ParseString(toParse)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parsed
}

func (p *PrometheusExpositionFormatParserService) ParseString(toParse string) (parsed *PrometheusParsedMetrics, err error) {
	if toParse == "" {
		return nil, tracederrors.TracedErrorEmptyString("toParse")
	}

	stringReader := strings.NewReader(toParse)

	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(stringReader)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Failed to parse text into metric families: '%w'",
			err,
		)
	}

	parsed = NewPrometheusParsedMetrics()

	err = parsed.SetNativeMetricFamilies(metricFamilies)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
