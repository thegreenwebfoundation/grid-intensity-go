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
	CarbonIntensityOrgUK = "CarbonIntensityOrgUK"
	ElectricityMap       = "ElectricityMap"
)

type CarbonIntensity struct {
	DataProvider  string    `json:"data_provider"`
	EmissionsType string    `json:"emissions_type"`
	MetricType    string    `json:"metric_type"`
	Region        string    `json:"region"`
	Units         string    `json:"units"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	Value         float64   `json:"value"`
}

type Interface interface {
	GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error)
}
