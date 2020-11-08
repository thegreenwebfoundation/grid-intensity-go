package carbonintensity

type CarbonIntensityResponse struct {
	Data []IntensityData `json:"data"`
}

type IntensityData struct {
	//From      time.Time `json:"from"`
	//To        time.Time `json:"to"`
	Intensity *Intensity `json:"intensity"`
}

type Intensity struct {
	Forecast float64 `json:"forecast"`
	Actual   float64 `json:"actual"`
	Index    string  `json:"index"`
}
