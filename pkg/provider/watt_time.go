package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WattTimeClient struct {
	cache       *cacheStore
	client      *http.Client
	apiURL      string
	apiUser     string
	apiPassword string
	token       string
}

type WattTimeConfig struct {
	Client      *http.Client
	APIURL      string
	APIUser     string
	APIPassword string
	CacheFile   string
}

func NewWattTime(config WattTimeConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = "https://api2.watttime.org/v2"
	}

	c := cacheConfig{
		CacheFile: config.CacheFile,
	}
	cache, err := NewCacheStore(c)
	if err != nil {
		return nil, fmt.Errorf("could not make cache %v", err)
	}

	w := &WattTimeClient{
		cache:       cache,
		client:      config.Client,
		apiURL:      config.APIURL,
		apiUser:     config.APIUser,
		apiPassword: config.APIPassword,
	}

	return w, nil
}

func (w *WattTimeClient) GetCarbonIntensity(ctx context.Context, location string) ([]CarbonIntensity, error) {
	result, err := w.fetchCarbonIntensityData(ctx, location)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w *WattTimeClient) fetchCarbonIntensityData(ctx context.Context, location string) ([]CarbonIntensity, error) {
	result, err := w.cache.getCacheData(ctx, location)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	if w.token == "" {
		token, err := w.getAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		w.token = token
	}

	indexData, err := w.getCarbonIntensityData(ctx, location)
	if errors.Is(err, ErrReceived403Forbidden) {
		token, err := w.getAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		w.token = token

		indexData, err = w.getCarbonIntensityData(ctx, location)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	result, ttl, err := parseCarbonIntensityData(ctx, location, indexData)
	if err != nil {
		return nil, err
	}

	err = w.cache.setCacheData(ctx, location, result, ttl)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w *WattTimeClient) getAccessToken(ctx context.Context) (string, error) {
	loginURL, err := w.loginURL()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loginURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(w.apiUser, w.apiPassword)

	log.Printf("calling %s", req.URL)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errBadStatus(resp)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	loginResp := wattTimeLoginResp{}
	err = json.Unmarshal(bytes, &loginResp)
	if err != nil {
		return "", nil
	}

	return loginResp.Token, nil
}

func (w *WattTimeClient) getCarbonIntensityData(ctx context.Context, location string) (*wattTimeIndexData, error) {
	indexURL, err := w.indexURL(location)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, indexURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.token))

	log.Printf("calling %s", req.URL)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrReceived403Forbidden
	} else if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	indexData := wattTimeIndexData{}
	err = json.Unmarshal(bytes, &indexData)
	if err != nil {
		return nil, err
	}

	return &indexData, nil
}

func (w *WattTimeClient) indexURL(location string) (string, error) {
	indexPath := fmt.Sprintf("/index?ba=%s", location)
	return buildURL(w.apiURL, indexPath)
}

func (w *WattTimeClient) loginURL() (string, error) {
	return buildURL(w.apiURL, "/login")
}

func parseCarbonIntensityData(ctx context.Context, location string, indexData *wattTimeIndexData) ([]CarbonIntensity, time.Time, error) {
	freq, err := strconv.ParseInt(indexData.Freq, 0, 64)
	if err != nil {
		return nil, time.Time{}, err
	}

	validFrom := indexData.PointTime
	validTo := validFrom.Add(time.Duration(freq) * time.Second)

	ttl := validTo
	if ttl.Before(time.Now()) {
		// The TTL calculated from the point time is in the past. So reset the
		// TTL using the current time plus the frequency provided by the API.
		// UTC is used to match the WattTime API.
		ttl = time.Now().UTC().Add(time.Duration(freq) * time.Second)
	}

	result := []CarbonIntensity{}

	if indexData.Percent != "" {
		percent, err := strconv.ParseFloat(indexData.Percent, 64)
		if err != nil {
			return nil, time.Time{}, err
		}
		relative := CarbonIntensity{
			EmissionsType: MarginalEmissionsType,
			MetricType:    RelativeMetricType,
			Provider:      WattTime,
			Location:      location,
			Units:         Percent,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         percent,
		}
		result = append(result, relative)
	}

	if indexData.MOER != "" {
		moer, err := strconv.ParseFloat(indexData.MOER, 64)
		if err != nil {
			return nil, time.Time{}, err
		}
		marginal := CarbonIntensity{
			EmissionsType: MarginalEmissionsType,
			MetricType:    AbsoluteMetricType,
			Provider:      WattTime,
			Location:      location,
			Units:         LbCO2EPerMWh,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         moer,
		}
		result = append(result, marginal)
	}

	return result, ttl, nil
}

type wattTimeIndexData struct {
	BA        string    `json:"ba"`
	Freq      string    `json:"freq"`
	MOER      string    `json:"moer"`
	Percent   string    `json:"percent"`
	PointTime time.Time `json:"point_time"`
}

type wattTimeLoginResp struct {
	Token string `json:"token"`
}
