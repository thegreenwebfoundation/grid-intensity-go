package carbonintensity_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
	"github.com/thegreenwebfoundation/grid-intensity-go/carbonintensity"
)

var MockResponse = `{
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

func TestSimpleRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, MockResponse)
	}))
	defer ts.Close()

	c, err := carbonintensity.New(carbonintensity.WithAPIURL(ts.URL))
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	resp, err := c.GetCarbonIndex(context.Background(), "UK")
	if err != nil {
		t.Errorf("Got error on GetCarbonIndex: %s", err)
		return
	}

	if resp != gridintensity.MODERATE {
		t.Errorf("Expected MODERATE, got %s", resp)
	}
}
