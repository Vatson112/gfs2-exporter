package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"github.com/Vatson112/gfs2-exporter/gfs2"
)

type collector struct{}

var (
	scrapeDurationDesc = prometheus.NewDesc(
		"gfs2_scrape_duration_seconds",
		"gfs2_exporter: Duration of scraping gfs2.",
		nil,
		nil,
	)

	scrapeSuccessDesc = prometheus.NewDesc(
		"gfs2_scrape_success",
		"gfs2_exporter: Whether scraping gfs2 succeeded.",
		nil,
		nil,
	)

	glocksDesc = prometheus.NewDesc(
		"gfs2_glocks_total",
		"gfs2_glocks_total: Count of incore GFS2 glock data structures by state.",
		[]string{"cluster", "fs", "state"},
		nil,
	)

	// dcnt in glstats file
	dlmRequestsDesc = prometheus.NewDesc(
		"gfs2_dlm_requests",
		"gfs2_dlm_requests: Number of dlm requests made.",
		[]string{"cluster", "fs"},
		nil,
	)
	// qcnt in glstats file
	glocksRequestsDesc = prometheus.NewDesc(
		"gfs2_glocks_requests",
		"gfs2_glocks_requests: Number of glock requests queued.",
		[]string{"cluster", "fs"},
		nil,
	)
)

func (c *collector) Describe(descChan chan<- *prometheus.Desc) {
	descChan <- scrapeDurationDesc
	descChan <- scrapeSuccessDesc
	descChan <- glocksDesc
	descChan <- dlmRequestsDesc
	descChan <- glocksRequestsDesc
}

func (c *collector) Collect(metricChan chan<- prometheus.Metric) {
	start := time.Now()
	clusterMetric, err := gfs2.GetMetric()
	duration := time.Since(start)
	if err == nil && len(clusterMetric) == 0 {
		err = errors.New("Error on parsing or empty output.")
	}
	metricChan <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds())
	if err != nil {
		metricChan <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 0)
		log.Error(err)
		return
	}
	metricChan <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 1)
	for clusterName, cluster := range clusterMetric {
		for fsName, fs := range cluster {
			metricChan <- prometheus.MustNewConstMetric(glocksDesc, prometheus.CounterValue, float64(fs.Value), clusterName, fsName, fs.State)
			metricChan <- prometheus.MustNewConstMetric(dlmRequestsDesc, prometheus.CounterValue, float64(fs.Value), clusterName, fsName)
			metricChan <- prometheus.MustNewConstMetric(glocksRequestsDesc, prometheus.CounterValue, float64(fs.Value), clusterName, fsName)
		}

	}
}

func main() {
	var (
		listenAddress string = ":9457"
		metricsPath   string = "/metrics"
	)

	// listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9457").String()
	flag.StringVar(&listenAddress, "-listen_address", listenAddress, "Address on which to expose metrics and web interface. Default: :9457 ")
	// metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	flag.StringVar(&metricsPath, "-metrics_path", metricsPath, "Address on which to expose metrics and web interface. Default: /metrics ")
	flag.Parse()

	log.Infoln("Starting gfs2_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>gfs2 exporter</title></head>
			<body>
			<h1>gfs2 exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Infoln("Listening on", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)

	if err != nil {
		log.Fatal(err)
	}
}
