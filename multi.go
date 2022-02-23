package gridintensity

import (
	"context"
	"sync"
)

// GetCarbonIntensityMap accept a Provider, and a series of regions
func GetCarbonIntensityMap(ctx context.Context, p Provider, regions ...string) (*CarbonMap, error) {

	if len(regions) == 0 {
		return nil, ErrNoRegionProvided
	}

	requestCounter := len(regions)

	// create our Map of carbon intesity values
	// for each region we are interested in
	intensityMap := &CarbonMap{
		m: make(map[string]float64, requestCounter),
	}
	errChan := make(chan error, requestCounter)

	// for each of the regions, fetch the carbon intensity with its
	// own channel, and store the result in  our CarbinIntensityMap
	for _, region := range regions {

		// declare our goroutine. then call it right away
		go func(r string) {
			intensity, err := p.GetCarbonIntensity(ctx, r)

			// sent the result back to the errChan.
			// an empty result means we decrement
			// our counter
			errChan <- err
			if err != nil {
				return
			}
			intensityMap.Set(r, intensity)
		}(region)
	}

	// how does this know to stop? - I can'
	for {
		select {
		// keep reading from the errChan channel,
		// until there are no more goroutines to
		// finish running, and return the populated
		// intensityMap
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
			requestCounter--
			if requestCounter == 0 {
				return intensityMap, nil
			}
		//
		case <-ctx.Done():
			return nil, ErrTimeout
		}
	}
}

type CarbonMap struct {
	m  map[string]float64
	mu sync.RWMutex
}

// Set the value for region to intensity on Carbon Map.
// Uses a mutex to avoid  multiple requests setting the same memory
func (c *CarbonMap) Set(region string, intensity float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[region] = intensity
}

// GetAll return the populated map, locking it to avoid
// other requests accessing the same memory
func (c *CarbonMap) GetAll() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.m
}
