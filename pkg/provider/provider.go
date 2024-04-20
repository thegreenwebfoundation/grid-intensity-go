package provider

import (
	"context"
	"net/url"
	"path"
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
	ElectricityMaps      = "ElectricityMaps"
	Ember                = "Ember"
	WattTime             = "WattTime"
)

type CarbonIntensity struct {
	EmissionsType string    `json:"emissions_type"`
	MetricType    string    `json:"metric_type"`
	Provider      string    `json:"provider"`
	Location      string    `json:"location"`
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
	GetCarbonIntensity(ctx context.Context, location string) ([]CarbonIntensity, error)
}

func GetProviderDetails() []Details {
	return []Details{
		{
			Name: CarbonIntensityOrgUK,
			URL:  "carbonintensity.org.uk",
		},
		{
			Name: ElectricityMaps,
			URL:  "electricitymaps.com",
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

func buildURL(apiURL, relativePath string) (string, error) {
	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return "", err
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return "", err
	}

	// Safely add relative path.
	baseURL.Path = path.Join(baseURL.Path, relativeURL.Path)

	// Safely merge query strings.
	baseQuery := baseURL.Query()

	for param, values := range relativeURL.Query() {
		for _, value := range values {
			baseQuery.Add(param, value)
		}
	}

	baseURL.RawQuery = baseQuery.Encode()
	return baseURL.String(), nil
}
