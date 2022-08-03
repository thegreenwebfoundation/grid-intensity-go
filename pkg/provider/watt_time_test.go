package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var MockWattTimeIndexResponse = `{
		"ba": "CAISO_NORTH",
		"freq": "300",
		"moer": "916",
		"percent": "78",
		"point_time": "2022-07-06T16:25:00Z"
}`
var MockWattTimeLoginResponse = `{"token":"mytoken"}`

func makeWattTimeTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			fmt.Fprintln(w, MockWattTimeLoginResponse)
		case "/index":
			fmt.Fprintln(w, MockWattTimeIndexResponse)
		default:
			t.Errorf("unknown path %#q", r.URL.Path)
		}
	}))
}

func Test_WattTime_SimpleRequest(t *testing.T) {
	ts := makeWattTimeTestServer(t)
	defer ts.Close()

	c := WattTimeConfig{
		APIURL:      ts.URL,
		APIUser:     "user",
		APIPassword: "password",
	}
	w, err := NewWattTime(c)
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	result, err := w.GetCarbonIntensity(context.Background(), "CAISO_NORTH")
	if err != nil {
		t.Errorf("Got error on GetCarbonIntensity: %s", err)
		return
	}

	expected := []CarbonIntensity{
		{
			EmissionsType: "marginal",
			MetricType:    "relative",
			Provider:      "watttime.org",
			Region:        "CAISO_NORTH",
			Units:         "percent",
			ValidFrom:     time.Date(2022, 7, 6, 16, 25, 0, 0, time.UTC),
			ValidTo:       time.Date(2022, 7, 6, 16, 30, 0, 0, time.UTC),
			Value:         78,
		},
		{
			EmissionsType: "marginal",
			MetricType:    "absolute",
			Provider:      "watttime.org",
			Region:        "CAISO_NORTH",
			Units:         "lbCO2e per MWh",
			ValidFrom:     time.Date(2022, 7, 6, 16, 25, 0, 0, time.UTC),
			ValidTo:       time.Date(2022, 7, 6, 16, 30, 0, 0, time.UTC),
			Value:         916,
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("want matching \n %s", cmp.Diff(result, expected))
	}
}
