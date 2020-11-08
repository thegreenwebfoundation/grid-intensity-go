package energymap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

type ApiOption func(*ApiClient) error

func New(token string, opts ...ApiOption) (gridintensity.Provider, error) {
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
		a.apiURL = "https://api.electricitymap.org/v3"
	}
	a.token = token

	if a.moderate_threshold == 0.0 {
		a.moderate_threshold = 150.0
	}
	if a.high_threshold == 0.0 {
		a.high_threshold = 300.0
	}

	return a, nil
}

type ApiClient struct {
	client             *http.Client
	apiURL             string
	token              string
	moderate_threshold float64
	high_threshold     float64
}

func (a *ApiClient) GetCarbonIndex(ctx context.Context, region string) (gridintensity.CarbonIndex, error) {
	intensity, err := a.GetCarbonIntensity(ctx, region)
	if err != nil {
		return gridintensity.UNKNOWN, err
	}

	if intensity >= a.high_threshold {
		return gridintensity.HIGH, nil
	}

	if intensity >= a.moderate_threshold {
		return gridintensity.MODERATE, nil
	}

	return gridintensity.LOW, nil
}

func (a *ApiClient) MinIntensity(ctx context.Context, regions ...string) (string, error) {

	if len(regions) == 0 {
		return "", ErrNoRegionProvided
	}

	requestCounter := len(regions)

	intensityMap := &IntensityMap{
		m: make(map[string]float64, requestCounter),
	}
	errChan := make(chan error, requestCounter)

	for _, region := range regions {
		go func(r string) {
			intensity, err := a.GetCarbonIntensity(ctx, r)
			errChan <- err
			if err != nil {
				return
			}
			intensityMap.Set(r, intensity)
		}(region)
	}

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return "", err
			}
			requestCounter--
			if requestCounter == 0 {
				r, err := intensityMap.Min()
				if err != nil {
					return "", err
				}
				return r, nil
			}
		case <-ctx.Done():
			return "", ErrTimeout
		}
	}
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.intensityURLWithZone(region), nil)
	if err != nil {
		return 0, err
	}
	resp, err := a.do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errBadStatus(resp)
	}

	respObj := &CarbonIntensityResp{}
	err = json.NewDecoder(resp.Body).Decode(respObj)
	if err != nil {
		return 0, err
	}
	return respObj.CarbonIntensity, nil
}

func (a *ApiClient) intensityURLWithZone(zone string) string {
	return fmt.Sprintf("%s/carbon-intensity/latest?zone=%s", a.apiURL, zone)
}

func (a *ApiClient) do(req *http.Request) (*http.Response, error) {
	// Add auth headers
	req.Header.Set("auth-token", a.token)
	return a.client.Do(req)
}
