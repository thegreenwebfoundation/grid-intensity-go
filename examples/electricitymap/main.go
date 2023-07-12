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
	token := os.Getenv("ELECTRICITY_MAP_API_TOKEN")
	if token == "" {
		log.Fatalln("please set the env variable `ELECTRICITY_MAP_API_TOKEN`")
	}

	c := provider.ElectricityMapConfig{
		Token: token,
	}
	url := os.Getenv("ELECTRICITY_MAP_API_URL")
	if url != "" {
		c.APIURL = url
	}
	e, err := provider.NewElectricityMap(c)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}

	res, err := e.GetCarbonIntensity(context.Background(), "AU-SA")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	bytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	fmt.Println(string(bytes))
}
