package provider

import (
	"context"
	"time"
)

const (
	// Supported emissions types.
	AverageEmissionsType  = "average"
	MarginalEmissionsType = "marginal"

	// Supported metric types.
	AbsoluteMetricType = "absolute"
	RelativeMetricType = "relative"

	// Supported units.
	GramsCO2EPerkWh = "gCO2e per kWh"
	LbCO2EPerMWh    = "lbCO2e per MWh"
	Percent         = "percent"

	// Supported providers
	CarbonIntensityOrgUK = "CarbonIntensityOrgUK"
	ElectricityMap       = "ElectricityMap"
	Ember                = "Ember"
	WattTime             = "WattTime"
)

type CarbonIntensity struct {
	EmissionsType string    `json:"emissions_type"`
	MetricType    string    `json:"metric_type"`
	Provider      string    `json:"provider"`
	Region        string    `json:"region"`
	Units         string    `json:"units"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	Value         float64   `json:"value"`
}

type Details struct {
	Name string
	URL  string
}

type Interface interface {
	GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error)
}

func GetProviderDetails() []Details {
	return []Details{
		{
			Name: CarbonIntensityOrgUK,
			URL:  "carbonintensity.org.uk",
		},
		{
			Name: ElectricityMap,
			URL:  "electricitymap.org",
		},
		{
			Name: Ember,
			URL:  "ember-climate.org",
		},
		{
			Name: WattTime,
			URL:  "watttime.org",
		},
	}
}
