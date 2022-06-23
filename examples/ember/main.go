package main

import (
	"context"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/ember"
)

func main() {
	c, err := ember.New()
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIntensity(context.Background(), "ESP")
	log.Println(a, err)
}
