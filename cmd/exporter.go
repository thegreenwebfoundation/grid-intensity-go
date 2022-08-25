package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

const (
	labelLocation = "location"
	labelProvider = "provider"
	labelUnits    = "units"
	namespace     = "grid_intensity"
)

func init() {
	exporterCmd.Flags().StringP(locationKey, "l", "", "Location code for provider")
	exporterCmd.Flags().StringP(providerKey, "p", provider.Ember, "Provider of carbon intensity data")

	rootCmd.AddCommand(exporterCmd)
}

var (
	averageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "average"),
		"Average carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelProvider,
			labelUnits,
		},
		nil,
	)

	marginalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "marginal"),
		"Marginal carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelProvider,
			labelUnits,
		},
		nil,
	)

	relativeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "relative"),
		"Relative carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelProvider,
			labelUnits,
		},
		nil,
	)

	exporterCmd = &cobra.Command{
		Use:   "exporter",
		Short: "Metrics for carbon intensity data for electricity grids",
		Long: `A prometheus exporter for getting the carbon intensity data for
electricity grids.

This can be used to make your software carbon aware so it runs at times when
the grid is greener or at locations where carbon intensity is lower.

	grid-intensity exporter --provider Ember --location ARG
	grid-intensity exporter -p Ember -l BOL`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(providerKey, cmd.Flags().Lookup(providerKey))
			viper.BindPFlag(locationKey, cmd.Flags().Lookup(locationKey))
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := runExporter()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

type Exporter struct {
	client   provider.Interface
	location string
	provider string
}

func NewExporter(providerName, locationName string) (*Exporter, error) {
	var client provider.Interface
	var err error

	if locationName == "" {
		return nil, fmt.Errorf("location must be set")
	}

	// Cache filename is empty so we use in-memory cache.
	client, err = getClient(providerName, "")
	if err != nil {
		return nil, err
	}

	e := &Exporter{
		client:   client,
		location: locationName,
		provider: providerName,
	}

	return e, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	if e.provider == provider.WattTime {
		result, err := e.client.GetCarbonIntensity(ctx, e.location)
		if err != nil {
			log.Printf("failed to get carbon intensity %#v", err)
		}

		for _, data := range result {
			if data.MetricType == provider.AbsoluteMetricType {
				ch <- prometheus.MustNewConstMetric(
					marginalDesc,
					prometheus.GaugeValue,
					data.Value,
					data.Location,
					data.Provider,
					data.Units,
				)
			}
			if data.MetricType == provider.RelativeMetricType {
				ch <- prometheus.MustNewConstMetric(
					relativeDesc,
					prometheus.GaugeValue,
					data.Value,
					data.Location,
					data.Provider,
					data.Units,
				)
			}
		}
	} else {
		result, err := e.client.GetCarbonIntensity(ctx, e.location)
		if err != nil {
			log.Printf("failed to get carbon intensity %#v", err)
		}

		averageIntensity := result[0]
		ch <- prometheus.MustNewConstMetric(
			averageDesc,
			prometheus.GaugeValue,
			averageIntensity.Value,
			averageIntensity.Location,
			averageIntensity.Provider,
			averageIntensity.Units,
		)
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.provider == provider.WattTime {
		ch <- marginalDesc
		ch <- relativeDesc
	} else {
		ch <- averageDesc
	}
}

func runExporter() error {
	providerName, locationCode, err := readConfig()
	if err != nil {
		return err
	}

	exporter, err := NewExporter(providerName, locationCode)
	if err != nil {
		return err
	}

	err = writeConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Using provider %q with location %q\n", providerName, locationCode)
	fmt.Println("Metrics available at :8000/metrics")

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8000", nil))

	return nil
}
