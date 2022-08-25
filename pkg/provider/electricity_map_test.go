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

var MockElectricityMapResponse = `{
	"zone": "IN-KA",
	"carbonIntensity": 312,
	"datetime": "2020-01-01T00:00:00.000Z",
	"updatedAt": "2020-01-01T00:00:01.000Z"
}`

func Test_ElectricityMap_SimpleRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, MockElectricityMapResponse)
	}))
	defer ts.Close()

	c := ElectricityMapConfig{
		APIURL: ts.URL,
		Token:  "token",
	}
	a, err := NewElectricityMap(c)
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	res, err := a.GetCarbonIntensity(context.Background(), "IN-KA")
	if err != nil {
		t.Fatalf("got error on GetCarbonIntensity: %s", err)
	}

	expected := []CarbonIntensity{
		{
			EmissionsType: "average",
			MetricType:    "absolute",
			Provider:      "ElectricityMap",
			Location:      "IN-KA",
			Units:         "gCO2e per kWh",
			ValidFrom:     time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
			ValidTo:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Value:         312,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("want matching \n %s", cmp.Diff(res, expected))
	}
}
