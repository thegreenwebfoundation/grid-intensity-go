package carbonintensity

type CarbonIntensityResponse struct {
	Data []IntensityData `json:"data"`
}

type IntensityData struct {
	From      string     `json:"from"`
	To        string     `json:"to"`
	Intensity *Intensity `json:"intensity"`
}

type Intensity struct {
	Forecast float64 `json:"forecast"`
	Actual   float64 `json:"actual"`
	Index    string  `json:"index"`
}
