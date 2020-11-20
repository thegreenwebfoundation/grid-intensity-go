package carbonintensity

type CarbonIntensityResponse struct {
	Data []IntensityData `json:"data"`
}

type IntensityData struct {
	Intensity *Intensity `json:"intensity"`
}

type Intensity struct {
	Forecast float64 `json:"forecast"`
	Actual   float64 `json:"actual"`
	Index    string  `json:"index"`
}
