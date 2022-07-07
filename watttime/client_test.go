package watttime_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/thegreenwebfoundation/grid-intensity-go/watttime"
)

var MockIndexResponse = `{
		"ba": "CAISO_NORTH",
		"freq": "300",
		"moer": "916",
		"percent": "78",
		"point_time": "2022-07-06T16:25:00Z"
}`
var MockLoginResponse = `{"token":"mytoken"}`

func makeTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			fmt.Fprintln(w, MockLoginResponse)
		case "/index":
			fmt.Fprintln(w, MockIndexResponse)
		default:
			t.Errorf("unknown path %#q", r.URL.Path)
		}
	}))
}

func Test_GetCarbonIntensity(t *testing.T) {
	ts := makeTestServer(t)
	defer ts.Close()

	c, err := watttime.New("user", "password", watttime.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	result, err := c.GetCarbonIntensity(context.Background(), "CAISO_NORTH")
	if err != nil {
		t.Errorf("Got error on GetCarbonIntensity: %s", err)
		return
	}

	moer := 916.0
	if result != moer {
		t.Errorf("Expected %.2f, got %.2f", moer, result)
	}
}

func Test_GetCarbonIntensityData(t *testing.T) {
	ts := makeTestServer(t)
	defer ts.Close()

	c, err := watttime.New("user", "password", watttime.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	result, err := c.GetCarbonIntensityData(context.Background(), "CAISO_NORTH")
	if err != nil {
		t.Errorf("Got error on GetCarbonIntensity: %s", err)
		return
	}

	data := &watttime.IndexData{
		BA:        "CAISO_NORTH",
		Freq:      "300",
		MOER:      "916",
		Percent:   "78",
		PointTime: time.Date(2022, 7, 6, 16, 25, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, data) {
		t.Errorf("Expected %#v, got %#v", data, result)
	}
}

func Test_GetRelativeCarbonIntensity(t *testing.T) {
	ts := makeTestServer(t)
	defer ts.Close()

	c, err := watttime.New("user", "password", watttime.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	result, err := c.GetRelativeCarbonIntensity(context.Background(), "CAISO_NORTH")
	if err != nil {
		t.Errorf("Got error on GetCarbonIntensity: %s", err)
		return
	}

	percent := 78
	if result != int64(percent) {
		t.Errorf("Expected %d, got %d", percent, result)
	}
}
