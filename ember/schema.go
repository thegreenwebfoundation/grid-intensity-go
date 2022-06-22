package ember

type GridIntensity struct {
	CountryCode                  string  `json:"country_code"`
	CountryOrRegion              string  `json:"country_or_region"`
	Year                         int     `json:"year"`
	LatestYear                   int     `json:"latest_year"`
	EmissionsIntensityGCO2PerKWH float64 `json:"emissions_intensity_gco2_per_kwh"`
}
