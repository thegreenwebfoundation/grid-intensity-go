package data

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"strconv"
)

//go:embed co2-intensities-ember-2021.csv
var emberData []byte

func GetEmberGridIntensity() (map[string]EmberGridIntensity, error) {
	data := map[string]EmberGridIntensity{}

	reader := bytes.NewReader(emberData)
	rows, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		countryCodeISO2 := row[0]
		countryCodeISO3 := row[1]
		if countryCodeISO2 == "" || countryCodeISO2 == "country_code_iso_2" {
			continue
		}

		year, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, err
		}

		latestYear, err := strconv.Atoi(row[4])
		if err != nil {
			return nil, err
		}

		intensity, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, err
		}

		country := EmberGridIntensity{
			CountryCodeISO2:              countryCodeISO2,
			CountryCodeISO3:              countryCodeISO3,
			CountryOrRegion:              row[2],
			Year:                         year,
			LatestYear:                   latestYear,
			EmissionsIntensityGCO2PerKWH: intensity,
		}

		// Add both country codes to allow lookups with either format.
		data[countryCodeISO2] = country
		data[countryCodeISO3] = country
	}

	return data, nil
}
