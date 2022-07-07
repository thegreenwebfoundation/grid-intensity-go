package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

const (
	countryCode     = "country_code"
	countryCodeISO2 = "country_code_iso_2"
	countryCodeISO3 = "country_code_iso_3"
)

func main() {
	err := mapCountryCodes(os.Args[1])
	if err != nil {
		panic(err)
	}
}

// mapCountryCodes takes in an Ember data file with country_code and 3 char
// ISO codes. The data is transformed to have both 2 and 3 char ISO codes so
// users of the CLI can use either format.
func mapCountryCodes(inputFile string) error {
	countries, err := getCountryLookups()
	if err != nil {
		return err
	}

	data, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	rows, err := csv.NewReader(data).ReadAll()
	if err != nil {
		return err
	}

	err = updateHeader(rows[0])
	if err != nil {
		return err
	}

	// Skip header row as it is output by updateHeader.
	rows = rows[1:]

	for _, row := range rows {
		iso_3 := row[0]
		iso_2, ok := countries[iso_3]
		if !ok {
			if iso_3 != "" {
				return fmt.Errorf("country %#q not found", iso_3)
			}
		}

		output := []string{
			iso_2,
		}
		output = append(output, row...)

		fmt.Println(strings.Join(output, ","))
	}

	return nil
}

func getCountryLookups() (map[string]string, error) {
	countries := map[string]string{}

	data, err := os.Open("hack/countries.csv")
	if err != nil {
		return nil, err
	}

	rows, err := csv.NewReader(data).ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		iso_2 := row[1]
		iso_3 := row[2]
		countries[iso_3] = iso_2
	}

	return countries, nil
}

func updateHeader(header []string) error {
	if header[0] == countryCodeISO2 {
		return fmt.Errorf("data already processed - %#q should not be present", countryCodeISO2)
	} else if header[0] == countryCode {
		output := []string{
			countryCodeISO2,
		}
		header[0] = countryCodeISO3
		output = append(output, header...)

		fmt.Println(strings.Join(output, ","))
	} else {
		return fmt.Errorf("header %#q not recognized", header)
	}

	return nil
}
