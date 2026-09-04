package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func main() {
	c := provider.GridCarbonConfig{}
	p, err := provider.NewGridCarbon(c)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}

	// Any gridcarbon zone code: FR, DE-LU, US-NYISO, GB, ... (https://api.gridcarbon.dev/v1/zones)
	res, err := p.GetCarbonIntensity(context.Background(), "FR")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	bytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatalln("could not get carbon intensity", err)
	}

	fmt.Println(string(bytes))
}
