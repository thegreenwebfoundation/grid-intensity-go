//go:build dockerrequired
// +build dockerrequired

package exporter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func Test_GridIntensityMetric(t *testing.T) {
	var metrics string

	metricsURL := "http://localhost:8000/metrics"

	o := func() error {
		resp, err := http.Get(metricsURL)
		if err != nil {
			return fmt.Errorf("could not retrieve %s: %v", metricsURL, err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected status code %d: got %d", http.StatusOK, resp.StatusCode)
		}

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read body %v", err)
		}

		metrics = string(respBytes)

		return nil
	}

	n := func(err error, d time.Duration) {
		t.Logf("failed to get metrics from %s: retrying in %s", metricsURL, d)
	}

	b := backoff.NewExponentialBackOff()
	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	platform := os.Getenv("TEST_PLATFORM")

	var node, region string

	if platform == "kubernetes" {
		node = "chart-testing-control-plane"
	} else if platform == "nomad" {
		node = "dev01"
		region = "global"
	}

	expectedMetricText := fmt.Sprintf("grid_intensity_carbon_average{is_estimated=\"true\",location=\"GBR\",node=\"%s\",provider=\"Ember\",region=\"%s\",units=\"gCO2e per kWh\"}",
		node, region)

	if !strings.Contains(metrics, expectedMetricText) {
		t.Fatalf("expected metric text %q not found got %q", expectedMetricText, metrics)
	}
}
