package data

type EmberGridIntensity struct {
	CountryCodeISO2              string  `json:"country_code_iso_2"`
	CountryCodeISO3              string  `json:"country_code_iso_3"`
	CountryOrRegion              string  `json:"country_or_region"`
	Year                         int     `json:"year"`
	LatestYear                   int     `json:"latest_year"`
	EmissionsIntensityGCO2PerKWH float64 `json:"emissions_intensity_gco2_per_kwh"`
}
