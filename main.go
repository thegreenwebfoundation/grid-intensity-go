package main

import (
	"fmt"
	"log"

	"github.com/Xuanwo/go-locale"
)

func main() {
	tag, err := locale.Detect()
	if err != nil {
		log.Fatal(err)
	}

	region, _ := tag.Region()
	fmt.Printf("Looks like your region is %s\n", region.ISO3())
}
