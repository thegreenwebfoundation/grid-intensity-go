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

	// Supported providers
	CarbonIntensityOrgUK = "carbonintensity.org.uk"
	Ember                = "ember-climate.org"
	ElectricityMap       = "electricitymap.org"
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

type Interface interface {
	GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error)
}
