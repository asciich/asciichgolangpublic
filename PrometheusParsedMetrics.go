package asciichgolangpublic

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type PrometheusParsedMetrics struct {
	nativeMetricFamilies map[string]*dto.MetricFamily
}

func NewPrometheusParsedMetrics() (p *PrometheusParsedMetrics) {
	return new(PrometheusParsedMetrics)
}

// Get the value of the metric 'metricName' as float64.
// If the metric is not unique (e.g. a Vector with more than one value) this function will return an error.
func (p *PrometheusParsedMetrics) GetMetricValueAsFloat64(metricName string) (metricValue float64, err error) {
	if metricName == "" {
		return -1, errors.TracedErrorEmptyString("metricName")
	}

	nativeMetricFamilies, err := p.GetNativeMetricFamilies()
	if err != nil {
		return -1, err
	}

	for k, v := range nativeMetricFamilies {
		if k == metricName {
			metricFamily, err := GetPrometheusMetricFamilyFromNativeMetricFamily(v)
			if err != nil {
				return -1, err
			}

			metricValue, err := metricFamily.GetValueAsFloat64()
			if err != nil {
				return -1, err
			}

			return metricValue, nil
		}
	}

	return -1, errors.TracedErrorf("Metric '%s' not found.", metricName)
}

func (p *PrometheusParsedMetrics) GetNativeMetricFamilies() (nativeMetricFamilies map[string]*dto.MetricFamily, err error) {
	if p.nativeMetricFamilies == nil {
		return nil, errors.TracedError("nativeMetricFamilies not set")
	}
	return p.nativeMetricFamilies, nil
}

func (p *PrometheusParsedMetrics) MustGetMetricValueAsFloat64(metricName string) (metricValue float64) {
	metricValue, err := p.GetMetricValueAsFloat64(metricName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return metricValue
}

func (p *PrometheusParsedMetrics) MustGetNativeMetricFamilies() (nativeMetricFamilies map[string]*dto.MetricFamily) {
	nativeMetricFamilies, err := p.GetNativeMetricFamilies()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeMetricFamilies
}

func (p *PrometheusParsedMetrics) MustSetNativeMetricFamilies(nativeMetricFamilies map[string]*dto.MetricFamily) {
	err := p.SetNativeMetricFamilies(nativeMetricFamilies)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusParsedMetrics) SetNativeMetricFamilies(nativeMetricFamilies map[string]*dto.MetricFamily) (err error) {
	p.nativeMetricFamilies = nativeMetricFamilies

	return nil
}
