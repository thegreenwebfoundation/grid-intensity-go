package watttime

import "net/http"

func WithHTTPClient(h *http.Client) ApiOption {
	return func(a *ApiClient) error {
		a.client = h
		return nil
	}
}

func WithAPIURL(url string) ApiOption {
	return func(a *ApiClient) error {
		a.apiURL = url
		return nil
	}
}

func WithCacheFile(cacheFile string) ApiOption {
	return func(a *ApiClient) error {
		a.cacheFile = cacheFile
		return nil
	}
}
