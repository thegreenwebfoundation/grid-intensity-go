package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go/api"
	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

const (
	labelProvider = "provider"
	labelRegion   = "region"
	labelUnits    = "units"
	namespace     = "grid_intensity"
)

func init() {
	rootCmd.AddCommand(exporterCmd)
}

var (
	averageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "average"),
		"Average carbon intensity for the electricity grid in this region.",
		[]string{
			labelProvider,
			labelRegion,
			labelUnits,
		},
		nil,
	)

	marginalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "marginal"),
		"Marginal carbon intensity for the electricity grid in this region.",
		[]string{
			labelProvider,
			labelRegion,
			labelUnits,
		},
		nil,
	)

	relativeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "relative"),
		"Relative carbon intensity for the electricity grid in this region.",
		[]string{
			labelProvider,
			labelRegion,
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

	grid-intensity exporter --provider PROVIDER --region ARG
	grid-intensity exporter -p ember-climate.org -r BOL`,
		Run: func(cmd *cobra.Command, args []string) {
			err := runExporter()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

type Exporter struct {
	apiClient gridintensity.Provider
	client    provider.Interface
	provider  string
	region    string
	units     string
}

func NewExporter(providerName, regionName string) (*Exporter, error) {
	var client provider.Interface

	var units string
	var err error

	if regionName == "" {
		return nil, fmt.Errorf("region must be set")
	}

	switch providerName {
	case provider.CarbonIntensityOrgUK:
		c := provider.CarbonIntensityUKConfig{}
		client, err = provider.NewCarbonIntensityUK(c)
		if err != nil {
			return nil, err
		}
	case provider.ElectricityMap:
		token := os.Getenv(electricityMapAPITokenEnvVar)
		if token == "" {
			return nil, fmt.Errorf("%q env var must be set", electricityMapAPITokenEnvVar)
		}

		c := provider.ElectricityMapConfig{
			Token: token,
		}
		client, err = provider.NewElectricityMap(c)
		if err != nil {
			return nil, err
		}
	case provider.Ember:
		client, err = provider.NewEmber()
		if err != nil {
			return nil, err
		}
	case provider.WattTime:
		user := os.Getenv(wattTimeUserEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimeUserEnvVar)
		}

		password := os.Getenv(wattTimePasswordEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimePasswordEnvVar)
		}

		c := provider.WattTimeConfig{
			APIUser:     user,
			APIPassword: password,
		}
		client, err = provider.NewWattTime(c)
		if err != nil {
			return nil, fmt.Errorf("could not make provider %v", err)
		}
	default:
		return nil, fmt.Errorf("provider %q not supported", providerName)
	}

	e := &Exporter{
		client:   client,
		provider: providerName,
		region:   regionName,
		units:    units,
	}

	return e, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	if e.provider == provider.WattTime {
		result, err := e.client.GetCarbonIntensity(ctx, e.region)
		if err != nil {
			log.Printf("failed to get carbon intensity %#v", err)
		}

		for _, data := range result {
			if data.MetricType == provider.AbsoluteMetricType {
				ch <- prometheus.MustNewConstMetric(
					marginalDesc,
					prometheus.GaugeValue,
					data.Value,
					data.Provider,
					data.Region,
					data.Units,
				)
			}
			if data.MetricType == provider.RelativeMetricType {
				ch <- prometheus.MustNewConstMetric(
					relativeDesc,
					prometheus.GaugeValue,
					data.Value,
					data.Provider,
					data.Region,
					data.Units,
				)
			}
		}
	} else {
		result, err := e.client.GetCarbonIntensity(ctx, e.region)
		if err != nil {
			log.Printf("failed to get carbon intensity %#v", err)
		}

		averageIntensity := result[0]
		ch <- prometheus.MustNewConstMetric(
			averageDesc,
			prometheus.GaugeValue,
			averageIntensity.Value,
			averageIntensity.Provider,
			averageIntensity.Region,
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
	providerName, regionCode, err := readConfig()
	if err != nil {
		return err
	}

	exporter, err := NewExporter(providerName, regionCode)
	if err != nil {
		return err
	}

	err = writeConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Using provider %q with region %q\n", providerName, regionCode)
	fmt.Println("Metrics available at :8000/metrics")

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8000", nil))

	return nil
}
