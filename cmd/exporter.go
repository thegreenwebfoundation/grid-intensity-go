package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

const (
	labelLocation = "location"
	labelNode     = "node"
	labelProvider = "provider"
	labelRegion   = "region"
	labelUnits    = "units"
	namespace     = "grid_intensity"
	nodeKey       = "node"
	regionKey     = "region"
)

func init() {
	exporterCmd.Flags().StringP(locationKey, "l", "", "Location code for provider")
	exporterCmd.Flags().StringP(nodeKey, "n", "", "Node where the exporter is running")
	exporterCmd.Flags().StringP(providerKey, "p", provider.Ember, "Provider of carbon intensity data")
	exporterCmd.Flags().StringP(regionKey, "r", "", "Region where the exporter is running")

	// Also support environment variables.
	viper.SetEnvPrefix("grid_intensity")
	viper.BindEnv(locationKey)
	viper.BindEnv(providerKey)
	viper.BindEnv(regionKey)
	viper.BindEnv(nodeKey)

	rootCmd.AddCommand(exporterCmd)
}

var (
	averageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "average"),
		"Average carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelNode,
			labelProvider,
			labelRegion,
			labelUnits,
		},
		nil,
	)

	marginalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "marginal"),
		"Marginal carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelNode,
			labelProvider,
			labelRegion,
			labelUnits,
		},
		nil,
	)

	relativeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "relative"),
		"Relative carbon intensity for the electricity grid in this location.",
		[]string{
			labelLocation,
			labelNode,
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

	grid-intensity exporter --provider Ember --location IE --region eu-west-1 --node worker-1
	grid-intensity exporter -p Ember -l BOL`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(locationKey, cmd.Flags().Lookup(locationKey))
			viper.BindPFlag(nodeKey, cmd.Flags().Lookup(nodeKey))
			viper.BindPFlag(providerKey, cmd.Flags().Lookup(providerKey))
			viper.BindPFlag(regionKey, cmd.Flags().Lookup(regionKey))
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
	node     string
	provider string
	region   string
}

type ExporterConfig struct {
	Location string
	Node     string
	Provider string
	Region   string
}

func NewExporter(config ExporterConfig) (*Exporter, error) {
	var client provider.Interface
	var err error

	if config.Location == "" {
		return nil, fmt.Errorf("location must be set")
	}

	// Cache filename is empty so we use in-memory cache.
	client, err = getClient(config.Provider, "")
	if err != nil {
		return nil, err
	}

	e := &Exporter{
		client:   client,
		location: config.Location,
		node:     config.Node,
		provider: config.Provider,
		region:   config.Region,
	}

	return e, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	var result []provider.CarbonIntensity
	locationCodes := strings.Split(e.location, ",")

	for _, locationCode := range locationCodes {
		res, err := e.client.GetCarbonIntensity(ctx, locationCode)
		if err != nil {
			log.Printf("could not get carbon intensity for location %s, %#v", locationCode, err)
		}
		result = append(result, res...)
	}

	for _, data := range result {
		desc, err := getMetricDesc(data)
		if err != nil {
			log.Printf("failed to get metric description %#v", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			data.Value,
			data.Location,
			e.node,
			data.Provider,
			e.region,
			data.Units,
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

func getMetricDesc(data provider.CarbonIntensity) (*prometheus.Desc, error) {
	switch data.MetricType {
	case provider.AbsoluteMetricType:
		switch data.EmissionsType {
		case provider.AverageEmissionsType:
			return averageDesc, nil
		case provider.MarginalEmissionsType:
			return marginalDesc, nil
		default:
			return nil, fmt.Errorf("unknown emissions type %s", data.EmissionsType)
		}
	case provider.RelativeMetricType:
		return relativeDesc, nil
	}

	return nil, fmt.Errorf("unknown metric type %s", data.MetricType)
}

func runExporter() error {
	providerName, err := readConfig(providerKey)
	if err != nil {
		return err
	}
	locationCode, err := readConfig(locationKey)
	if err != nil {
		return err
	}
	node, err := readConfig(nodeKey)
	if err != nil {
		return err
	}
	region, err := readConfig(regionKey)
	if err != nil {
		return err
	}

	c := ExporterConfig{
		Location: locationCode,
		Node:     node,
		Provider: providerName,
		Region:   region,
	}
	exporter, err := NewExporter(c)
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
