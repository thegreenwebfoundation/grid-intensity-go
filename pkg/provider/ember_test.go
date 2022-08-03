package provider

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_GetGridIntensityForCountry(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		region      string
		result      []CarbonIntensity
		expectedErr string
	}{
		{
			name:   "country exists",
			region: "ESP",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Region:        "ESP",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         193.737,
				},
			},
		},
		{
			name:   "2 char country code",
			region: "ES",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Region:        "ES",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         193.737,
				},
			},
		},
		{
			name:   "lower case country code",
			region: "gbr",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Region:        "GBR",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         268.255,
				},
			},
		},
		{
			name:        "invalid country code",
			region:      "AAA",
			expectedErr: "region \"AAA\" not found",
		},
	}

	p, err := NewEmber()
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := p.GetCarbonIntensity(ctx, tc.region)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil want non-nil")
			}

			if result != nil && len(result) != len(tc.result) {
				t.Fatalf("expected %d result got %d", len(tc.result), len(result))
			}
			if result != nil && !reflect.DeepEqual(tc.result, result) {
				t.Errorf("want matching \n %s", cmp.Diff(result, tc.result))
			}
			if tc.expectedErr != "" && tc.expectedErr != err.Error() {
				t.Fatalf("expected error %q got %q", tc.expectedErr, err.Error())
			}
		})
	}
}
