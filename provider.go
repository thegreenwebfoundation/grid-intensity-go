package gridintensity

import "context"

type Provider interface {
	GetCarbonIndex(ctx context.Context, region string) (CarbonIndex, error)
	MinIntensity(ctx context.Context, regions ...string) (region string, err error)
}
