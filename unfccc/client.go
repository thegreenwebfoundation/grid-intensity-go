package unfccc

import (
	"context"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
)

type ApiOption func(*ApiClient) error

// New returns an instance of the client, with the local map populated with
// carbon intensity values
func New(opts ...ApiOption) (gridintensity.Provider, error) {

	a := &ApiClient{}

	// load JSON file

	// populate our map with the intensities

	return a, nil
}

type ApiClient struct {
	sourceFilePath string
	regionMap      map[string]float64
}

func (a *ApiClient) GetCarbonIntensity(ctx context.Context, region string) (float64, error) {
	var ok bool

	intensity, ok := a.regionMap[region]; ok

	if ok == false {
		return (intensity, nil)
	}

	return (0, ErrNoMatchingRegion)

}
