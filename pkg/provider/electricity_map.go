package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ElectricityMapClient struct {
	client *http.Client
	apiURL string
	token  string
}

type ElectricityMapConfig struct {
	Client *http.Client
	APIURL string
	Token  string
}

func NewElectricityMap(config ElectricityMapConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.electricitymap.org/v3"
	}

	c := &ElectricityMapClient{
		apiURL: config.APIURL,
		client: config.Client,
		token:  config.Token,
	}

	return c, nil
}

func (e *ElectricityMapClient) GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.intensityURLWithZone(region), nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	data := &electricityMapData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	validFrom, err := time.Parse(time.RFC3339Nano, data.UpdatedAt)
	if err != nil {
		return nil, err
	}
	validTo, err := time.Parse(time.RFC3339Nano, data.DateTime)
	if err != nil {
		return nil, err
	}

	return []CarbonIntensity{
		{
			EmissionsType: AverageEmissionsType,
			MetricType:    AbsoluteMetricType,
			Provider:      ElectricityMap,
			Region:        region,
			Units:         GramsCO2EPerkWh,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         data.CarbonIntensity,
		},
	}, nil
}

func (e *ElectricityMapClient) intensityURLWithZone(zone string) string {
	return fmt.Sprintf("%s/carbon-intensity/latest?zone=%s", e.apiURL, zone)
}

type electricityMapData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
	DateTime        string  `json:"datetime"`
	UpdatedAt       string  `json:"updatedAt"`
}
