package unfccc

import (
	"context"
	"encoding/json"
	"io/ioutil"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

// New returns an instance of the client, with the local map populated with
// carbon intensity values
func New(pathtoJSONFile string) (gridintensity.Provider, error) {
	a := &ApiClient{}

	// the map to hold our carbon intensity entries
	var regionMap map[string]CarbonIntensityReading

	var regionJSON []byte
	regionJSON, err := ioutil.ReadFile(pathtoJSONFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(regionJSON, &regionMap)
	if err != nil {
		return nil, err
	}

	a.regionMap = regionMap

	return a, nil
}

type ApiClient struct {
	regionMap map[string]CarbonIntensityReading
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	if intensity, ok := a.regionMap[region]; ok {
		return intensity.CarbonIntensity, nil
	}

	return 0, ErrNoMatchingRegion
}
