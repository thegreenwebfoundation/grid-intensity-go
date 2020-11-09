package energymap_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
	"github.com/thegreenwebfoundation/grid-intensity-go/energymap"
)

var responseTable = map[string]float64{
	"IN-KA": 312,
	"IN-AP": 301,
}

var MockResponses = map[string]string{
	"IN-KA": `{
		"zone": "IN-KA",
		"carbonIntensity": 312,
		"datetime": "2020-01-01T00:00:00.000Z",
		"updatedAt": "2020-01-01T00:00:01.000Z"
	}`,
	"IN-AP": `{
			"zone": "IN-AP",
			"carbonIntensity": 301,
			"datetime": "2020-01-01T00:00:00.000Z",
			"updatedAt": "2020-01-01T00:00:01.000Z"
	}`,
}

func makeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zone := r.URL.Query().Get("zone")
		i := time.Duration(100 + rand.Intn(200))
		time.Sleep(i * time.Millisecond)
		if val, ok := MockResponses[zone]; ok {
			fmt.Fprintln(w, val)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error": "unknown zone"}`)
	}))
}

func TestSimpleRequest(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	e, err := energymap.New("fake_token", energymap.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	resp, err := e.GetCarbonIntensity(context.Background(), "IN-KA")
	if err != nil {
		t.Errorf("Got error on GetCarbonIndex: %s", err)
		return
	}

	if resp != 312 {
		t.Errorf("Expected HIGH, got %.2f", resp)
	}
}

func TestMinIntensity(t *testing.T) {
	ts := makeTestServer()
	defer ts.Close()

	e, err := energymap.New("fake_token", energymap.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	resp, err := gridintensity.GetCarbonIntensityMap(context.Background(), e, "IN-KA", "IN-AP")
	if err != nil {
		t.Errorf("Got error on GetCarbonIndex: %s", err)
		return
	}
	carbonIntensityMap := resp.GetAll()

	for region, value := range responseTable {
		if _, ok := carbonIntensityMap[region]; !ok {
			t.Errorf("Expected to find %s in the map", region)
			return
		}
		if carbonIntensityMap[region] != value {
			t.Errorf("Expected region %s to have %.2f gco2e/kwh intensity, got %.2f", region, value, carbonIntensityMap[region])
			return
		}
	}

}
