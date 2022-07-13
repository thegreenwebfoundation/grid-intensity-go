package carbonintensity

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const ProviderName = "carbonintensity.org.uk"

type ApiOption func(*ApiClient) error

func New(opts ...ApiOption) (*ApiClient, error) {
	a := &ApiClient{}
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}

	if a.client == nil {
		a.client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	if a.apiURL == "" {
		a.apiURL = "https://api.carbonintensity.org.uk/intensity/"
	}

	return a, nil
}

type ApiClient struct {
	client *http.Client
	apiURL string
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	data, err := a.getCarbonIntensityData(ctx, region)
	if err != nil {
		return 0, err
	}
	if data.Intensity == nil {
		return 0, ErrNoResponse
	}
	return data.Intensity.Actual, nil
}

func (a *ApiClient) GetCarbonIntensityData(ctx context.Context, region string) (*IntensityData, error) {
	data, err := a.getCarbonIntensityData(ctx, region)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *ApiClient) getCarbonIntensityData(ctx context.Context, region string) (*IntensityData, error) {
	if region != "UK" {
		return nil, ErrOnlyUK
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.apiURL, nil)
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

	respObj := &CarbonIntensityResponse{}

	err = json.NewDecoder(resp.Body).Decode(respObj)
	if err != nil {
		return nil, err
	}

	if len(respObj.Data) == 0 {
		return nil, ErrNoResponse
	}
	return &respObj.Data[0], nil
}
