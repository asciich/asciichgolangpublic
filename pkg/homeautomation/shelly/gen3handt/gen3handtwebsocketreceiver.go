package gen3handt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/homeautomation/shelly/messagestructure"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ---------------------------------------------------------------------------
// Metrics
// ---------------------------------------------------------------------------

type sensorRequestCounter struct {
	mu       sync.Mutex
	counters map[string]*int64
}

func newSensorRequestCounter() *sensorRequestCounter {
	return &sensorRequestCounter{
		counters: make(map[string]*int64),
	}
}

// Inc increments the counter for the given sensor name and returns the new value.
func (c *sensorRequestCounter) Inc(sensorName string) int64 {
	c.mu.Lock()
	if _, ok := c.counters[sensorName]; !ok {
		v := int64(0)
		c.counters[sensorName] = &v
	}
	ptr := c.counters[sensorName]
	c.mu.Unlock()
	return atomic.AddInt64(ptr, 1)
}

// Snapshot returns a copy of all counters at this point in time.
func (c *sensorRequestCounter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.counters))
	for k, v := range c.counters {
		out[k] = atomic.LoadInt64(v)
	}
	return out
}

// ---------------------------------------------------------------------------
// Value types
// ---------------------------------------------------------------------------

type ValueWithTimestamp struct {
	Value     float64    `json:"value"`
	TimeStamp *time.Time `json:"timestamp"`
}

type BatteryValue struct {
	Voltage   float64    `json:"value"`
	Perscent  float64    `json:"percent"`
	TimeStamp *time.Time `json:"timestamp"`
}

type SensorValues struct {
	TemperatureCelsius ValueWithTimestamp `json:"temperature_celsius"`
	HumidityPercent    ValueWithTimestamp `json:"humidity_percent"`
	Battery            BatteryValue       `json:"battery"`
}

// ---------------------------------------------------------------------------
// Receiver
// ---------------------------------------------------------------------------

type ShellyGen3HAndTWebsocketReceiver struct {
	Port         int
	SensorNames  []string
	SensorValues map[string]*SensorValues

	// metrics
	totalRequests     int64 // accessed via atomic
	sensorReqCounters *sensorRequestCounter
	sensorValuesMu    sync.RWMutex // protects SensorValues for concurrent reads from /metrics
}

// ---------------------------------------------------------------------------
// Sensor update helpers
// ---------------------------------------------------------------------------

func (s *ShellyGen3HAndTWebsocketReceiver) UpdateSensorBattery(ctx context.Context, sensorName string, voltage float64, percent float64) error {
	if sensorName == "" {
		return tracederrors.TracedErrorEmptyString("sensorName")
	}

	s.sensorValuesMu.Lock()
	defer s.sensorValuesMu.Unlock()

	if s.SensorValues == nil {
		s.SensorValues = map[string]*SensorValues{}
	}

	ts := time.Now()

	if s.SensorValues[sensorName] == nil {
		s.SensorValues[sensorName] = &SensorValues{}
	}

	s.SensorValues[sensorName].Battery = BatteryValue{
		Voltage:   voltage,
		Perscent:  percent,
		TimeStamp: &ts,
	}

	logging.LogInfoByCtxf(ctx, "Updated sensor '%s' battery to '%f' volts and '%f' percent.", sensorName, voltage, percent)

	return nil
}

func (s *ShellyGen3HAndTWebsocketReceiver) UpdateSensorTemperatureCelsius(ctx context.Context, sensorName string, value float64) error {
	if sensorName == "" {
		return tracederrors.TracedErrorEmptyString("sensorName")
	}

	s.sensorValuesMu.Lock()
	defer s.sensorValuesMu.Unlock()

	if s.SensorValues == nil {
		s.SensorValues = map[string]*SensorValues{}
	}

	if s.SensorValues[sensorName] == nil {
		s.SensorValues[sensorName] = &SensorValues{}
	}

	ts := time.Now()

	s.SensorValues[sensorName].TemperatureCelsius = ValueWithTimestamp{
		Value:     value,
		TimeStamp: &ts,
	}

	logging.LogInfoByCtxf(ctx, "Updated sensor '%s' temperature to '%f' degrees Celsius.", sensorName, value)

	return nil
}

func (s *ShellyGen3HAndTWebsocketReceiver) UpdateSensorHumidityPercent(ctx context.Context, sensorName string, value float64) error {
	if sensorName == "" {
		return tracederrors.TracedErrorEmptyString("sensorName")
	}

	s.sensorValuesMu.Lock()
	defer s.sensorValuesMu.Unlock()

	if s.SensorValues == nil {
		s.SensorValues = map[string]*SensorValues{}
	}

	if s.SensorValues[sensorName] == nil {
		s.SensorValues[sensorName] = &SensorValues{}
	}

	ts := time.Now()

	s.SensorValues[sensorName].HumidityPercent = ValueWithTimestamp{
		Value:     value,
		TimeStamp: &ts,
	}

	logging.LogInfoByCtxf(ctx, "Updated sensor '%s' humidity to '%f' percent.", sensorName, value)

	return nil
}

func (s *ShellyGen3HAndTWebsocketReceiver) UpdateSensorTemperatureCelsiusAndHumidityPercent(ctx context.Context, sensorName string, temperatureCelsius float64, humidityPercent float64) error {
	if sensorName == "" {
		return tracederrors.TracedErrorEmptyString("sensorName")
	}

	s.sensorValuesMu.Lock()
	defer s.sensorValuesMu.Unlock()

	if s.SensorValues == nil {
		s.SensorValues = map[string]*SensorValues{}
	}

	ts := time.Now()

	s.SensorValues[sensorName] = &SensorValues{
		TemperatureCelsius: ValueWithTimestamp{
			Value:     temperatureCelsius,
			TimeStamp: &ts,
		},
		HumidityPercent: ValueWithTimestamp{
			Value:     humidityPercent,
			TimeStamp: &ts,
		},
	}

	logging.LogInfoByCtxf(ctx, "Updated sensor '%s' temperature to '%f' degrees Celsius and humidity to '%f' percent.", sensorName, temperatureCelsius, humidityPercent)

	return nil
}

// ---------------------------------------------------------------------------
// WebSocket handler
// ---------------------------------------------------------------------------

func (s *ShellyGen3HAndTWebsocketReceiver) shellyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, sensorName string) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.LogGoError(tracederrors.TracedErrorf("Upgrade error: %w", err))
		return
	}
	defer conn.Close()

	logging.LogInfoByCtxf(ctx, "Shelly sensor '%s' connected from %s", sensorName, conn.RemoteAddr())

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "connection reset by peer") {
				logging.LogInfoByCtxf(ctx, "Sensor connection closed.")
			} else {
				logging.LogGoError(tracederrors.TracedErrorf("Connection closed: %w", err))
			}
			return
		}

		var msg messagestructure.ShellyMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("JSON parse error: %v — raw: %s", err, raw))
			continue
		}

		logging.LogInfoByCtxf(ctx, "[%s] method=%s src=%s", conn.RemoteAddr(), msg.Method, msg.Src)

		switch msg.Method {
		case "NotifyStatus":
			s.handleNotifyStatus(ctx, msg.Params, sensorName)

		case "NotifyEvent":
			s.handleNotifyEvent(ctx, msg.Params, sensorName)

		case "NotifyFullStatus":
			s.handleNotifyFullStatus(ctx, msg.Params, sensorName)

		default:
			logging.LogGoError(tracederrors.TracedErrorf("   Unknown method: %s — params: %s", msg.Method, msg.Params))
		}
	}
}

// ---------------------------------------------------------------------------
// Message handlers
// ---------------------------------------------------------------------------

func (s *ShellyGen3HAndTWebsocketReceiver) handleNotifyStatus(ctx context.Context, raw json.RawMessage, sensorName string) {
	var params messagestructure.NotifyStatusParams
	if err := json.Unmarshal(raw, &params); err != nil {
		logging.LogGoError(tracederrors.TracedErrorf("NotifyStatus parse error: %v", err))
		return
	}
	if params.Temperature != nil {
		err := s.UpdateSensorTemperatureCelsius(ctx, sensorName, params.Temperature.TC)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorTemperatureCelsius error: %v", err))
			return
		}
	}
	if params.Humidity != nil {
		err := s.UpdateSensorHumidityPercent(ctx, sensorName, params.Humidity.RH)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorHumidityPercent error: %v", err))
			return
		}
	}
	if params.DevicePower != nil {
		err := s.UpdateSensorBattery(ctx, sensorName, params.DevicePower.Battery.V, params.DevicePower.Battery.Percent)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorBattery error: %v", err))
			return
		}
	}
}

func (s *ShellyGen3HAndTWebsocketReceiver) handleNotifyEvent(ctx context.Context, raw json.RawMessage, sensorName string) {
	var params messagestructure.NotifyEventParams
	if err := json.Unmarshal(raw, &params); err != nil {
		logging.LogGoError(tracederrors.TracedErrorf("NotifyEvent parse error: %v", err))
		return
	}
	for _, e := range params.Events {
		logging.LogInfoByCtxf(ctx, "Event: component=%s event=%s", e.Component, e.Event)
	}
}

func (s *ShellyGen3HAndTWebsocketReceiver) handleNotifyFullStatus(ctx context.Context, raw json.RawMessage, sensorName string) {
	var params messagestructure.NotifyFullStatusParams
	if err := json.Unmarshal(raw, &params); err != nil {
		logging.LogGoError(tracederrors.TracedErrorf("NotifyFullStatus parse error: %v", err))
		return
	}

	if params.Temperature != nil {
		err := s.UpdateSensorTemperatureCelsius(ctx, sensorName, params.Temperature.TC)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorTemperatureCelsius error: %v", err))
			return
		}
	}
	if params.Humidity != nil {
		err := s.UpdateSensorHumidityPercent(ctx, sensorName, params.Humidity.RH)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorHumidityPercent error: %v", err))
			return
		}
	}
	if params.DevicePower != nil {
		err := s.UpdateSensorBattery(ctx, sensorName, params.DevicePower.Battery.V, params.DevicePower.Battery.Percent)
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("UpdateSensorBattery error: %v", err))
			return
		}
	}

	if params.Cloud != nil {
		logging.LogInfoByCtxf(ctx, "Cloud connected: %v", params.Cloud.Connected)
	}

	if params.MQTT != nil {
		logging.LogInfoByCtxf(ctx, "MQTT connected: %v", params.MQTT.Connected)
	}

	if params.WS != nil {
		logging.LogInfoByCtxf(ctx, "WebSocket connected: %v", params.WS.Connected)
	}

	if params.WiFi != nil {
		logging.LogInfoByCtxf(ctx, "WiFi: ssid=%s ip=%s rssi=%d", params.WiFi.SSID, params.WiFi.StaIP, params.WiFi.RSSI)
	}

	if params.Sys != nil {
		logging.LogInfoByCtxf(ctx, "Sys: mac=%s uptime=%.0fs wakeup_reason=%s/%s",
			params.Sys.MAC,
			params.Sys.Uptime,
			params.Sys.WakeupReason.Boot,
			params.Sys.WakeupReason.Cause,
		)
	}
}

func (s *ShellyGen3HAndTWebsocketReceiver) GetPort() (int, error) {
	if s.Port <= 0 {
		return 0, tracederrors.TracedErrorf("Port value is invalid: %d", s.Port)
	}

	return s.Port, nil
}

func (s *ShellyGen3HAndTWebsocketReceiver) GetSensorNames() ([]string, error) {
	if len(s.SensorNames) <= 0 {
		return nil, tracederrors.TracedErrorf("Empty sensor names")
	}

	return s.SensorNames, nil
}

// ---------------------------------------------------------------------------
// /metrics handler  (Prometheus text format 0.0.4, no extra dependencies)
// ---------------------------------------------------------------------------

// metricsHandler writes a Prometheus-compatible text exposition.
//
// Metrics exposed:
//
//	shelly_http_requests_total                          – all HTTP requests (counter)
//	shelly_sensor_requests_total{sensor="<name>"}       – requests per sensor endpoint (counter)
//	shelly_temperature_celsius{sensor="<name>"}         – latest temperature reading (gauge)
//	shelly_humidity_percent{sensor="<name>"}            – latest humidity reading (gauge)
//	shelly_battery_voltage{sensor="<name>"}             – latest battery voltage (gauge)
//	shelly_battery_percent{sensor="<name>"}             – latest battery percent (gauge)
func (s *ShellyGen3HAndTWebsocketReceiver) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	var b strings.Builder

	// ---- total requests ----
	total := atomic.LoadInt64(&s.totalRequests)
	b.WriteString("# HELP shelly_http_requests_total Total number of HTTP requests received (all endpoints).\n")
	b.WriteString("# TYPE shelly_http_requests_total counter\n")
	fmt.Fprintf(&b, "shelly_http_requests_total %d\n", total)

	// ---- per-sensor request counters ----
	b.WriteString("# HELP shelly_sensor_requests_total Total number of HTTP requests received per sensor endpoint.\n")
	b.WriteString("# TYPE shelly_sensor_requests_total counter\n")
	if s.sensorReqCounters != nil {
		for sensor, count := range s.sensorReqCounters.Snapshot() {
			fmt.Fprintf(&b, "shelly_sensor_requests_total{sensor=%q} %d\n", sensor, count)
		}
	}

	// ---- sensor values ----
	b.WriteString("# HELP shelly_temperature_celsius Latest temperature reading in degrees Celsius.\n")
	b.WriteString("# TYPE shelly_temperature_celsius gauge\n")
	b.WriteString("# HELP shelly_humidity_percent Latest relative humidity reading in percent.\n")
	b.WriteString("# TYPE shelly_humidity_percent gauge\n")
	b.WriteString("# HELP shelly_battery_voltage Latest battery voltage in Volts.\n")
	b.WriteString("# TYPE shelly_battery_voltage gauge\n")
	b.WriteString("# HELP shelly_battery_percent Latest battery charge level in percent.\n")
	b.WriteString("# TYPE shelly_battery_percent gauge\n")

	s.sensorValuesMu.RLock()
	defer s.sensorValuesMu.RUnlock()

	for sensor, sv := range s.SensorValues {
		if sv == nil {
			continue
		}
		if sv.TemperatureCelsius.TimeStamp != nil {
			fmt.Fprintf(&b, "shelly_temperature_celsius{sensor=%q} %g\n", sensor, sv.TemperatureCelsius.Value)
		}
		if sv.HumidityPercent.TimeStamp != nil {
			fmt.Fprintf(&b, "shelly_humidity_percent{sensor=%q} %g\n", sensor, sv.HumidityPercent.Value)
		}
		if sv.Battery.TimeStamp != nil {
			fmt.Fprintf(&b, "shelly_battery_voltage{sensor=%q} %g\n", sensor, sv.Battery.Voltage)
			fmt.Fprintf(&b, "shelly_battery_percent{sensor=%q} %g\n", sensor, sv.Battery.Perscent)
		}
	}

	fmt.Fprint(w, b.String())
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

func (s *ShellyGen3HAndTWebsocketReceiver) Run(ctx context.Context) error {
	port, err := s.GetPort()
	if err != nil {
		return err
	}

	// Initialise metrics helpers.
	s.sensorReqCounters = newSensorRequestCounter()

	logging.LogInfoByCtxf(ctx, "Run Shelly Gen3 H&T websocket receiver on port '%d' started.", port)

	addr := "0.0.0.0:" + strconv.Itoa(port)

	mux := http.NewServeMux()

	sensorNames, err := s.GetSensorNames()
	if err != nil {
		return err
	}

	for _, sensorName := range sensorNames {
		mux.HandleFunc("/"+sensorName, func(w http.ResponseWriter, r *http.Request) {
			s.sensorReqCounters.Inc(sensorName)
			ctxSensor := contextutils.WithLogLinePrefix(ctx, fmt.Sprintf("sensor: %s", sensorName))
			s.shellyHandler(ctxSensor, w, r, sensorName)
		})
		logging.LogInfoByCtxf(ctx, "Added endpoint '/%s' for sensor '%s'.", sensorName, sensorName)
	}

	mux.HandleFunc("/values.json", func(w http.ResponseWriter, r *http.Request) {
		type DataStruct struct {
			Timestamp    time.Time                `json:"time"`
			SensorValues map[string]*SensorValues `json:"sensor_values"`
		}

		s.sensorValuesMu.RLock()
		data := &DataStruct{
			Timestamp:    time.Now(),
			SensorValues: s.SensorValues,
		}
		s.sensorValuesMu.RUnlock()

		ret, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			logging.LogGoError(tracederrors.TracedErrorf("Failed to marshal response: %w", err))
			return
		}

		w.Write(ret)
	})

	mux.HandleFunc("/metrics", s.metricsHandler)
	logging.LogInfoByCtxf(ctx, "Added endpoint '/metrics' for Prometheus scraping.")

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logging.LogInfoByCtxf(ctx, "Received call to '%s', will be redirected to /values.json", r.URL.String())
		http.Redirect(w, r, "/values.json", http.StatusMovedPermanently)
	})

	// Middleware: count every request (including /metrics itself, invalid paths, etc.)
	var loggingMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt64(&s.totalRequests, 1)
			logging.LogInfoByCtxf(ctx, "Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}

	logging.LogInfoByCtxf(ctx, "Shelly WebSocket server listening on %s", addr)
	server := &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logging.LogInfoByCtxf(ctx, "Server shutdown error: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return tracederrors.TracedErrorf("Server error: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Run Shelly Gen3 H&T websocket receiver on port '%d' finished.", port)

	return nil
}
