package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ElectricityMapsClient struct {
	client *http.Client
	apiURL string
	token  string
}

type ElectricityMapsConfig struct {
	Client *http.Client
	APIURL string
	Token  string
}

func NewElectricityMaps(config ElectricityMapsConfig) (Interface, error) {
	if config.Client == nil {
		config.Client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.electricitymap.org/v3"
	}

	c := &ElectricityMapsClient{
		apiURL: config.APIURL,
		client: config.Client,
		token:  config.Token,
	}

	return c, nil
}

func (e *ElectricityMapsClient) GetCarbonIntensity(ctx context.Context, location string) ([]CarbonIntensity, error) {
	intensityURL, err := e.historicIntensityURLWithZone(location)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, intensityURL, nil)
	req.Header.Add("auth-token", e.token)
	if err != nil {
		return nil, err
	}

	log.Printf("calling %s", req.URL)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errBadStatus(resp)
	}

	historyResponse := &electricityMapsHistoryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&historyResponse)
	if err != nil {
		return nil, err
	}

	var carbonIntensityPoints []CarbonIntensity
	var recentDatapoints = NewElectricityMapsDatapoints()

	for _, dataPoint := range historyResponse.History {
		// We get the most recent value (which is usually estimated)
		// and the most recent value which is registered (not estimated)
		// and they will end up in the recentDatapoints variable
		recentDatapoints.update(dataPoint)

	}
	// We need to check if each value exists, because sometimes there are no estimated datapoints
	// and sometimes there are only estimated datapoints (depending on the location):
	if recentDatapoints.estimatedFound {
		estimatedCarbonIntensity, err := toCarbonIntensity(location, recentDatapoints.estimated)
		if err == nil {
			carbonIntensityPoints = append(carbonIntensityPoints, *estimatedCarbonIntensity)
		}
	}
	if recentDatapoints.realFound {
		realCarbonIntensity, err := toCarbonIntensity(location, recentDatapoints.real)
		if err == nil {
			carbonIntensityPoints = append(carbonIntensityPoints, *realCarbonIntensity)
		}
	}

	return carbonIntensityPoints, nil
}

// Helper struct to remove clutter in the calling function
// while finding the latest (and greatest) data points
type electricityMapsDatapoints struct {
	estimated      electricityMapsData
	real           electricityMapsData
	estimatedFound bool
	realFound      bool
}

func NewElectricityMapsDatapoints() electricityMapsDatapoints {
	dataPoints := electricityMapsDatapoints{}
	dataPoints.estimatedFound = false
	dataPoints.realFound = false
	return dataPoints
}

func (m *electricityMapsDatapoints) setEstimated(dataPoint electricityMapsData) {
	m.estimated = dataPoint
	m.estimatedFound = true
}

func (m *electricityMapsDatapoints) setReal(dataPoint electricityMapsData) {
	m.real = dataPoint
	m.realFound = true
}

func (m *electricityMapsDatapoints) update(dataPoint electricityMapsData) error {

	dataPointDateTime, err := stringToTime(dataPoint.DateTime)
	if err != nil {
		return err
	}

	// If the current datapoint is estimated,
	// update if it's the most recent:
	if dataPoint.IsEstimated {
		if !m.estimatedFound {
			m.setEstimated(dataPoint)
		} else {
			estimatedDateTime, err := stringToTime(m.estimated.DateTime)
			if err != nil {
				return err
			}
			if estimatedDateTime.Before(dataPointDateTime) {
				m.setEstimated(dataPoint)
			}
		}
	}

	if !dataPoint.IsEstimated {
		if !m.realFound {
			m.setReal(dataPoint)
		} else {
			realDateTime, err := stringToTime(m.real.DateTime)
			if err != nil {
				return err
			}
			if realDateTime.Before(dataPointDateTime) {
				m.setReal(dataPoint)
			}
		}
	}

	return nil
}

func stringToTime(dateTimeString string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, dateTimeString)
}

func toCarbonIntensity(location string, dataPoint electricityMapsData) (*CarbonIntensity, error) {
	validFrom, err := time.Parse(time.RFC3339Nano, dataPoint.DateTime)
	if err != nil {
		log.Printf("Error parsing datetime %s", dataPoint.DateTime)
		return nil, err
	}
	validTo := validFrom.Add(60 * time.Minute)
	carbonIntensityDataPoint := CarbonIntensity{
		EmissionsType: AverageEmissionsType,
		MetricType:    AbsoluteMetricType,
		Provider:      ElectricityMaps,
		Location:      location,
		Units:         GramsCO2EPerkWh,
		ValidFrom:     validFrom,
		ValidTo:       validTo,
		Value:         dataPoint.CarbonIntensity,
		IsEstimated:   dataPoint.IsEstimated,
	}
	return &carbonIntensityDataPoint, nil
}

func (e *ElectricityMapsClient) historicIntensityURLWithZone(zone string) (string, error) {
	zoneURL := fmt.Sprintf("/carbon-intensity/history?zone=%s", zone)
	return buildURL(e.apiURL, zoneURL)
}

type electricityMapsData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
	DateTime        string  `json:"datetime"`
	UpdatedAt       string  `json:"updatedAt"`
	IsEstimated     bool    `json:"isEstimated"`
}

type electricityMapsHistoryResponse struct {
	Zone    string
	History []electricityMapsData
}
