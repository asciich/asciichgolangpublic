package ollamaproxy

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// Name of the environment variable containing the shared API key clients have to send.
	EnvVarNameOllamaApiKey = "OLLAMA_API_KEY"

	// Default values used when no explicit option is given.
	DefaultOllamaProxyListenAddress = ":11435"
	DefaultOllamaProxyTargetUrl     = "http://localhost:11434"
	DefaultOllamaProxyMetricsPath   = "/metrics"

	// Prefix for all exposed prometheus metrics.
	OllamaProxyMetricsPrefix = "ollama_proxy_"
)

// OllamaProxyOptions defines the behaviour of an OllamaProxy.
type OllamaProxyOptions struct {
	// Address the proxy listens on. Defaults to DefaultOllamaProxyListenAddress.
	ListenAddress string

	// Base URL of the Ollama server to forward the requests to.
	// Defaults to DefaultOllamaProxyTargetUrl.
	TargetUrl string

	// Shared API key the clients have to send as "Authorization: Bearer <key>".
	// If empty the value of the environment variable OLLAMA_API_KEY is used.
	ApiKey string

	// Path the prometheus metrics are exposed on. Defaults to DefaultOllamaProxyMetricsPath.
	// The metrics endpoint is intentionally not protected by the API key so
	// prometheus can scrape it without additional configuration.
	MetricsPath string

	// Timeout for a single request forwarded to Ollama. Zero means no timeout
	// which is the recommended setting for long running prompts.
	UpstreamTimeout time.Duration

	// If set the "Authorization" header of the incoming request is forwarded to
	// Ollama as well. By default the header is stripped since the upstream
	// Ollama server does not need (and should not see) the shared key.
	ForwardAuthorizationHeader bool
}

// GetDeepCopy returns an independent copy of the given options.
func (o *OllamaProxyOptions) GetDeepCopy() *OllamaProxyOptions {
	if o == nil {
		return new(OllamaProxyOptions)
	}

	copy := *o

	return &copy
}

// GetApiKeyOrDefault returns the configured API key or the one defined by the
// environment variable OLLAMA_API_KEY.
func (o *OllamaProxyOptions) GetApiKeyOrDefault(ctx context.Context) (string, error) {
	if o == nil {
		return "", tracederrors.TracedErrorNil("o")
	}

	if o.ApiKey != "" {
		return o.ApiKey, nil
	}

	apiKey := strings.TrimSpace(os.Getenv(EnvVarNameOllamaApiKey))
	if apiKey == "" {
		return "", tracederrors.TracedErrorf(
			"Unable to get API key for the ollama proxy. Environment variable '%s' is empty or unset.",
			EnvVarNameOllamaApiKey,
		)
	}

	logging.LogInfoByCtxf(ctx, "Using API key from environment variable '%s'.", EnvVarNameOllamaApiKey)

	return apiKey, nil
}

// GetListenAddressOrDefault returns the address the proxy should listen on.
func (o *OllamaProxyOptions) GetListenAddressOrDefault(ctx context.Context) (string, error) {
	if o == nil {
		return "", tracederrors.TracedErrorNil("o")
	}

	if o.ListenAddress != "" {
		return o.ListenAddress, nil
	}

	logging.LogInfoByCtxf(ctx, "Using default listen address '%s'.", DefaultOllamaProxyListenAddress)

	return DefaultOllamaProxyListenAddress, nil
}

// GetTargetUrlOrDefault returns the parsed URL of the upstream ollama server.
func (o *OllamaProxyOptions) GetTargetUrlOrDefault(ctx context.Context) (*url.URL, error) {
	if o == nil {
		return nil, tracederrors.TracedErrorNil("o")
	}

	targetUrl := o.TargetUrl
	if targetUrl == "" {
		targetUrl = DefaultOllamaProxyTargetUrl
		logging.LogInfoByCtxf(ctx, "Using default ollama target url '%s'.", targetUrl)
	}

	parsed, err := url.Parse(targetUrl)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse ollama target url '%s': %w", targetUrl, err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, tracederrors.TracedErrorf("Ollama target url '%s' is not a valid absolute url.", targetUrl)
	}

	return parsed, nil
}

// GetMetricsPathOrDefault returns the path the prometheus metrics are exposed on.
func (o *OllamaProxyOptions) GetMetricsPathOrDefault() (string, error) {
	if o == nil {
		return "", tracederrors.TracedErrorNil("o")
	}

	if o.MetricsPath == "" {
		return DefaultOllamaProxyMetricsPath, nil
	}

	if !strings.HasPrefix(o.MetricsPath, "/") {
		return "/" + o.MetricsPath, nil
	}

	return o.MetricsPath, nil
}

// ollamaProxyMetrics bundles all prometheus metrics of the proxy.
// All metric names use the prefix "ollama_proxy_".
type ollamaProxyMetrics struct {
	registry *prometheus.Registry

	// Number of prompts currently waiting for the ollama server to become free.
	promtsWaiting prometheus.Gauge

	// 1 if a prompt is currently processed by ollama, 0 otherwise.
	promtRunning prometheus.Gauge

	// Counts all prompts performed since the proxy was started.
	promtCounter prometheus.Counter

	// Additional/ optional metrics which are handy in daily operations.
	promtErrorCounter prometheus.Counter
	unauthorizedCount prometheus.Counter
	promtDuration     prometheus.Histogram
}

func newOllamaProxyMetrics() (*ollamaProxyMetrics, error) {
	metrics := &ollamaProxyMetrics{
		registry: prometheus.NewRegistry(),

		promtsWaiting: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: OllamaProxyMetricsPrefix + "promts_waiting",
			Help: "Number of prompts currently waiting to be forwarded to the ollama server.",
		}),

		promtRunning: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: OllamaProxyMetricsPrefix + "promt_running",
			Help: "1 if there is an active promt currently processed, 0 otherwise.",
		}),

		promtCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: OllamaProxyMetricsPrefix + "promt_counter",
			Help: "Counts all promts performed by this proxy.",
		}),

		promtErrorCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: OllamaProxyMetricsPrefix + "promt_error_counter",
			Help: "Counts all promts which failed while being forwarded to the ollama server.",
		}),

		unauthorizedCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: OllamaProxyMetricsPrefix + "unauthorized_counter",
			Help: "Counts all requests rejected with http status 401.",
		}),

		promtDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    OllamaProxyMetricsPrefix + "promt_duration_seconds",
			Help:    "Duration of the promts forwarded to the ollama server.",
			Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60, 120, 300, 600},
		}),
	}

	collectors := []prometheus.Collector{
		metrics.promtsWaiting,
		metrics.promtRunning,
		metrics.promtCounter,
		metrics.promtErrorCounter,
		metrics.unauthorizedCount,
		metrics.promtDuration,
	}

	for _, collector := range collectors {
		if err := metrics.registry.Register(collector); err != nil {
			return nil, tracederrors.TracedErrorf("Unable to register prometheus collector: %w", err)
		}
	}

	// Initialize the gauges so they are present right after startup.
	metrics.promtsWaiting.Set(0)
	metrics.promtRunning.Set(0)

	return metrics, nil
}

// OllamaProxy is a reverse proxy in front of an ollama server which
//   - authenticates the clients using a shared API key (OpenAI style bearer token),
//   - ensures only one request at a time is forwarded to ollama,
//   - exposes prometheus metrics.
type OllamaProxy struct {
	options   *OllamaProxyOptions
	apiKey    string
	targetUrl *url.URL
	metrics   *ollamaProxyMetrics

	// ollamaMutex ensures only one request is forwarded to ollama at a time.
	ollamaMutex sync.Mutex

	reverseProxy *httputil.ReverseProxy
	handler      http.Handler
}

// NewOllamaProxy returns a ready to use OllamaProxy.
func NewOllamaProxy(ctx context.Context, options *OllamaProxyOptions) (*OllamaProxy, error) {
	optionsToUse := options.GetDeepCopy()

	apiKey, err := optionsToUse.GetApiKeyOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	targetUrl, err := optionsToUse.GetTargetUrlOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	metricsPath, err := optionsToUse.GetMetricsPathOrDefault()
	if err != nil {
		return nil, err
	}

	metrics, err := newOllamaProxyMetrics()
	if err != nil {
		return nil, err
	}

	proxy := &OllamaProxy{
		options:   optionsToUse,
		apiKey:    apiKey,
		targetUrl: targetUrl,
		metrics:   metrics,
	}

	proxy.reverseProxy = &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(targetUrl)
			request.Out.Host = targetUrl.Host

			// The shared key is only valid between client and proxy.
			if !optionsToUse.ForwardAuthorizationHeader {
				request.Out.Header.Del("Authorization")
			}

			request.SetXForwarded()
		},

		// Streaming responses (stream: true) have to be flushed immediately.
		FlushInterval: -1,

		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			metrics.promtErrorCounter.Inc()

			logging.LogErrorByCtxf(ctx, "Forwarding request '%s %s' to ollama failed: %v", request.Method, request.URL.Path, err)

			if errors.Is(err, context.Canceled) {
				// Client is gone, nothing to answer anymore.
				return
			}

			writeOpenAiStyleError(
				writer,
				http.StatusBadGateway,
				"Unable to reach the upstream ollama server.",
				"upstream_error",
				"upstream_unavailable",
			)
		},
	}

	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/", proxy.handleProxyRequest)

	proxy.handler = mux

	logging.LogInfoByCtxf(
		ctx,
		"Ollama proxy prepared. Target='%s', metrics path='%s'.",
		targetUrl.String(),
		metricsPath,
	)

	return proxy, nil
}

// GetHttpHandler returns the http.Handler of the proxy.
// Handy for tests (httptest.NewServer) or to embed the proxy into another server.
func (o *OllamaProxy) GetHttpHandler() (http.Handler, error) {
	if o == nil {
		return nil, tracederrors.TracedErrorNil("o")
	}

	if o.handler == nil {
		return nil, tracederrors.TracedError("Ollama proxy handler is not initialized. Use NewOllamaProxy to create the proxy.")
	}

	return o.handler, nil
}

// ServeHTTP implements http.Handler.
func (o *OllamaProxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if o == nil || o.handler == nil {
		writeOpenAiStyleError(
			writer,
			http.StatusInternalServerError,
			"Ollama proxy is not initialized.",
			"server_error",
			"not_initialized",
		)

		return
	}

	o.handler.ServeHTTP(writer, request)
}

// Run starts the proxy and blocks until the given context is cancelled or the
// http server fails.
func (o *OllamaProxy) Run(ctx context.Context) error {
	if o == nil {
		return tracederrors.TracedErrorNil("o")
	}

	listenAddress, err := o.options.GetListenAddressOrDefault(ctx)
	if err != nil {
		return err
	}

	handler, err := o.GetHttpHandler()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              listenAddress,
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second,
		// No write timeout on purpose: prompts may run for a very long time.
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to listen on '%s': %w", listenAddress, err)
	}

	logging.LogInfoByCtxf(
		ctx,
		"Ollama proxy started. Listening on '%s' and forwarding to '%s'.",
		listener.Addr().String(),
		o.targetUrl.String(),
	)

	shutdownFinished := make(chan struct{})

	go func() {
		defer close(shutdownFinished)

		<-ctx.Done()

		logging.LogInfoByCtxf(ctx, "Ollama proxy shutdown started.")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()

		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			logging.LogErrorByCtxf(ctx, "Ollama proxy shutdown failed: %v", shutdownErr)
		}
	}()

	err = server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return tracederrors.TracedErrorf("Ollama proxy failed: %w", err)
	}

	<-shutdownFinished

	logging.LogInfoByCtxf(ctx, "Ollama proxy finished.")

	return nil
}

// handleProxyRequest authenticates the client and forwards the request to the
// ollama server. Only one request is forwarded at a time.
func (o *OllamaProxy) handleProxyRequest(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if !o.isAuthorized(request) {
		o.metrics.unauthorizedCount.Inc()

		logging.LogWarnByCtxf(
			ctx,
			"Rejected unauthorized request '%s %s' from '%s'.",
			request.Method,
			request.URL.Path,
			request.RemoteAddr,
		)

		writeOpenAiStyleError(
			writer,
			http.StatusUnauthorized,
			"Incorrect API key provided. You can set it using the 'Authorization: Bearer <OLLAMA_API_KEY>' header.",
			"invalid_request_error",
			"invalid_api_key",
		)

		return
	}

	// From here on the request is going to be forwarded to ollama.
	// Since only one request at a time is allowed all others are waiting.
	o.metrics.promtsWaiting.Inc()

	o.ollamaMutex.Lock()

	o.metrics.promtsWaiting.Dec()
	o.metrics.promtRunning.Set(1)
	o.metrics.promtCounter.Inc()

	tStart := time.Now()

	defer func() {
		duration := time.Since(tStart)

		o.metrics.promtDuration.Observe(duration.Seconds())
		o.metrics.promtRunning.Set(0)

		o.ollamaMutex.Unlock()

		logging.LogInfoByCtxf(
			ctx,
			"Forwarded request '%s %s' to ollama. The request took '%v'.",
			request.Method,
			request.URL.Path,
			duration,
		)
	}()

	requestToUse := request

	if o.options.UpstreamTimeout > 0 {
		upstreamCtx, cancel := context.WithTimeout(ctx, o.options.UpstreamTimeout)
		defer cancel()

		requestToUse = request.WithContext(upstreamCtx)
	}

	o.reverseProxy.ServeHTTP(writer, requestToUse)
}

// isAuthorized returns true if and only if the request contains the correct
// shared API key as OpenAI style bearer token.
func (o *OllamaProxy) isAuthorized(request *http.Request) bool {
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return false
	}

	const bearerPrefix = "bearer "

	if len(authorizationHeader) <= len(bearerPrefix) {
		return false
	}

	if !strings.EqualFold(authorizationHeader[:len(bearerPrefix)], bearerPrefix) {
		return false
	}

	providedKey := strings.TrimSpace(authorizationHeader[len(bearerPrefix):])
	if providedKey == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(providedKey), []byte(o.apiKey)) == 1
}

// openAiError is the error format used by the OpenAI API.
type openAiError struct {
	Error openAiErrorDetails `json:"error"`
}

type openAiErrorDetails struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    *string `json:"code"`
}

// writeOpenAiStyleError writes an error response in the format used by the OpenAI API.
func writeOpenAiStyleError(writer http.ResponseWriter, statusCode int, message string, errorType string, code string) {
	writer.Header().Set("Content-Type", "application/json")

	if statusCode == http.StatusUnauthorized {
		// OpenAI compatible clients expect this header on 401.
		writer.Header().Set("WWW-Authenticate", `Bearer realm="ollama-proxy"`)
	}

	writer.WriteHeader(statusCode)

	errorCode := code

	payload := openAiError{
		Error: openAiErrorDetails{
			Message: message,
			Type:    errorType,
			Param:   nil,
			Code:    &errorCode,
		},
	}

	// Nothing useful can be done if the connection is already gone.
	_ = json.NewEncoder(writer).Encode(payload)
}

// RunOllamaProxy is a convenience function to start a proxy with the given options.
func RunOllamaProxy(ctx context.Context, options *OllamaProxyOptions) error {
	proxy, err := NewOllamaProxy(ctx, options)
	if err != nil {
		return err
	}

	return proxy.Run(ctx)
}
