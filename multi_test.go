package gridintensity_test

import (
	"context"
	"errors"
	"testing"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

var responseTable = map[string]float64{
	"IN-KA": 312.56,
	"IN-AP": 301.23,
}

type MockProvider struct {
}

func (m *MockProvider) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	if val, ok := responseTable[region]; ok {
		return val, nil
	}
	return 0, errors.New("region not found")
}

func TestGetCarbonIntensityMap(t *testing.T) {
	p := &MockProvider{}

	resp, err := gridintensity.GetCarbonIntensityMap(context.Background(), p, "IN-KA", "IN-AP")
	if err != nil {
		t.Errorf("could not get carbonIntensityMap: %s", err)
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
