package main

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
	"github.com/thegreenwebfoundation/grid-intensity-go/carbonintensity"
)

const (
	labelProvider = "provider"
	labelRegion   = "region"
	namespace     = "grid_intensity"
)

var (
	actualDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "actual"),
		"Actual carbon intensity for this region.",
		[]string{
			labelProvider,
			labelRegion,
		},
		nil,
	)
)

type Exporter struct {
	apiClient gridintensity.Provider
	provider  string
	region    string
}

func NewExporter(provider, region string) *Exporter {
	apiClient, err := carbonintensity.New()
	if err != nil {
		log.Fatalln("could not make provider", err)
	}

	return &Exporter{
		apiClient: apiClient,
		provider:  provider,
		region:    region,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- actualDesc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	actualIntensity, err := e.apiClient.GetCarbonIntensity(ctx, e.region)
	if err != nil {
		log.Printf("failed to get carbon intensity %#v", err)
	}

	ch <- prometheus.MustNewConstMetric(
		actualDesc,
		prometheus.GaugeValue,
		actualIntensity,
		e.provider,
		e.region,
	)
}

func main() {
	exporter := NewExporter("carbonintensity.org.uk", "UK")
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
