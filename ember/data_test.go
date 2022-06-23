package ember

import (
	"reflect"
	"testing"
)

func Test_GetGridIntensityForCountry(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		result      *GridIntensity
		expectedErr string
	}{
		{
			name:        "country exists",
			countryCode: "ESP",
			result: &GridIntensity{
				CountryCodeISO2:              "ES",
				CountryCodeISO3:              "ESP",
				CountryOrRegion:              "Spain",
				Year:                         2021,
				LatestYear:                   2021,
				EmissionsIntensityGCO2PerKWH: 193.737,
			},
		},
		{
			name:        "2 char country code",
			countryCode: "ES",
			result: &GridIntensity{
				CountryCodeISO2:              "ES",
				CountryCodeISO3:              "ESP",
				CountryOrRegion:              "Spain",
				Year:                         2021,
				LatestYear:                   2021,
				EmissionsIntensityGCO2PerKWH: 193.737,
			},
		},
		{
			name:        "lower case country code",
			countryCode: "gbr",
			result: &GridIntensity{
				CountryCodeISO2:              "GB",
				CountryCodeISO3:              "GBR",
				CountryOrRegion:              "United Kingdom",
				Year:                         2021,
				LatestYear:                   2021,
				EmissionsIntensityGCO2PerKWH: 268.255,
			},
		},
		{
			name:        "invalid country code",
			countryCode: "AAA",
			expectedErr: "country code \"AAA\" not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetGridIntensityForCountry(tc.countryCode)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil want non-nil")
			}

			if !reflect.DeepEqual(tc.result, result) {
				t.Fatalf("expected %#v got %#v", tc.result, result)
			}
			if tc.expectedErr != "" && tc.expectedErr != err.Error() {
				t.Fatalf("expected error %q got %q", tc.expectedErr, err.Error())
			}
		})
	}
}
