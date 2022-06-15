package data

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type GridIntensity struct {
	CountryCode                  string  `json:"country_code"`
	CountryOrRegion              string  `json:"country_or_region"`
	Year                         int     `json:"year"`
	LatestYear                   int     `json:"latest_year"`
	EmissionsIntensityGCO2PerKWH float64 `json:"emissions_intensity_gco2_per_kwh"`
}

//go:embed co2-intensities-ember-2021.csv
var data []byte

func GetGridIntensityForCountry(countryCode string) (*GridIntensity, error) {
	reader := bytes.NewReader(data)
	rows, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if strings.EqualFold(row[0], countryCode) {
			year, err := strconv.Atoi(row[2])
			if err != nil {
				return nil, err
			}

			latestYear, err := strconv.Atoi(row[3])
			if err != nil {
				return nil, err
			}

			intensity, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				return nil, err
			}

			return &GridIntensity{
				CountryCode:                  row[0],
				CountryOrRegion:              row[1],
				Year:                         year,
				LatestYear:                   latestYear,
				EmissionsIntensityGCO2PerKWH: intensity,
			}, nil
		}
	}

	return nil, fmt.Errorf("country code %q not found", countryCode)
}
