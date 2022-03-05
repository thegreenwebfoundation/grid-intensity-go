package unfccc_test

import (
	"context"
	"testing"

	"github.com/thegreenwebfoundation/grid-intensity-go/unfccc"
)

func TestFetchIntensity(t *testing.T) {
	u, err := unfccc.New("sampleRegionData.json")
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	resp, err := u.GetCarbonIntensity(context.Background(), "de")
	if err != nil {
		t.Errorf("Could not get carbon intensity: %s", err)
	}

	if resp != 312 {
		t.Errorf("Expected 312, got %.2f", resp)
	}
}
