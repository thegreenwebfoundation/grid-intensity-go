package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func main() {
	c := provider.CarbonIntensityUKConfig{}
	p, err := provider.NewCarbonIntensityUK(c)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}

	res, err := p.GetCarbonIntensity(context.Background(), "UK")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	bytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	fmt.Println(string(bytes))
}
