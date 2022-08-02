package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func main() {
	p, err := provider.NewEmber()
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	result, err := p.GetCarbonIntensity(context.Background(), "ESP")
	if err != nil {
		log.Fatalln("could not get carbon intesity", err)
	}

	bytes, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	fmt.Println(string(bytes))
}
