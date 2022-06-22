package api

import (
	"context"
)

type Provider interface {
	GetAllRegionsCarbonIntensity(ctx context.Context) (map[string]float64, error)
	GetCarbonIntensity(ctx context.Context, region string) (float64, error)
}
