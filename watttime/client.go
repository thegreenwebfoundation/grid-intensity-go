package watttime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const ProviderName = "watttime.org"

type ApiOption func(*ApiClient) error

func New(user, password string, opts ...ApiOption) (Provider, error) {
	a := &ApiClient{
		user:     user,
		password: password,
	}

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
		a.apiURL = "https://api2.watttime.org/v2"
	}

	return a, nil
}

type ApiClient struct {
	client   *http.Client
	apiURL   string
	user     string
	password string
	token    string
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	result, err := a.fetchCarbonIntensityData(ctx, region)
	if err != nil {
		return -1, err
	}

	if result.MOER == "" {
		return -1, ErrNoMarginalIntensityPresent
	}
	moer, err := strconv.ParseFloat(result.MOER, 64)
	if err != nil {
		return -1, err
	}

	return moer, nil
}

func (a *ApiClient) GetCarbonIntensityData(ctx context.Context, region string) (*IndexData, error) {
	result, err := a.fetchCarbonIntensityData(ctx, region)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ApiClient) GetRelativeCarbonIntensity(ctx context.Context, region string) (int64, error) {
	result, err := a.fetchCarbonIntensityData(ctx, region)
	if err != nil {
		return -1, err
	}

	if result.Percent == "" {
		return -1, ErrNoRelativeIntensityPresent
	}
	percent, err := strconv.ParseInt(result.Percent, 0, 64)
	if err != nil {
		return -1, err
	}

	return percent, nil
}

func (a *ApiClient) fetchCarbonIntensityData(ctx context.Context, region string) (*IndexData, error) {
	if a.token == "" {
		token, err := a.getAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		a.token = token
	}

	result, err := a.getCarbonIntensityData(ctx, region)
	if errors.Is(err, ErrReceived403Forbidden) {
		token, err := a.getAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		a.token = token

		result, err = a.getCarbonIntensityData(ctx, region)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *ApiClient) getAccessToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.loginURL(), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(a.user, a.password)
	resp, err := a.client.Do(req)
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

	loginResp := LoginResp{}
	err = json.Unmarshal(bytes, &loginResp)
	if err != nil {
		return "", nil
	}

	return loginResp.Token, nil
}

func (a *ApiClient) getCarbonIntensityData(ctx context.Context, region string) (*IndexData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.indexURL(region), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	resp, err := a.client.Do(req)
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

	indexData := IndexData{}
	err = json.Unmarshal(bytes, &indexData)
	if err != nil {
		return nil, err
	}

	return &indexData, nil
}

func (a *ApiClient) indexURL(region string) string {
	return fmt.Sprintf("%s/index?ba=%s", a.apiURL, region)
}

func (a *ApiClient) loginURL() string {
	return fmt.Sprintf("%s/login", a.apiURL)
}
