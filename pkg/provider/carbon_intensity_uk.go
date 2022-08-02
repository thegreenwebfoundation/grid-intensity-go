package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CarbonIntensityUKClient struct {
	client *http.Client
	apiURL string
}

type CarbonIntensityUKConfig struct {
	Client *http.Client
	APIURL string
}

func NewCarbonIntensityUK(config CarbonIntensityUKConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.carbonintensity.org.uk/intensity/"
	}

	c := &CarbonIntensityUKClient{
		client: config.Client,
		apiURL: config.APIURL,
	}

	return c, nil
}

func (a *CarbonIntensityUKClient) GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error) {
	if region != "UK" {
		return nil, ErrInvalidRegion
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	respObj := &carbonIntensityUKResponse{}

	err = json.NewDecoder(resp.Body).Decode(respObj)
	if err != nil {
		return nil, err
	}

	if len(respObj.Data) == 0 {
		return nil, ErrNoResponse
	}

	data := &respObj.Data[0]
	if data.Intensity == nil {
		return nil, ErrNoResponse
	}

	layout := "2006-01-02T15:04Z"
	validFrom, err := time.Parse(layout, data.From)
	if err != nil {
		return nil, err
	}
	validTo, err := time.Parse(layout, data.To)
	if err != nil {
		return nil, err
	}

	return []CarbonIntensity{
		{
			DataProvider:  CarbonIntensityOrgUK,
			EmissionsType: AverageEmissionsType,
			MetricType:    AbsoluteMetricType,
			Region:        region,
			Units:         GramsCO2EPerkWh,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         data.Intensity.Actual,
		},
	}, nil
}

type carbonIntensityUKResponse struct {
	Data []carbonIntensityUKData `json:"data"`
}

type carbonIntensityUKData struct {
	From      string                      `json:"from"`
	To        string                      `json:"to"`
	Intensity *carbonIntensityUKIntensity `json:"intensity"`
}

type carbonIntensityUKIntensity struct {
	Forecast float64 `json:"forecast"`
	Actual   float64 `json:"actual"`
	Index    string  `json:"index"`
}
