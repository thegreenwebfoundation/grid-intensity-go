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
		location    string
		result      []CarbonIntensity
		expectedErr string
	}{
		{
			name:     "country exists",
			location: "ESP",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Location:      "ESP",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         193.737,
				},
			},
		},
		{
			name:     "2 char country code",
			location: "ES",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Location:      "ES",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         193.737,
				},
			},
		},
		{
			name:     "lower case country code",
			location: "gbr",
			result: []CarbonIntensity{
				{
					EmissionsType: "average",
					MetricType:    "absolute",
					Provider:      "Ember",
					Location:      "GBR",
					Units:         "gCO2e per kWh",
					ValidFrom:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					ValidTo:       time.Date(2021, 12, 31, 23, 59, 0, 0, time.UTC),
					Value:         268.255,
				},
			},
		},
		{
			name:        "invalid country code",
			location:    "AAA",
			expectedErr: "location \"AAA\" not found",
		},
	}

	p, err := NewEmber()
	if err != nil {
		t.Errorf("Could not make provider: %s", err)
		return
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := p.GetCarbonIntensity(ctx, tc.location)
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
