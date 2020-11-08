package energymap

import "sync"

type CarbonIntensityResp struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}

type IntensityMap struct {
	m  map[string]float64
	mu sync.Mutex
}

func (i *IntensityMap) Set(region string, intensity float64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.m[region] = intensity
}

func (i *IntensityMap) Min() (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.m) == 0 {
		return "", ErrNoRegionProvided
	}
	var minRegion string = ""
	// start with an impossibly high value
	var minIntensity float64 = 999999.0
	for region, intensity := range i.m {
		if intensity < minIntensity {
			minRegion = region
			minIntensity = intensity
		}
	}

	if minRegion == "" {
		return "", ErrNoRegionFound
	}
	return minRegion, nil
}
