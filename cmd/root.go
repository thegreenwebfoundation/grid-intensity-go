package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Xuanwo/go-locale"
	"github.com/spf13/cobra"

	"github.com/thegreenwebfoundation/grid-intensity-go/ember"
)

const (
	countryCode = "country-code"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grid-intensity",
	Short: "Get carbon intensity data for electricity grids",
	Long: `A tool for getting the carbon intensity data for electricity grids.

This can be used to make your sofware carbon aware so it runs at times when the
grid is greener or at locations where carbon intensity is lower.

	grid-intensity --country-code ARG
	grid-intensity -c BOL`,

	Run: func(cmd *cobra.Command, args []string) {
		country, err := cmd.Flags().GetString(countryCode)
		if err != nil {
			log.Fatal(err)
		}

		err = getGridIntensityForCountry(country)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(countryCode, "c", "", "Country code in ISO 3 character format")
}

func getGridIntensityForCountry(countryCode string) error {
	if countryCode == "" {
		// No country code provided so try to detect it from the user's locale.
		tag, err := locale.Detect()
		if err != nil {
			return err
		}

		region, _ := tag.Region()
		countryCode = region.ISO3()
	}

	result, err := ember.GetGridIntensityForCountry(countryCode)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}
