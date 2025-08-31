package gitlabmetricscmd

import (
	"context"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewExposeMetricsCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "expose-prometheus-metrics",
		Short: "Expose selected gitlab metrics using prometheus exporter.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify the Gitlab --url.")
			}

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if port <= 0 {
				logging.LogFatalf("Pleases specify a valid --port. '%d' is not valid.", port)
			}

			group, err := cmd.Flags().GetString("group")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			cliExposeMetrics(ctx, url, port, group)
		},
	}

	cmd.PersistentFlags().String("url", "", "Url of the Gitlab instance to expose.")
	cmd.PersistentFlags().Int("port", 9123, "Port to expose the metrics.")
	cmd.PersistentFlags().String("group", "", "Path to group to collect project metrics. (optional)")

	return cmd
}

func cliExposeMetrics(ctx context.Context, url string, port int, group string) {
	logging.LogInfoByCtxf(ctx, "Going to export metrics of Gitlab instance '%s'", url)

	gitlab := mustutils.Must(asciichgolangpublic.GetGitlabByFQDN(url))

	accessToken := os.Getenv("GITLAB_ACCESS_TOKEN")
	if accessToken == "" {
		logging.LogFatal("GITLAB_ACCESS_TOKEN is not set as environment variable.")
	} else {
		logging.LogInfoByCtx(ctx, "Gitlab token read from env var GITLAB_ACCESS_TOKEN.")
	}

	if group == "" {
		logging.LogInfoByCtx(ctx, "No '--group' specified to export projects.")
	}

	mustutils.Must0(gitlab.Authenticate(ctx, &asciichgolangpublic.GitlabAuthenticationOptions{AccessToken: accessToken}))

	gitlab_fqdn := mustutils.Must(gitlab.GetFqdn())

	instanceMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_instance",
			Help: "Gitlab instance specific metrics",
		},
		[]string{"name", "value"},
	)
	error_in_last_loop := prometheus.NewGauge(prometheus.GaugeOpts{Name: "gitlab_collect_metrics_errors", Help: "Number of errors during last collection loop."})

	gitlab_project_n_scheduled_pipelines := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_project_n_scheduled_pipelines",
			Help: "Number of scheduled pipelines per gitlab project",
		},
		[]string{"project"},
	)

	scheduled_pipelines_status := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_project_scheduled_pipeline_status",
			Help: "Status of last run scheduled pipeline. 0=pending, 1=success, -1=error",
		},
		[]string{"project", "name"},
	)

	go func() {
		for {
			error_counter := 0

			version, revision, err := gitlab.GetVersionAnRevisionAsString(ctx)
			if err != nil {
				logging.LogGoError(err)
				error_counter += 1
			}

			instanceMetrics.Reset()
			instanceMetrics.WithLabelValues("instance", gitlab_fqdn).Set(1)
			instanceMetrics.WithLabelValues("version", version).Set(1)
			instanceMetrics.WithLabelValues("revision", revision).Set(1)

			nScheduledPipelines := map[string]int{}

			type ScheduledStatus struct {
				Url         string
				Name        string
				StatusValue int
			}

			scheduleStatuses := []ScheduledStatus{}

			if group != "" {
				g, err := gitlab.GetGroupByPath(ctx, group)
				if err != nil {
					logging.LogGoError(err)
					error_counter += 1
					g = nil
				}

				if g != nil {
					projects, err := g.ListProjects(contextutils.WithSilent(ctx), &asciichgolangpublic.GitlabListProjectsOptions{Recursive: true})
					if err != nil {
						logging.LogGoError(err)
						error_counter += 1
						projects = nil
					}

					logging.LogInfoByCtxf(ctx, "Collected '%d' projects in group '%s'", len(projects), group)

					for _, project := range projects {
						url, err := project.GetProjectUrl(ctx)
						if err != nil {
							logging.LogGoError(err)
							error_counter += 1
							continue
						}

						scheduledPipelines, err := project.ListScheduledPipelines(ctx)
						if err != nil {
							logging.LogGoError(err)
							error_counter += 1
							continue
						}

						nScheduledPipelines[url] = len(scheduledPipelines)

						for _, sp := range scheduledPipelines {
							status, err := sp.GetLastPipelineStatus(ctx)
							if err != nil {
								logging.LogGoError(err)
								error_counter += 1
								continue
							}

							statusValue := -1
							if status == "success" {
								statusValue = 1
							}

							if slices.Contains([]string{"pending", "running"}, status) {
								statusValue = 0
							}

							name, err := sp.GetCachedName()
							if err != nil {
								logging.LogGoError(err)
								error_counter += 1
								continue
							}

							scheduleStatuses = append(scheduleStatuses, ScheduledStatus{url, name, statusValue})
						}
					}
				}

				scheduled_pipelines_status.Reset()
				for _, s := range scheduleStatuses {
					scheduled_pipelines_status.WithLabelValues(s.Url, s.Name).Set(float64(s.StatusValue))
				}

				gitlab_project_n_scheduled_pipelines.Reset()
				for url, nScheduled := range nScheduledPipelines {
					gitlab_project_n_scheduled_pipelines.WithLabelValues(url).Set(float64(nScheduled))
				}
			}

			error_in_last_loop.Set(float64(error_counter))

			time.Sleep(10 * time.Second)
		}
	}()

	prometheus.MustRegister(
		instanceMetrics,
		error_in_last_loop,
		gitlab_project_n_scheduled_pipelines,
		scheduled_pipelines_status,
	)

	http.Handle("/metrics", promhttp.Handler())
	logging.LogInfoByCtxf(ctx, "Exposing metrics on http://localhost:%d/metrics", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
