package unfccc

// CarbonIntensityReading a representation of the annual carbon intensity for the given country,
// as identified by region.
// Region is usually a 2 character country code, but with special entries for global or EU wide figures.
type CarbonIntensityReading struct {
	Region          string  `json:"region"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}
