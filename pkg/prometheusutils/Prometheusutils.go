package prometheusutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils/httputilsparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func MustGetMetricValueFromMetricPage(ctx context.Context, url string, metricName string) (value float64) {
	value, err := GetMetricValueFromMetricPage(ctx, url, metricName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func GetMetricValueFromMetricPage(ctx context.Context, url string, metricName string) (value float64, err error) {
	if url == "" {
		return 0.0, tracederrors.TracedErrorEmptyString("url")
	}

	if metricName == "" {
		return 0.0, tracederrors.TracedErrorEmptyString("metricName")
	}

	logging.LogInfoByCtxf(ctx, "Collect metric value '%s' from %s started.", metricName, url)

	metrics, err := GetParsedMetricPage(ctx, url)
	if err != nil {
		return 0.0, err
	}

	value, err = metrics.GetMetricValueAsFloat64(metricName)
	if err != nil {
		return 0.0, err
	}

	logging.LogInfoByCtxf(ctx, "Collect metric value '%s' from %s finished.", metricName, url)

	return value, nil
}

// Use a get request to retrieve the exposed metrics from the given url.
// Then the already parsed PrometheusParsedMetrics are returned.
func GetParsedMetricPage(ctx context.Context, url string) (metrics *PrometheusParsedMetrics, err error) {
	if url == "" {
		return nil, tracederrors.TracedErrorEmptyString("url")
	}

	m, err := httputils.SendRequestAndGetBodyAsString(
		ctx,
		&httputilsparameteroptions.RequestOptions{
			Url:    url,
			Method: "GET",
		},
	)
	if err != nil {
		return nil, err
	}

	metrics, err = PrometheusExpositionFormatParser().ParseString(m)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
