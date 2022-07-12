package watttime

import (
	"context"
	"time"
)

type Provider interface {
	GetCarbonIntensity(ctx context.Context, region string) (float64, error)
	GetCarbonIntensityData(ctx context.Context, region string) (*IndexData, error)
	GetRelativeCarbonIntensity(ctx context.Context, region string) (int64, error)
}

type CacheData struct {
	Data *IndexData `json:"data"`
	TTL  time.Time  `json:"ttl"`
}

type IndexData struct {
	BA        string    `json:"ba"`
	Freq      string    `json:"freq"`
	MOER      string    `json:"moer"`
	Percent   string    `json:"percent"`
	PointTime time.Time `json:"point_time"`
}

type LoginResp struct {
	Token string `json:"token"`
}
