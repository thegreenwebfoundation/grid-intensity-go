package carbonintensity

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

type ApiOption func(*ApiClient) error

func New(opts ...ApiOption) (gridintensity.Provider, error) {
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

func (a *ApiClient) GetCarbonIndex(ctx context.Context, region string) (gridintensity.CarbonIndex, error) {
	if region != "UK" {
		return gridintensity.UNKNOWN, ErrOnlyUK
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.apiURL, nil)
	if err != nil {
		return gridintensity.UNKNOWN, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return gridintensity.UNKNOWN, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gridintensity.UNKNOWN, errBadStatus(resp)
	}

	respObj := &CarbonIntensityResponse{}

	err = json.NewDecoder(resp.Body).Decode(respObj)
	if err != nil {
		return gridintensity.UNKNOWN, err
	}

	if len(respObj.Data) == 0 || respObj.Data[0].Intensity == nil {
		return gridintensity.UNKNOWN, ErrNoResponse
	}

	latestData := respObj.Data[0]
	switch latestData.Intensity.Index {
	case "very low", "low":
		return gridintensity.LOW, nil
	case "moderate":
		return gridintensity.MODERATE, nil
	case "high", "very high":
		return gridintensity.HIGH, nil
	}
	return gridintensity.UNKNOWN, errUnknownIndex(latestData.Intensity.Index)
}

func (a *ApiClient) MinIntensity(ctx context.Context, regions ...string) (string, error) {
	return "", ErrOnlyUK
}
