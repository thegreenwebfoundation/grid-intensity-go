package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go/api"
	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
	"github.com/thegreenwebfoundation/grid-intensity-go/watttime"
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
	apiClient      gridintensity.Provider
	client         provider.Interface
	provider       string
	region         string
	units          string
	wattTimeClient watttime.Provider
}

func NewExporter(providerName, regionName string) (*Exporter, error) {
	var client provider.Interface
	var wattTimeClient watttime.Provider

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
	case watttime.ProviderName:
		user := os.Getenv(wattTimeUserEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimeUserEnvVar)
		}

		password := os.Getenv(wattTimePasswordEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimePasswordEnvVar)
		}

		wattTimeClient, err = watttime.New(user, password)
		if err != nil {
			return nil, fmt.Errorf("could not make provider %v", err)
		}
		units = "lb CO2 per MWh"
	default:
		return nil, fmt.Errorf("provider %q not supported", providerName)
	}

	e := &Exporter{
		provider: providerName,
		region:   regionName,
		units:    units,
	}
	if providerName == watttime.ProviderName {
		e.wattTimeClient = wattTimeClient
	} else {
		e.client = client
	}

	return e, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	if e.provider == watttime.ProviderName {
		data, err := e.wattTimeClient.GetCarbonIntensityData(ctx, e.region)
		if err != nil {
			log.Printf("failed to get carbon intensity %#v", err)
		}

		if data.MOER != "" {
			marginalIntensity, err := strconv.ParseFloat(data.MOER, 64)
			if err != nil {
				log.Printf("failed to parse marginal intensity %#v", err)
			}

			ch <- prometheus.MustNewConstMetric(
				marginalDesc,
				prometheus.GaugeValue,
				marginalIntensity,
				e.provider,
				e.region,
				e.units,
			)
		}

		if data.Percent != "" {
			relativeIntensity, err := strconv.ParseFloat(data.Percent, 64)
			if err != nil {
				log.Printf("failed to parse relative intensity %#v", err)
			}

			ch <- prometheus.MustNewConstMetric(
				relativeDesc,
				prometheus.GaugeValue,
				relativeIntensity,
				e.provider,
				e.region,
				"percent",
			)
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
	if e.provider == watttime.ProviderName {
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
