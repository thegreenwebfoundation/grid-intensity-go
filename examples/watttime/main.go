package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func main() {
	// Register via the API
	// https://www.watttime.org/api-documentation/#register-new-user
	c := provider.WattTimeConfig{
		APIUser:     "your-user",
		APIPassword: "your-password",
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
