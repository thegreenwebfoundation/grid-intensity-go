package main

import (
	"context"
	"log"

	"github.com/thegreenwebfoundation/grid-intensity-go/watttime"
)

func main() {
	// Register via the API
	// https://www.watttime.org/api-documentation/#register-new-user
	wattTimeUser := "your-user"
	wattTimePassword := "your-password"

	c, err := watttime.New(wattTimeUser, wattTimePassword)
	if err != nil {
		log.Fatalln("could not make provider", err)
	}
	a, err := c.GetCarbonIntensity(context.Background(), "CAISO_NORTH")
	log.Println(a, err)
}
