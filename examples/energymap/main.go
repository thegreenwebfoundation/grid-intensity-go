package main

import (
	"context"
	"log"
	"os"

	"github.com/thegreenwebfoundation/grid-intensity-go/energymap"
)

func main() {
	token := os.Getenv("ENERGY_MAP_API_TOKEN")
	if token == "" {
		log.Fatalln("please set the env variable `ENERGY_MAP_API_TOKEN`")
	}

	c, err := energymap.New(token)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIndex(context.Background(), "IN-KA")
	log.Println(a, err)

	lowest, err := c.MinIntensity(context.Background(), "IN-KA", "IN-AP")
	log.Println(lowest, err)
}
