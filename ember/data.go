package ember

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go/api"
)

//go:embed co2-intensities-ember-2021.csv
var emberData []byte

type ApiClient struct {
	data map[string]GridIntensity
}

func New() (gridintensity.Provider, error) {
	emberData, err := getGridIntensityData()
	if err != nil {
		return nil, err
	}

	return &ApiClient{data: emberData}, nil
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	result, ok := a.data[strings.ToUpper(region)]
	if !ok {
		return 0, fmt.Errorf("region %q not found", region)
	}

	return result.EmissionsIntensityGCO2PerKWH, nil
}

func GetGridIntensityForCountry(countryCode string) (*GridIntensity, error) {
	emberData, err := getGridIntensityData()
	if err != nil {
		return nil, err
	}

	result, ok := emberData[strings.ToUpper(countryCode)]
	if !ok {
		return nil, fmt.Errorf("country code %q not found", countryCode)
	}

	return &result, nil
}

func getGridIntensityData() (map[string]GridIntensity, error) {
	data := map[string]GridIntensity{}

	reader := bytes.NewReader(emberData)
	rows, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		countryCode := row[0]
		if countryCode == "" || countryCode == "country_code" {
			continue
		}

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

		data[countryCode] = GridIntensity{
			CountryCode:                  countryCode,
			CountryOrRegion:              row[1],
			Year:                         year,
			LatestYear:                   latestYear,
			EmissionsIntensityGCO2PerKWH: intensity,
		}
	}

	return data, nil
}
