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

	resp, err := e.GetCarbonIndex(context.Background(), "IN-KA")
	if err != nil {
		t.Errorf("Got error on GetCarbonIndex: %s", err)
		return
	}

	if resp != gridintensity.HIGH {
		t.Errorf("Expected HIGH, got %s", resp)
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

	resp, err := e.MinIntensity(context.Background(), "IN-KA", "IN-AP")
	if err != nil {
		t.Errorf("Got error on GetCarbonIndex: %s", err)
		return
	}

	if resp != "IN-AP" {
		t.Errorf("Expected IN-AP, got %s", resp)
	}
}
