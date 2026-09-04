package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	gridCarbonDefaultAPIURL = "https://api.gridcarbon.dev"
	gridCarbonUserAgent     = "grid-intensity-go"
)

// GridCarbonClient reads hourly average carbon intensity from gridcarbon
// (https://gridcarbon.dev): a keyless public API covering 33 European bidding
// zones (ENTSO-E), 11 US balancing authorities (EIA-930) and Great Britain (NESO).
type GridCarbonClient struct {
	client *http.Client
	apiURL string
}

type GridCarbonConfig struct {
	Client *http.Client
	// APIURL is the API origin without a trailing path. Defaults to https://api.gridcarbon.dev.
	APIURL string
}

func NewGridCarbon(config GridCarbonConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = gridCarbonDefaultAPIURL
	}

	c := &GridCarbonClient{
		client: config.Client,
		apiURL: strings.TrimRight(config.APIURL, "/"),
	}

	return c, nil
}

// GetCarbonIntensity returns the newest published hourly value for one zone.
// location is a gridcarbon zone code such as FR, DE-LU, US-NYISO or GB; the full
// list is served at /v1/zones. Unknown zones return ErrInvalidLocation.
//
// Values are average, production-based and lifecycle (IPCC AR5) for every zone
// except GB, which republishes NESO's operational (combustion-only) figure and is
// therefore not comparable in level with the other zones. GB intervals are 30
// minutes and the newest one is NESO's forecast, reported here as estimated.
func (a *GridCarbonClient) GetCarbonIntensity(ctx context.Context, location string) ([]CarbonIntensity, error) {
	location = strings.ToUpper(strings.TrimSpace(location))
	if location == "" {
		return nil, ErrInvalidLocation
	}

	reqURL := fmt.Sprintf("%s/v1/intensity/latest?zone=%s", a.apiURL, url.QueryEscape(location))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", gridCarbonUserAgent)

	log.Printf("calling %s", req.URL)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// {"error":"unknown or empty zone: XX"}
		return nil, ErrInvalidLocation
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	respObj := &gridCarbonResponse{}
	if err := json.NewDecoder(resp.Body).Decode(respObj); err != nil {
		return nil, err
	}
	if len(respObj.Data) == 0 {
		return nil, ErrNoResponse
	}

	result := make([]CarbonIntensity, 0, len(respObj.Data))
	for _, d := range respObj.Data {
		validFrom, err := time.Parse(time.RFC3339, d.TS)
		if err != nil {
			return nil, err
		}
		// Every zone is hourly except GB, which is settled in half-hours.
		step := time.Hour
		if d.Zone == "GB" {
			step = 30 * time.Minute
		}

		result = append(result, CarbonIntensity{
			EmissionsType: AverageEmissionsType,
			MetricType:    AbsoluteMetricType,
			Provider:      GridCarbon,
			Location:      d.Zone,
			Units:         GramsCO2EPerkWh,
			ValidFrom:     validFrom,
			ValidTo:       validFrom.Add(step),
			Value:         d.GCO2eqPerKWh,
			// method is "computed:v1" for values derived from published generation and
			// "upstream:uk-neso:forecast" when the newest GB half-hour is still a forecast.
			IsEstimated: strings.HasSuffix(d.Method, ":forecast"),
		})
	}

	return result, nil
}

type gridCarbonResponse struct {
	Unit string           `json:"unit"`
	Data []gridCarbonData `json:"data"`
}

type gridCarbonData struct {
	Zone         string  `json:"zone"`
	TS           string  `json:"ts"`
	GCO2eqPerKWh float64 `json:"gco2eq_kwh"`
	Method       string  `json:"method"`
}
