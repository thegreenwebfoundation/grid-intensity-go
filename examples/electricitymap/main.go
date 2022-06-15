package main

import (
	"context"
	"log"
	"os"

	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go/api"
	"github.com/thegreenwebfoundation/grid-intensity-go/electricitymap"
)

func main() {
	token := os.Getenv("ELECTRICITY_MAP_API_TOKEN")
	if token == "" {
		log.Fatalln("please set the env variable `ELECTRICITY_MAP_API_TOKEN`")
	}

	c, err := electricitymap.New(token)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIntensity(context.Background(), "IN-KA")
	log.Println(a, err)

	carbonIntensityMap, err := gridintensity.GetCarbonIntensityMap(context.Background(), c, "IN-KA", "IN-AP")
	log.Println(carbonIntensityMap.GetAll(), err)
}
