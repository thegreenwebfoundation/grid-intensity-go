package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var MockGridCarbonResponse = `{"unit":"gCO2eq/kWh","data":[{"zone":"FR","ts":"2026-09-03T05:00:00Z","gco2eq_kwh":32.2,"method":"computed:v1"}]}`

var MockGridCarbonGBResponse = `{"unit":"gCO2eq/kWh","data":[{"zone":"GB","ts":"2026-09-03T06:30:00Z","gco2eq_kwh":82,"method":"upstream:uk-neso:forecast"}]}`

func Test_GridCarbon_SimpleRequest(t *testing.T) {
	var gotPath, gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.RequestURI()
		gotUA = r.Header.Get("User-Agent")
		fmt.Fprintln(w, MockGridCarbonResponse)
	}))
	defer ts.Close()

	a, err := NewGridCarbon(GridCarbonConfig{APIURL: ts.URL + "/"})
	if err != nil {
		t.Fatalf("Could not make provider: %s", err)
	}

	res, err := a.GetCarbonIntensity(context.Background(), "fr")
	if err != nil {
		t.Fatalf("got error on GetCarbonIntensity: %s", err)
	}

	if gotPath != "/v1/intensity/latest?zone=FR" {
		t.Errorf("want path /v1/intensity/latest?zone=FR, got %q", gotPath)
	}
	if gotUA != gridCarbonUserAgent {
		t.Errorf("want user agent %q, got %q", gridCarbonUserAgent, gotUA)
	}

	expected := []CarbonIntensity{
		{
			Provider:      "GridCarbon",
			EmissionsType: "average",
			MetricType:    "absolute",
			Location:      "FR",
			Units:         "gCO2e per kWh",
			ValidFrom:     time.Date(2026, 9, 3, 5, 0, 0, 0, time.UTC),
			ValidTo:       time.Date(2026, 9, 3, 6, 0, 0, 0, time.UTC),
			Value:         32.2,
			IsEstimated:   false,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("want matching \n %s", cmp.Diff(res, expected))
	}
}

func Test_GridCarbon_GBHalfHourForecast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, MockGridCarbonGBResponse)
	}))
	defer ts.Close()

	a, _ := NewGridCarbon(GridCarbonConfig{APIURL: ts.URL})
	res, err := a.GetCarbonIntensity(context.Background(), "GB")
	if err != nil {
		t.Fatalf("got error on GetCarbonIntensity: %s", err)
	}

	expected := []CarbonIntensity{
		{
			Provider:      "GridCarbon",
			EmissionsType: "average",
			MetricType:    "absolute",
			Location:      "GB",
			Units:         "gCO2e per kWh",
			ValidFrom:     time.Date(2026, 9, 3, 6, 30, 0, 0, time.UTC),
			ValidTo:       time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC),
			Value:         82,
			IsEstimated:   true,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("want matching \n %s", cmp.Diff(res, expected))
	}
}

func Test_GridCarbon_UnknownZone(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"unknown or empty zone: XX"}`)
	}))
	defer ts.Close()

	a, _ := NewGridCarbon(GridCarbonConfig{APIURL: ts.URL})
	_, err := a.GetCarbonIntensity(context.Background(), "XX")
	if !errors.Is(err, ErrInvalidLocation) {
		t.Errorf("want ErrInvalidLocation, got %v", err)
	}

	_, err = a.GetCarbonIntensity(context.Background(), "")
	if !errors.Is(err, ErrInvalidLocation) {
		t.Errorf("want ErrInvalidLocation for empty location, got %v", err)
	}
}

func Test_GridCarbon_EmptyData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"unit":"gCO2eq/kWh","data":[]}`)
	}))
	defer ts.Close()

	a, _ := NewGridCarbon(GridCarbonConfig{APIURL: ts.URL})
	_, err := a.GetCarbonIntensity(context.Background(), "FR")
	if !errors.Is(err, ErrNoResponse) {
		t.Errorf("want ErrNoResponse, got %v", err)
	}
}

func Test_GridCarbon_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, `{"error":"upstream stale"}`)
	}))
	defer ts.Close()

	a, _ := NewGridCarbon(GridCarbonConfig{APIURL: ts.URL})
	_, err := a.GetCarbonIntensity(context.Background(), "FR")
	if !errors.Is(err, ErrReceivedNon200Status) {
		t.Errorf("want ErrReceivedNon200Status, got %v", err)
	}
}
