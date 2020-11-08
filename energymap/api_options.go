package energymap

import "net/http"

func WithHTTPClient(h *http.Client) ApiOption {
	return func(a *ApiClient) error {
		a.client = h
		return nil
	}
}

func WithThresholds(moderate, high float64) ApiOption {
	return func(a *ApiClient) error {
		a.moderate_threshold = moderate
		a.high_threshold = high
		return nil
	}
}

func WithAPIURL(url string) ApiOption {
	return func(a *ApiClient) error {
		a.apiURL = url
		return nil
	}
}
