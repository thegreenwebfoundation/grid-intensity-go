package unfccc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// New returns an instance of the client, with the local map populated with
// carbon intensity values
func New(pathtoJSONFile string) (gridintensity.Provider, error) {

	// the map to hold our carbon intensity entries
	var regionMap map[string]CarbonIntensityReading

	// the file to the serialised region carbon intensity data
	var regionJSON []byte

	a := &ApiClient{}

	// figure out how to load JSON file
	regionJSON, err := os.ReadFile(pathtoJSONFile)
	if err != nil {
		panic(err)
	}

	// populate our map with the intensities
	err = json.Unmarshal(regionJSON, &regionMap)
	if err != nil {
		panic(err)
	}

	a.regionMap = regionMap

	fmt.Println(regionMap)

	// assign our regionMap to the "API client". Somehow.

	return a, nil
}

type ApiClient struct {
	sourceFilePath string
	regionMap      map[string]CarbonIntensityReading
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {

	fmt.Println(a.regionMap)

	if intensity, ok := a.regionMap[region]; ok {
		return intensity.CarbonIntensity, nil
	}

	return 0, ErrNoMatchingRegion

}
