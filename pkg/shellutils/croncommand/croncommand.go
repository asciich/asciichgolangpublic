package croncommand

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	robfigcron "github.com/robfig/cron/v3"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Run the `command` as often as defined in `cron`
//
// `cron` can be:
//   - a standard crontab entry like: '*/5 * * * *' for every minute.
//   - a crontab entry with seconds precision, just prepend the seconds block additionally: `*/10 * * * * *` for every 10 seconds.
//
// This function blocks until there is an error executing 'command'. In case of the 'command' exits non successful the error is returned.
//
// If 'metricsPort' is greater than 0 prometheus metrics will be exposed on this port.
func RunCronCommand(ctx context.Context, name string, cron string, command []string, metricsPort int) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if cron == "" {
		return tracederrors.TracedErrorEmptyString("cron")
	}

	if command == nil {
		return tracederrors.TracedErrorNil("command")
	}

	joinedCommand, err := shelllinehandler.Join(command)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Going to periodically run cron command '%s' with command '%s' based on cron interval '%s' started.", name, joinedCommand, cron)

	// -------------------------------------------------------------------------
	// Prometheus metrics setup
	// -------------------------------------------------------------------------
	labels := prometheus.Labels{"name": name}

	// How many times the job has been executed
	jobRunsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cron_command_runs_total",
			Help: "Total number of times the cron job has been executed.",
		},
		[]string{"name"},
	)

	// Duration of the last job run in seconds
	jobLastDurationSeconds := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_command_last_duration_seconds",
			Help: "Duration in seconds of the last cron job execution. -1 if the job has never run.",
		},
		[]string{"name"},
	)

	// Unix timestamp of when the last job run finished
	jobLastFinishTimestamp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_command_last_finish_timestamp_seconds",
			Help: "Unix timestamp of when the last cron job execution finished. -1 if the job has never run.",
		},
		[]string{"name"},
	)

	// Unix timestamp of when the last job run started
	jobLastStartTimestamp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_command_last_start_timestamp_seconds",
			Help: "Unix timestamp of when the last cron job execution started. -1 if the job has never run.",
		},
		[]string{"name"},
	)

	// 1 if the job is currently running, 0 otherwise
	jobIsRunning := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_command_is_running",
			Help: "Indicates if the cron job is currently running (1 = running, 0 = idle).",
		},
		[]string{"name"},
	)

	// Unix timestamp of when RunCronCommand was called (i.e. the process start time)
	jobStartTimestamp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_command_start_timestamp_seconds",
			Help: "Unix timestamp of when the RunCronCommand function was called (process start time).",
		},
		[]string{"name"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		jobRunsTotal,
		jobLastDurationSeconds,
		jobLastFinishTimestamp,
		jobLastStartTimestamp,
		jobIsRunning,
		jobStartTimestamp,
	)

	// -------------------------------------------------------------------------
	// Initialize all metrics with useful default values so they are visible
	// in Prometheus/Grafana before the first run happens.
	// -1 is used for timestamps and durations to clearly indicate "never run".
	// -------------------------------------------------------------------------
	jobRunsTotal.With(labels).Add(0)            // initializes the counter at 0
	jobLastDurationSeconds.With(labels).Set(-1)  // -1 = never run
	jobLastFinishTimestamp.With(labels).Set(-1)  // -1 = never run
	jobLastStartTimestamp.With(labels).Set(-1)   // -1 = never run
	jobIsRunning.With(labels).Set(0)             // 0 = idle
	jobStartTimestamp.With(labels).Set(float64(time.Now().Unix()))

	// -------------------------------------------------------------------------
	// Start the metrics HTTP server if metricsPort > 0
	// -------------------------------------------------------------------------
	if metricsPort > 0 {
		metricsAddr := fmt.Sprintf(":%d", metricsPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

		metricsServer := &http.Server{
			Addr:    metricsAddr,
			Handler: mux,
		}

		go func() {
			logging.LogInfoByCtxf(ctx, "Exposing Prometheus metrics for cron job '%s' on http://0.0.0.0%s/metrics", name, metricsAddr)
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logging.LogErrorByCtxf(ctx, "Metrics server for cron job '%s' failed: %v", name, err)
			}
		}()

		// Shut down the metrics server when the context is done
		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = metricsServer.Shutdown(shutdownCtx)
		}()
	}

	// -------------------------------------------------------------------------
	// Cron scheduler setup
	// -------------------------------------------------------------------------

	// Detect if the cron string has 6 fields (with seconds) or 5 fields (standard)
	cronFields := strings.Fields(cron)

	var scheduler *robfigcron.Cron
	if len(cronFields) == 6 {
		scheduler = robfigcron.New(robfigcron.WithSeconds())
	} else {
		scheduler = robfigcron.New()
	}

	// jobErrCh is used to communicate a job execution error back to the blocking caller
	jobErrCh := make(chan error, 1)

	_, err = scheduler.AddFunc(cron, func() {
		tStart := time.Now()

		// Mark job as running and record last start timestamp
		jobIsRunning.With(labels).Set(1)
		jobLastStartTimestamp.With(labels).Set(float64(tStart.Unix()))

		logging.LogInfoByCtxf(ctx, "Cron job '%s' triggered, running command: '%s'.", name, joinedCommand)

		// Increment the run counter
		jobRunsTotal.With(labels).Inc()

		_, runErr := commandexecutorexec.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)

		// Always mark job as no longer running, regardless of success or failure
		jobIsRunning.With(labels).Set(0)

		if runErr != nil {
			jobErrCh <- tracederrors.TracedErrorf("cron job '%s' failed with command '%s': %w", name, joinedCommand, runErr)
			return
		}

		duration := time.Since(tStart)

		// Update duration and finish timestamp metrics
		jobLastDurationSeconds.With(labels).Set(duration.Seconds())
		jobLastFinishTimestamp.With(labels).Set(float64(time.Now().Unix()))

		logging.LogInfoByCtxf(ctx, "Cron job '%s' finished for command: '%s' took '%s'.", name, joinedCommand, duration)
	})
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", cron, err)
	}

	scheduler.Start()
	logging.LogInfoByCtxf(ctx, "Cron job '%s' scheduler started.", name)

	// Block until either:
	// 1. The job returns an error
	// 2. The context is cancelled
	select {
	case jobErr := <-jobErrCh:
		scheduler.Stop()
		return jobErr
	case <-ctx.Done():
		logging.LogInfoByCtxf(ctx, "Context cancelled, stopping cron job '%s'.", name)
		scheduler.Stop()
		return ctx.Err()
	}
}
