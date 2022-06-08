package api

import (
	"context"
)

type Provider interface {
	GetCarbonIntensity(ctx context.Context, region string) (float64, error)
}
