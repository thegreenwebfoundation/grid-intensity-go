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

var MockCarbonIntensityOrgUKResponse = `{
    "data": [
        {
            "from": "2020-01-01T00:00Z",
            "to": "2020-01-01T00:30Z",
            "intensity": {
                "forecast": 186,
                "actual": 190,
                "index": "moderate"
            }
        }
    ]
}`

func Test_CarbonIntensityUK_SimpleRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, MockCarbonIntensityOrgUKResponse)
	}))
	defer ts.Close()

	c := CarbonIntensityUKConfig{
		APIURL: ts.URL,
	}
	a, err := NewCarbonIntensityUK(c)
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	res, err := a.GetCarbonIntensity(context.Background(), "UK")
	if err != nil {
		t.Fatalf("got error on GetCarbonIntensity: %s", err)
	}

	expected := []CarbonIntensity{
		{
			DataProvider:  "CarbonIntensityOrgUK",
			EmissionsType: "average",
			MetricType:    "absolute",
			Region:        "UK",
			Units:         "gCO2e per kWh",
			ValidFrom:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			ValidTo:       time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
			Value:         190,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("want matching \n %s", cmp.Diff(res, expected))
	}
}
