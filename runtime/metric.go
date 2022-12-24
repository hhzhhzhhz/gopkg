package runtime

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

const Metric = ":7424"

func StartMetric(addr string) error {
	log.Logger().Info(fmt.Sprintf("metrics is listening and serving on %s", addr))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Errorf("%s", http.ListenAndServe(addr, nil))
	}()
	return nil
}

type Example struct {
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	opt             *OptionMetric
}

type OptionMetric struct {
	Addr      string
	Namespace string
}

func NewExample(opt *OptionMetric) *Example {
	return &Example{opt: opt}
}

func (m *Example) defineMetrics() {
	hostname, _ := os.Hostname()
	m.requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: m.opt.Namespace,
		Subsystem: hostname,
		Name:      "request_count_total",
		Help:      "Counter of DNS requests made.",
	}, []string{"system"})

	m.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: m.opt.Namespace,
		Subsystem: hostname,
		Name:      "request_duration_seconds",
		Help:      "Histogram of the time (in seconds) each request took to resolve.",
		Buckets:   append([]float64{0.001, 0.003}, prometheus.DefBuckets...),
	}, []string{"system"})

	m.responseSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: m.opt.Namespace,
		Subsystem: hostname,
		Name:      "response_size_bytes",
		Help:      "Size of the returns response in bytes.",
		Buckets: []float64{0, 512, 1024, 1500, 2048, 4096,
			8192, 12288, 16384, 20480, 24576, 28672, 32768, 36864,
			40960, 45056, 49152, 53248, 57344, 61440, 65536,
		},
	}, []string{"system"})
}

func (m *Example) ReportRequestCount(req string, sys string) {
	if m.requestCount == nil {
		return
	}
	m.requestCount.WithLabelValues(string(sys)).Inc()
}

func (m *Example) Run() error {
	m.defineMetrics()
	prometheus.MustRegister(m.requestCount)
	prometheus.MustRegister(m.requestDuration)
	prometheus.MustRegister(m.responseSize)
	fmt.Println(fmt.Sprintf("metrics is listening and serving on :%s", m.opt.Addr))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Errorf("%s", http.ListenAndServe(m.opt.Addr, nil))
	}()
	return nil
}
