package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func main() {
	// Register via the API
	// https://www.watttime.org/api-documentation/#register-new-user
	apiUser := os.Getenv("WATT_TIME_API_USER")
	if apiUser == "" {
		log.Fatalln("please set the env variable `WATT_TIME_API_USER`")
	}
	apiPassword := os.Getenv("WATT_TIME_API_PASSWORD")
	if apiPassword != "" {
		log.Fatalln("please set the env variable `WATT_TIME_API_PASSWORD`")
	}

	c := provider.WattTimeConfig{
		APIUser:     apiUser,
		APIPassword: apiPassword,
	}
	w, err := provider.NewWattTime(c)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}

	res, err := w.GetCarbonIntensity(context.Background(), "CAISO_NORTH")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	bytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	fmt.Println(string(bytes))
}
