package main

import (
	"context"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/unfccc"
)

func main() {
	c, err := unfccc.New("../../unfccc/sampleRegionData.json")
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIntensity(context.Background(), "de")
	log.Println(a, err)
}
