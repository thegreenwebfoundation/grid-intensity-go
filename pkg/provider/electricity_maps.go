package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ElectricityMapsClient struct {
	client *http.Client
	apiURL string
	token  string
}

type ElectricityMapsConfig struct {
	Client *http.Client
	APIURL string
	Token  string
}

func NewElectricityMaps(config ElectricityMapsConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.electricitymap.org/v3"
	}

	c := &ElectricityMapsClient{
		apiURL: config.APIURL,
		client: config.Client,
		token:  config.Token,
	}

	return c, nil
}

func (e *ElectricityMapsClient) GetCarbonIntensity(ctx context.Context, location string) ([]CarbonIntensity, error) {
	intensityURL, err := e.intensityURLWithZone(location)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, intensityURL, nil)
	req.Header.Add("auth-token", e.token)
	if err != nil {
		return nil, err
	}

	log.Printf("calling %s", req.URL)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	data := &electricityMapsData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	validFrom, err := time.Parse(time.RFC3339Nano, data.DateTime)
	if err != nil {
		return nil, err
	}
	validTo := validFrom.Add(60 * time.Minute)

	return []CarbonIntensity{
		{
			EmissionsType: AverageEmissionsType,
			MetricType:    AbsoluteMetricType,
			Provider:      ElectricityMaps,
			Location:      location,
			Units:         GramsCO2EPerkWh,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         data.CarbonIntensity,
		},
	}, nil
}

func (e *ElectricityMapsClient) intensityURLWithZone(zone string) (string, error) {
	zonePath := fmt.Sprintf("/carbon-intensity/latest?zone=%s", zone)
	return buildURL(e.apiURL, zonePath)
}

type electricityMapsData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
	DateTime        string  `json:"datetime"`
	UpdatedAt       string  `json:"updatedAt"`
}
