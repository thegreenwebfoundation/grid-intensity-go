package electricitymap

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

	return a, nil
}

type ApiClient struct {
	client *http.Client
	apiURL string
	token  string
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
