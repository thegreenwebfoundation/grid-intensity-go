package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/internal/data"
)

const (
	emberDataYear = 2021
)

type EmberClient struct {
	data map[string]data.EmberGridIntensity
}

func NewEmber() (Interface, error) {
	data, err := data.GetEmberGridIntensity()
	if err != nil {
		return nil, err
	}

	c := &EmberClient{
		data: data,
	}

	return c, nil
}

func (a *EmberClient) GetCarbonIntensity(ctx context.Context, region string) ([]CarbonIntensity, error) {
	region = strings.ToUpper(region)
	result, ok := a.data[region]
	if !ok {
		return nil, fmt.Errorf("region %q not found", region)
	}

	validFrom := time.Date(emberDataYear, 1, 1, 0, 0, 0, 0, time.UTC)
	validTo := time.Date(emberDataYear, 12, 31, 23, 59, 0, 0, time.UTC)

	return []CarbonIntensity{
		{
			EmissionsType: AverageEmissionsType,
			MetricType:    AbsoluteMetricType,
			Provider:      Ember,
			Region:        region,
			Units:         GramsCO2EPerkWh,
			ValidFrom:     validFrom,
			ValidTo:       validTo,
			Value:         result.EmissionsIntensityGCO2PerKWH,
		},
	}, nil
}
