package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type PrometheusMetric struct {
	value float64
	name  string
	help  string
}

func NewPrometheusMetric() (p *PrometheusMetric) {
	return new(PrometheusMetric)
}

func (p *PrometheusMetric) GetHelp() (help string, err error) {
	if p.help == "" {
		return "", errors.TracedErrorf("help not set")
	}

	return p.help, nil
}

func (p *PrometheusMetric) GetName() (name string, err error) {
	if p.name == "" {
		return "", errors.TracedErrorf("name not set")
	}

	return p.name, nil
}

func (p *PrometheusMetric) GetValue() (value float64, err error) {

	return p.value, nil
}

func (p *PrometheusMetric) GetValueAsFloat64() (value float64, err error) {
	return p.GetValue()
}

func (p *PrometheusMetric) MustGetHelp() (help string) {
	help, err := p.GetHelp()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return help
}

func (p *PrometheusMetric) MustGetName() (name string) {
	name, err := p.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (p *PrometheusMetric) MustGetValue() (value float64) {
	value, err := p.GetValue()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func (p *PrometheusMetric) MustGetValueAsFloat64() (value float64) {
	value, err := p.GetValueAsFloat64()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func (p *PrometheusMetric) MustSetHelp(help string) {
	err := p.SetHelp(help)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetric) MustSetName(name string) {
	err := p.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetric) MustSetValue(value float64) {
	err := p.SetValue(value)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetric) MustSetValueByFloat64(value float64) {
	err := p.SetValueByFloat64(value)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PrometheusMetric) SetHelp(help string) (err error) {
	if help == "" {
		return errors.TracedErrorf("help is empty string")
	}

	p.help = help

	return nil
}

func (p *PrometheusMetric) SetName(name string) (err error) {
	if name == "" {
		return errors.TracedErrorf("name is empty string")
	}

	p.name = name

	return nil
}

func (p *PrometheusMetric) SetValue(value float64) (err error) {
	p.value = value

	return nil
}

func (p *PrometheusMetric) SetValueByFloat64(value float64) (err error) {
	return p.SetValue(value)
}
