package ollamaproxy_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils/ollamaproxy"
)

const (
	testApiKey      = "test-api-key-1234567890"
	testEnvVarName  = "OLLAMA_API_KEY"
	testUpstreamMsg = `{"model":"test-model","response":"hello from ollama","done":true}`
)

// ---------------------------------------------------------------------------
// test helpers
// ---------------------------------------------------------------------------

// testContext returns the context used in all tests.
func testContext() context.Context {
	return context.Background()
}

// upstreamRecorder records what the fake ollama server has seen.
type upstreamRecorder struct {
	mutex sync.Mutex

	// Number of requests received.
	requestCount int

	// Concurrency tracking.
	currentConcurrency int64
	maxConcurrency     int64

	// Headers and paths of the last request.
	lastAuthorizationHeader string
	lastPath                string
	lastBody                string
}

func (u *upstreamRecorder) enter(request *http.Request) {
	bodyBytes, _ := io.ReadAll(request.Body)

	current := atomic.AddInt64(&u.currentConcurrency, 1)

	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.requestCount++
	u.lastAuthorizationHeader = request.Header.Get("Authorization")
	u.lastPath = request.URL.Path
	u.lastBody = string(bodyBytes)

	if current > u.maxConcurrency {
		u.maxConcurrency = current
	}
}

func (u *upstreamRecorder) leave() {
	atomic.AddInt64(&u.currentConcurrency, -1)
}

func (u *upstreamRecorder) GetRequestCount() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.requestCount
}

func (u *upstreamRecorder) GetMaxConcurrency() int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.maxConcurrency
}

func (u *upstreamRecorder) GetLastAuthorizationHeader() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastAuthorizationHeader
}

func (u *upstreamRecorder) GetLastPath() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastPath
}

func (u *upstreamRecorder) GetLastBody() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return u.lastBody
}

// mustStartUpstream starts a fake ollama server using the given handler.
func mustStartUpstream(t *testing.T, recorder *upstreamRecorder, handler func(writer http.ResponseWriter, request *http.Request)) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		recorder.enter(request)
		defer recorder.leave()

		handler(writer, request)
	}))

	t.Cleanup(server.Close)

	return server
}

// mustStartDefaultUpstream starts a fake ollama server answering a static response.
func mustStartDefaultUpstream(t *testing.T, recorder *upstreamRecorder) *httptest.Server {
	t.Helper()

	return mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write([]byte(testUpstreamMsg))
	})
}

// mustStartProxy creates the proxy and exposes it using a httptest server.
func mustStartProxy(t *testing.T, options *ollamaproxy.OllamaProxyOptions) *httptest.Server {
	t.Helper()

	proxy, err := ollamaproxy.NewOllamaProxy(testContext(), options)
	require.NoError(t, err)
	require.NotNil(t, proxy)

	handler, err := proxy.GetHttpHandler()
	require.NoError(t, err)
	require.NotNil(t, handler)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server
}

// doRequest sends a request to the proxy using the given authorization header.
// An empty authorizationHeader means "do not set the header at all".
func doRequest(t *testing.T, proxyUrl string, path string, authorizationHeader string, body string) *http.Response {
	t.Helper()

	request, err := http.NewRequestWithContext(
		testContext(),
		http.MethodPost,
		proxyUrl+path,
		strings.NewReader(body),
	)
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	if authorizationHeader != "" {
		request.Header.Set("Authorization", authorizationHeader)
	}

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = response.Body.Close()
	})

	return response
}

func readBodyAsString(t *testing.T, response *http.Response) string {
	t.Helper()

	bodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(bodyBytes)
}

// openAiErrorResponse is the error format used by the OpenAI API.
type openAiErrorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// requireOpenAiErrorFormat validates that the body is a valid OpenAI style error.
func requireOpenAiErrorFormat(t *testing.T, body string) openAiErrorResponse {
	t.Helper()

	// The body must be a JSON object containing exactly the key "error".
	var raw map[string]json.RawMessage

	err := json.Unmarshal([]byte(body), &raw)
	require.NoError(t, err, "response body is not valid json: %s", body)
	require.Contains(t, raw, "error", "response body does not contain the 'error' key: %s", body)

	var parsed openAiErrorResponse

	err = json.Unmarshal([]byte(body), &parsed)
	require.NoError(t, err)

	require.NotEmpty(t, parsed.Error.Message, "the error message must not be empty: %s", body)
	require.NotEmpty(t, parsed.Error.Type, "the error type must not be empty: %s", body)

	return parsed
}

// getMetricsPage scrapes the /metrics page of the proxy.
func getMetricsPage(t *testing.T, proxyUrl string) string {
	t.Helper()

	request, err := http.NewRequestWithContext(testContext(), http.MethodGet, proxyUrl+"/metrics", nil)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer func() {
		_ = response.Body.Close()
	}()

	require.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(bodyBytes)
}

// getMetricValue returns the value of the given metric from the /metrics page.
func getMetricValue(t *testing.T, proxyUrl string, metricName string) float64 {
	t.Helper()

	value, found := findMetricValue(t, getMetricsPage(t, proxyUrl), metricName)
	require.True(t, found, "metric '%s' not found on the metrics page", metricName)

	return value
}

// findMetricValue parses the value of a metric out of a prometheus text page.
func findMetricValue(t *testing.T, metricsPage string, metricName string) (float64, bool) {
	t.Helper()

	expression := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(metricName) + `(?:\{[^}]*\})?\s+(\S+)\s*$`)

	match := expression.FindStringSubmatch(metricsPage)
	if match == nil {
		return 0, false
	}

	value, err := strconv.ParseFloat(match[1], 64)
	require.NoError(t, err)

	return value, true
}

// getWaitingPromptsMetricName returns the name of the "currently waiting prompts" metric.
//
// The specification does not fix the exact name of this metric (only its prefix),
// so both commonly used spellings are accepted.
func getWaitingPromptsMetricName(t *testing.T, proxyUrl string) string {
	t.Helper()

	metricsPage := getMetricsPage(t, proxyUrl)

	candidates := []string{
		"ollama_proxy_promts_waiting",
		"ollama_proxy_prompts_waiting",
		"ollama_proxy_promts_currently_waiting",
		"ollama_proxy_waiting_promts",
	}

	for _, candidate := range candidates {
		if _, found := findMetricValue(t, metricsPage, candidate); found {
			return candidate
		}
	}

	require.Failf(
		t,
		"missing metric",
		"None of the expected 'waiting prompts' metrics %v was found on the metrics page:\n%s",
		candidates,
		metricsPage,
	)

	return ""
}

// waitForMetricValue polls the metrics page until the metric reaches the expected value.
func waitForMetricValue(t *testing.T, proxyUrl string, metricName string, expectedValue float64, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	lastValue := float64(-1)

	for time.Now().Before(deadline) {
		lastValue = getMetricValue(t, proxyUrl, metricName)

		if lastValue == expectedValue {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	require.Failf(
		t,
		"metric did not reach the expected value",
		"metric '%s' is '%v' but '%v' was expected within '%v'.",
		metricName,
		lastValue,
		expectedValue,
		timeout,
	)
}

// ---------------------------------------------------------------------------
// construction
// ---------------------------------------------------------------------------

func TestNewOllamaProxy_UsesApiKeyFromEnvVar(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	// The key from the env var has to be accepted.
	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{"prompt":"hi"}`)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestNewOllamaProxy_FailsIfApiKeyEnvVarIsMissing(t *testing.T) {
	// Make sure the env var is unset for this test.
	original, wasSet := os.LookupEnv(testEnvVarName)

	require.NoError(t, os.Unsetenv(testEnvVarName))

	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv(testEnvVarName, original)

			return
		}

		_ = os.Unsetenv(testEnvVarName)
	})

	proxy, err := ollamaproxy.NewOllamaProxy(testContext(), &ollamaproxy.OllamaProxyOptions{
		TargetUrl: "http://localhost:11434",
	})

	require.Error(t, err)
	require.Nil(t, proxy)
	require.Contains(t, err.Error(), testEnvVarName)
}

func TestNewOllamaProxy_FailsIfApiKeyEnvVarIsEmpty(t *testing.T) {
	t.Setenv(testEnvVarName, "")

	proxy, err := ollamaproxy.NewOllamaProxy(testContext(), &ollamaproxy.OllamaProxyOptions{
		TargetUrl: "http://localhost:11434",
	})

	require.Error(t, err)
	require.Nil(t, proxy)
}

func TestNewOllamaProxy_ExplicitApiKeyOverridesEnvVar(t *testing.T) {
	t.Setenv(testEnvVarName, "key-from-env")

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
		ApiKey:    "key-from-options",
	})

	responseFromEnvKey := doRequest(t, proxyServer.URL, "/api/generate", "Bearer key-from-env", `{}`)
	require.Equal(t, http.StatusUnauthorized, responseFromEnvKey.StatusCode)

	responseFromOptionsKey := doRequest(t, proxyServer.URL, "/api/generate", "Bearer key-from-options", `{}`)
	require.Equal(t, http.StatusOK, responseFromOptionsKey.StatusCode)
}

func TestNewOllamaProxy_FailsOnInvalidTargetUrl(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	invalidUrls := []string{
		"localhost:11434",
		"not a url",
		"/only/a/path",
	}

	for _, invalidUrl := range invalidUrls {
		t.Run(invalidUrl, func(t *testing.T) {
			proxy, err := ollamaproxy.NewOllamaProxy(testContext(), &ollamaproxy.OllamaProxyOptions{
				TargetUrl: invalidUrl,
			})

			require.Error(t, err)
			require.Nil(t, proxy)
		})
	}
}

// ---------------------------------------------------------------------------
// authentication
// ---------------------------------------------------------------------------

func TestOllamaProxy_UnauthorizedRequestsAreRejected(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	testCases := []struct {
		name                string
		authorizationHeader string
	}{
		{"missing header", ""},
		{"empty bearer", "Bearer "},
		{"bearer without value", "Bearer"},
		{"wrong key", "Bearer wrong-key"},
		{"key with typo", "Bearer " + testApiKey + "x"},
		{"key without bearer prefix", testApiKey},
		{"basic auth instead of bearer", "Basic " + testApiKey},
		{"token prefix instead of bearer", "Token " + testApiKey},
		{"only whitespace", "   "},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := &upstreamRecorder{}
			upstream := mustStartDefaultUpstream(t, recorder)

			proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
				TargetUrl: upstream.URL,
			})

			response := doRequest(t, proxyServer.URL, "/api/generate", testCase.authorizationHeader, `{"prompt":"hi"}`)

			require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			require.Contains(t, response.Header.Get("Content-Type"), "application/json")

			// The request must never reach ollama.
			require.Equal(t, 0, recorder.GetRequestCount())

			// The answer has to use the OpenAI error format.
			parsedError := requireOpenAiErrorFormat(t, readBodyAsString(t, response))
			require.NotEmpty(t, parsedError.Error.Message)
		})
	}
}

func TestOllamaProxy_UnauthorizedErrorUsesOpenAiFormat(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	response := doRequest(t, proxyServer.URL, "/v1/chat/completions", "Bearer nope", `{}`)
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	parsedError := requireOpenAiErrorFormat(t, readBodyAsString(t, response))

	require.Equal(t, "invalid_request_error", parsedError.Error.Type)
	require.NotNil(t, parsedError.Error.Code)
	require.Equal(t, "invalid_api_key", *parsedError.Error.Code)
	require.Nil(t, parsedError.Error.Param)
}

func TestOllamaProxy_BearerPrefixIsCaseInsensitive(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	prefixes := []string{"Bearer ", "bearer ", "BEARER ", "BeArEr "}

	for _, prefix := range prefixes {
		t.Run(strings.TrimSpace(prefix), func(t *testing.T) {
			response := doRequest(t, proxyServer.URL, "/api/generate", prefix+testApiKey, `{}`)
			require.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}

func TestOllamaProxy_AuthorizedRequestIsForwardedToOllama(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	requestBody := `{"model":"test-model","prompt":"hi","stream":false}`

	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, requestBody)

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, testUpstreamMsg, readBodyAsString(t, response))

	require.Equal(t, 1, recorder.GetRequestCount())
	require.Equal(t, "/api/generate", recorder.GetLastPath())
	require.Equal(t, requestBody, recorder.GetLastBody())
}

func TestOllamaProxy_AuthorizationHeaderIsStrippedByDefault(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)
	require.Equal(t, http.StatusOK, response.StatusCode)

	// The shared key is only valid between client and proxy.
	require.Empty(t, recorder.GetLastAuthorizationHeader())
}

func TestOllamaProxy_AuthorizationHeaderIsForwardedIfEnabled(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl:                  upstream.URL,
		ForwardAuthorizationHeader: true,
	})

	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)
	require.Equal(t, http.StatusOK, response.StatusCode)

	require.Equal(t, "Bearer "+testApiKey, recorder.GetLastAuthorizationHeader())
}

// ---------------------------------------------------------------------------
// only one request to ollama at a time
// ---------------------------------------------------------------------------

func TestOllamaProxy_OnlyOneRequestToOllamaAtATime(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	const numberOfRequests = 15

	recorder := &upstreamRecorder{}

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		// Sleep long enough so overlapping requests would be detected.
		time.Sleep(20 * time.Millisecond)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(testUpstreamMsg))
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	var waitGroup sync.WaitGroup

	statusCodes := make([]int, numberOfRequests)

	for i := 0; i < numberOfRequests; i++ {
		waitGroup.Add(1)

		go func(index int) {
			defer waitGroup.Done()

			request, err := http.NewRequestWithContext(
				testContext(),
				http.MethodPost,
				proxyServer.URL+"/api/generate",
				strings.NewReader(`{"prompt":"hi"}`),
			)
			if err != nil {
				statusCodes[index] = -1

				return
			}

			request.Header.Set("Authorization", "Bearer "+testApiKey)

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				statusCodes[index] = -1

				return
			}

			defer func() {
				_ = response.Body.Close()
			}()

			_, _ = io.Copy(io.Discard, response.Body)

			statusCodes[index] = response.StatusCode
		}(i)
	}

	waitGroup.Wait()

	for index, statusCode := range statusCodes {
		require.Equalf(t, http.StatusOK, statusCode, "request '%d' failed", index)
	}

	require.Equal(t, numberOfRequests, recorder.GetRequestCount())
	require.EqualValues(t, 1, recorder.GetMaxConcurrency(), "more than one request reached ollama at the same time")
}

func TestOllamaProxy_UnauthorizedRequestIsNotBlockedByRunningPromt(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}

	releaseUpstream := make(chan struct{})
	upstreamEntered := make(chan struct{}, 1)

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		upstreamEntered <- struct{}{}

		<-releaseUpstream

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(testUpstreamMsg))
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	go func() {
		response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)
		_, _ = io.Copy(io.Discard, response.Body)
	}()

	// Wait until the long running prompt occupies the single ollama slot.
	select {
	case <-upstreamEntered:
	case <-time.After(5 * time.Second):
		require.Fail(t, "upstream was never reached")
	}

	// Authentication must not wait for the ollama slot.
	done := make(chan int, 1)

	go func() {
		response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer wrong-key", `{}`)
		done <- response.StatusCode
	}()

	select {
	case statusCode := <-done:
		require.Equal(t, http.StatusUnauthorized, statusCode)
	case <-time.After(2 * time.Second):
		require.Fail(t, "an unauthorized request was blocked by a running promt")
	}

	close(releaseUpstream)
}

// ---------------------------------------------------------------------------
// metrics
// ---------------------------------------------------------------------------

func TestOllamaProxy_MetricsEndpointIsExposed(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	metricsPage := getMetricsPage(t, proxyServer.URL)

	require.NotEmpty(t, metricsPage)

	// The required metrics have to be present right after startup.
	require.Contains(t, metricsPage, "ollama_proxy_promt_running")
	require.Contains(t, metricsPage, "ollama_proxy_promt_counter")

	waitingMetricName := getWaitingPromptsMetricName(t, proxyServer.URL)
	require.Contains(t, metricsPage, waitingMetricName)

	// Scraping the metrics page must not be forwarded to ollama.
	require.Equal(t, 0, recorder.GetRequestCount())
}

func TestOllamaProxy_AllMetricsUseTheRequiredPrefix(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	metricsPage := getMetricsPage(t, proxyServer.URL)

	for _, line := range strings.Split(metricsPage, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		metricName := ""

		switch {
		case strings.HasPrefix(line, "# HELP "):
			metricName = strings.Fields(strings.TrimPrefix(line, "# HELP "))[0]
		case strings.HasPrefix(line, "# TYPE "):
			metricName = strings.Fields(strings.TrimPrefix(line, "# TYPE "))[0]
		case strings.HasPrefix(line, "#"):
			continue
		default:
			metricName = strings.Fields(line)[0]
			if index := strings.Index(metricName, "{"); index >= 0 {
				metricName = metricName[:index]
			}
		}

		require.Truef(
			t,
			strings.HasPrefix(metricName, "ollama_proxy_"),
			"metric '%s' does not use the required prefix 'ollama_proxy_'",
			metricName,
		)
	}
}

func TestOllamaProxy_MetricsInitialValues(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	waitingMetricName := getWaitingPromptsMetricName(t, proxyServer.URL)

	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, waitingMetricName))
	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running"))
	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"))
}

func TestOllamaProxy_PromtCounterCountsAllPromts(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	const numberOfPromts = 5

	for i := 0; i < numberOfPromts; i++ {
		response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{"prompt":"hi"}`)
		require.Equal(t, http.StatusOK, response.StatusCode)

		_, _ = io.Copy(io.Discard, response.Body)

		require.EqualValues(
			t,
			i+1,
			getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"),
			"the promt counter was not incremented",
		)
	}

	require.EqualValues(t, numberOfPromts, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"))
}

func TestOllamaProxy_PromtCounterIgnoresUnauthorizedRequests(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	for i := 0; i < 3; i++ {
		response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer wrong-key", `{}`)
		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	}

	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"))
	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running"))
}

func TestOllamaProxy_PromtRunningIsOneWhileAPromtIsProcessed(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}

	releaseUpstream := make(chan struct{})
	upstreamEntered := make(chan struct{}, 1)

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		upstreamEntered <- struct{}{}

		<-releaseUpstream

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(testUpstreamMsg))
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running"))

	requestFinished := make(chan struct{})

	go func() {
		defer close(requestFinished)

		response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)
		_, _ = io.Copy(io.Discard, response.Body)
	}()

	select {
	case <-upstreamEntered:
	case <-time.After(5 * time.Second):
		require.Fail(t, "upstream was never reached")
	}

	// While the promt is processed the gauge has to be 1.
	require.EqualValues(t, 1, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running"))
	require.EqualValues(t, 1, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"))

	close(releaseUpstream)

	<-requestFinished

	// After the promt finished the gauge has to be 0 again.
	waitForMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running", 0, 5*time.Second)
}

func TestOllamaProxy_WaitingPromtsAreCounted(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	const numberOfRequests = 4

	recorder := &upstreamRecorder{}

	releaseUpstream := make(chan struct{})
	upstreamEntered := make(chan struct{}, numberOfRequests)

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		upstreamEntered <- struct{}{}

		<-releaseUpstream

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(testUpstreamMsg))
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	waitingMetricName := getWaitingPromptsMetricName(t, proxyServer.URL)

	require.EqualValues(t, 0, getMetricValue(t, proxyServer.URL, waitingMetricName))

	var waitGroup sync.WaitGroup

	for i := 0; i < numberOfRequests; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			request, err := http.NewRequestWithContext(
				testContext(),
				http.MethodPost,
				proxyServer.URL+"/api/generate",
				strings.NewReader(`{"prompt":"hi"}`),
			)
			if err != nil {
				return
			}

			request.Header.Set("Authorization", "Bearer "+testApiKey)

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return
			}

			defer func() {
				_ = response.Body.Close()
			}()

			_, _ = io.Copy(io.Discard, response.Body)
		}()
	}

	// One request is running, all others are waiting.
	waitForMetricValue(t, proxyServer.URL, waitingMetricName, numberOfRequests-1, 10*time.Second)
	require.EqualValues(t, 1, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running"))

	close(releaseUpstream)

	waitGroup.Wait()

	// Everything processed: no waiting prompts, nothing running, all counted.
	waitForMetricValue(t, proxyServer.URL, waitingMetricName, 0, 10*time.Second)
	waitForMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running", 0, 10*time.Second)
	require.EqualValues(t, numberOfRequests, getMetricValue(t, proxyServer.URL, "ollama_proxy_promt_counter"))
}

func TestOllamaProxy_MetricsPathIsConfigurable(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl:   upstream.URL,
		MetricsPath: "/internal/metrics",
	})

	request, err := http.NewRequestWithContext(testContext(), http.MethodGet, proxyServer.URL+"/internal/metrics", nil)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer func() {
		_ = response.Body.Close()
	}()

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, readBodyAsString(t, response), "ollama_proxy_promt_counter")
}

func TestOllamaProxy_MultipleProxiesCanBeCreated(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	// Each proxy has to use its own registry, otherwise this panics or errors
	// because of duplicate metric registration.
	for i := 0; i < 3; i++ {
		proxy, err := ollamaproxy.NewOllamaProxy(testContext(), &ollamaproxy.OllamaProxyOptions{
			TargetUrl: "http://localhost:11434",
		})

		require.NoError(t, err)
		require.NotNil(t, proxy)
	}
}

// ---------------------------------------------------------------------------
// upstream behaviour
// ---------------------------------------------------------------------------

func TestOllamaProxy_UnreachableUpstreamIsAnsweredWithOpenAiError(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	// Start and immediately close an upstream to get a port nobody listens on.
	recorder := &upstreamRecorder{}
	deadUpstream := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		recorder.enter(request)
	}))

	deadUpstreamUrl := deadUpstream.URL
	deadUpstream.Close()

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: deadUpstreamUrl,
	})

	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)

	require.Equal(t, http.StatusBadGateway, response.StatusCode)
	requireOpenAiErrorFormat(t, readBodyAsString(t, response))

	// Even a failed promt must release the single ollama slot.
	waitForMetricValue(t, proxyServer.URL, "ollama_proxy_promt_running", 0, 5*time.Second)
}

func TestOllamaProxy_UpstreamStatusCodeAndBodyArePassedThrough(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("X-Custom-Upstream-Header", "ollama")
		writer.WriteHeader(http.StatusNotFound)

		_, _ = writer.Write([]byte(`{"error":"model not found"}`))
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	response := doRequest(t, proxyServer.URL, "/api/generate", "Bearer "+testApiKey, `{}`)

	require.Equal(t, http.StatusNotFound, response.StatusCode)
	require.Equal(t, "ollama", response.Header.Get("X-Custom-Upstream-Header"))
	require.JSONEq(t, `{"error":"model not found"}`, readBodyAsString(t, response))
}

func TestOllamaProxy_StreamingResponseIsNotBuffered(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}

	secondChunkReleased := make(chan struct{})

	upstream := mustStartUpstream(t, recorder, func(writer http.ResponseWriter, request *http.Request) {
		flusher, ok := writer.(http.Flusher)
		require.True(t, ok)

		writer.Header().Set("Content-Type", "application/x-ndjson")
		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write([]byte(`{"response":"first","done":false}` + "\n"))
		flusher.Flush()

		<-secondChunkReleased

		_, _ = writer.Write([]byte(`{"response":"second","done":true}` + "\n"))
		flusher.Flush()
	})

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	request, err := http.NewRequestWithContext(
		testContext(),
		http.MethodPost,
		proxyServer.URL+"/api/generate",
		strings.NewReader(`{"prompt":"hi","stream":true}`),
	)
	require.NoError(t, err)

	request.Header.Set("Authorization", "Bearer "+testApiKey)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer func() {
		_ = response.Body.Close()
	}()

	require.Equal(t, http.StatusOK, response.StatusCode)

	reader := bufio.NewReader(response.Body)

	firstLineChannel := make(chan string, 1)

	go func() {
		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			firstLineChannel <- ""

			return
		}

		firstLineChannel <- line
	}()

	// The first chunk must arrive before the upstream finished the response.
	select {
	case firstLine := <-firstLineChannel:
		require.Contains(t, firstLine, "first")
	case <-time.After(3 * time.Second):
		close(secondChunkReleased)
		require.Fail(t, "the streamed chunk was buffered by the proxy")
	}

	close(secondChunkReleased)

	remaining, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Contains(t, string(remaining), "second")
}

func TestOllamaProxy_AllPathsAndMethodsRequireAuthentication(t *testing.T) {
	t.Setenv(testEnvVarName, testApiKey)

	recorder := &upstreamRecorder{}
	upstream := mustStartDefaultUpstream(t, recorder)

	proxyServer := mustStartProxy(t, &ollamaproxy.OllamaProxyOptions{
		TargetUrl: upstream.URL,
	})

	paths := []string{
		"/",
		"/api/generate",
		"/api/chat",
		"/api/tags",
		"/api/pull",
		"/v1/chat/completions",
		"/v1/models",
	}

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodHead,
	}

	for _, path := range paths {
		for _, method := range methods {
			t.Run(fmt.Sprintf("%s %s", method, path), func(t *testing.T) {
				request, err := http.NewRequestWithContext(testContext(), method, proxyServer.URL+path, nil)
				require.NoError(t, err)

				response, err := http.DefaultClient.Do(request)
				require.NoError(t, err)

				defer func() {
					_ = response.Body.Close()
				}()

				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			})
		}
	}

	require.Equal(t, 0, recorder.GetRequestCount())
}

// ---------------------------------------------------------------------------
// handler sanity checks
// ---------------------------------------------------------------------------

func TestOllamaProxy_ServeHttpOfNilProxyDoesNotPanic(t *testing.T) {
	var proxy *ollamaproxy.OllamaProxy

	recorderWriter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/tags", nil)

	require.NotPanics(t, func() {
		proxy.ServeHTTP(recorderWriter, request)
	})

	require.Equal(t, http.StatusInternalServerError, recorderWriter.Code)
}

func TestOllamaProxy_GetHttpHandlerOfNilProxyReturnsError(t *testing.T) {
	var proxy *ollamaproxy.OllamaProxy

	handler, err := proxy.GetHttpHandler()

	require.Error(t, err)
	require.Nil(t, handler)
}
