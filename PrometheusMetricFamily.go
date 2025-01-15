package asciichgolangpublic

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type PrometheusMetricFamily struct {
	name    string
	help    string
	metrics []*PrometheusMetric
}

func GetPrometheusMetricFamilyFromNativeMetricFamily(metricFamily *dto.MetricFamily) (p *PrometheusMetricFamily, err error) {
	if metricFamily == nil {
		return nil, tracederrors.TracedErrorNil("metricFamily")
	}

	p = NewPrometheusMetricFamily()

	err = p.SetName(*metricFamily.Name)
	if err != nil {
		return nil, err
	}

	for _, metric := range metricFamily.Metric {
		toAdd := NewPrometheusMetric()

		untyped := metric.GetUntyped()
		gauge := metric.GetGauge()
		if untyped != nil {
			err = toAdd.SetValueByFloat64(untyped.GetValue())
			if err != nil {
				return nil, err
			}
		} else if gauge != nil {
			err = toAdd.SetValueByFloat64(gauge.GetValue())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, tracederrors.TracedErrorf("Not implemented for '%v'", metric)
		}

		err = p.AddMetric(toAdd)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func MustGetPrometheusMetricFamilyFromNativeMetricFamily(metricFamily *dto.MetricFamily) (p *PrometheusMetricFamily) {
	p, err := GetPrometheusMetricFamilyFromNativeMetricFamily(metricFamily)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return p
}

func NewPrometheusMetricFamily() (p *PrometheusMetricFamily) {
	return new(PrometheusMetricFamily)
}

func (p *PrometheusMetricFamily) AddMetric(toAdd *PrometheusMetric) (err error) {
	if toAdd == nil {
		return tracederrors.TracedErrorNil("toAdd")
	}

	p.metrics = append(p.metrics, toAdd)

	return nil
}

func (p *PrometheusMetricFamily) GetHelp() (help string, err error) {
	if p.help == "" {
		return "", tracederrors.TracedErrorf("help not set")
	}

	return p.help, nil
}

func (p *PrometheusMetricFamily) GetMetrics() (metrics []*PrometheusMetric, err error) {
	if p.metrics == nil {
		return nil, tracederrors.TracedErrorf("metrics not set")
	}

	if len(p.metrics) <= 0 {
		return nil, tracederrors.TracedErrorf("metrics has no elements")
	}

	return p.metrics, nil
}

func (p *PrometheusMetricFamily) GetName() (name string, err error) {
	if p.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return p.name, nil
}

func (p *PrometheusMetricFamily) GetValueAsFloat64() (value float64, err error) {
	metrics, err := p.GetMetrics()
	if err != nil {
		return -1, err
	}

	metricName, err := p.GetName()
	if err != nil {
		return -1, err
	}

	if len(metrics) != 1 {
		return -1, tracederrors.TracedErrorf(
			"Expected exactly 1 metric for '%s' to get the value from but got '%d'.",
			metricName,
			len(metrics),
		)
	}

	value, err = metrics[0].GetValueAsFloat64()
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (p *PrometheusMetricFamily) MustAddMetric(toAdd *PrometheusMetric) {
	err := p.AddMetric(toAdd)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetricFamily) MustGetHelp() (help string) {
	help, err := p.GetHelp()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return help
}

func (p *PrometheusMetricFamily) MustGetMetrics() (metrics []*PrometheusMetric) {
	metrics, err := p.GetMetrics()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return metrics
}

func (p *PrometheusMetricFamily) MustGetName() (name string) {
	name, err := p.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (p *PrometheusMetricFamily) MustGetValueAsFloat64() (value float64) {
	value, err := p.GetValueAsFloat64()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func (p *PrometheusMetricFamily) MustSetHelp(help string) {
	err := p.SetHelp(help)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetricFamily) MustSetMetrics(metrics []*PrometheusMetric) {
	err := p.SetMetrics(metrics)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetricFamily) MustSetName(name string) {
	err := p.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetricFamily) SetHelp(help string) (err error) {
	if help == "" {
		return tracederrors.TracedErrorf("help is empty string")
	}

	p.help = help

	return nil
}

func (p *PrometheusMetricFamily) SetMetrics(metrics []*PrometheusMetric) (err error) {
	if metrics == nil {
		return tracederrors.TracedErrorf("metrics is nil")
	}

	if len(metrics) <= 0 {
		return tracederrors.TracedErrorf("metrics has no elements")
	}

	p.metrics = metrics

	return nil
}

func (p *PrometheusMetricFamily) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	p.name = name

	return nil
}
