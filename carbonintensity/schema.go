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
	Forecast int    `json:"forecast"`
	Actual   int    `json:"actual"`
	Index    string `json:"index"`
}
