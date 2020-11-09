package main

import (
	"context"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/carbonintensity"
)

func main() {
	c, err := carbonintensity.New()
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIntensity(context.Background(), "UK")
	log.Println(a, err)
}
