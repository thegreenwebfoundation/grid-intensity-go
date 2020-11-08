package gridintensity

import "context"

type Provider interface {
	GetCarbonIndex(ctx context.Context, region string) (CarbonIndex, error)
	GetCarbonIntensity(ctx context.Context, region string) (float64, error)
	MinIntensity(ctx context.Context, regions ...string) (region string, err error)
}
