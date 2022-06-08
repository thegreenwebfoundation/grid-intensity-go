package api

import (
	"context"
	"sync"
)

func GetCarbonIntensityMap(ctx context.Context, p Provider, regions ...string) (*CarbonMap, error) {

	if len(regions) == 0 {
		return nil, ErrNoRegionProvided
	}

	requestCounter := len(regions)

	intensityMap := &CarbonMap{
		m: make(map[string]float64, requestCounter),
	}
	errChan := make(chan error, requestCounter)

	for _, region := range regions {
		go func(r string) {
			intensity, err := p.GetCarbonIntensity(ctx, r)
			errChan <- err
			if err != nil {
				return
			}
			intensityMap.Set(r, intensity)
		}(region)
	}

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
			requestCounter--
			if requestCounter == 0 {
				return intensityMap, nil
			}
		case <-ctx.Done():
			return nil, ErrTimeout
		}
	}
}

type CarbonMap struct {
	m  map[string]float64
	mu sync.RWMutex
}

func (c *CarbonMap) Set(region string, intensity float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[region] = intensity
}

func (c *CarbonMap) GetAll() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.m
}
