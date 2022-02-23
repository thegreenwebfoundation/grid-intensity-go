package unfccc_test

import (
	"context"
	"testing"

	"github.com/thegreenwebfoundation/grid-intensity-go/unfccc"
)

// dummy values until we have a JSON file we can load in
var responseTable = map[string]float64{
	"DE": 312,
}

// sample intensity for Germany
func TestFetchIntensity(t *testing.T) {

	u, err := unfccc.New()
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	resp, err := u.GetCarbonIntensity(context.Background(), "DE")

	if resp != 312 {
		t.Errorf("Expected 312, got %.2f", resp)
	}

}
