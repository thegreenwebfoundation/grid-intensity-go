package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xuanwo/go-locale"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thegreenwebfoundation/grid-intensity-go/ember"
)

const (
	configDir      = ".config/grid-intensity"
	configFileName = "config.yaml"
	provider       = "provider"
	region         = "region"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grid-intensity",
	Short: "Get carbon intensity data for electricity grids",
	Long: `A tool for getting the carbon intensity data for electricity grids.

This can be used to make your sofware carbon aware so it runs at times when the
grid is greener or at locations where carbon intensity is lower.

	grid-intensity --region ARG
	grid-intensity -r BOL`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runWithError()
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
	rootCmd.Flags().StringP(provider, "p", ember.Provider, "Provider of carbon intensity data")
	rootCmd.Flags().StringP(region, "r", "", "Region code for provider")

	viper.BindPFlag(provider, rootCmd.Flags().Lookup(provider))
	viper.BindPFlag(region, rootCmd.Flags().Lookup(region))
}

func getEmberGridIntensityForCountry(countryCode string) error {
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

// getCountryCode prompts the user to enter a country code. We try to detect
// a country code from the user's locale but the user can enter another value.
func getCountryCode() (string, error) {
	tag, err := locale.Detect()
	if err != nil {
		return "", err
	}

	region, _ := tag.Region()
	country := region.ISO3()

	fmt.Printf("Provider %s needs an ISO country code as a region parameter.\n", ember.Provider)
	if country != "" {
		fmt.Printf("%s detected from your locale.\n", country)
	}

	var reader = bufio.NewReader(os.Stdin)
	country, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(country), nil
}

func runWithError() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configFile := filepath.Join(homeDir, configDir, configFileName)
	viper.SetConfigFile(configFile)

	err = viper.ReadInConfig()
	if _, ok := err.(*fs.PathError); ok {
		// Create config dir if it doesn't exist.
		err = os.Mkdir(filepath.Join(homeDir, configDir), os.ModePerm)
		if err != nil {
			return err
		}

		// Create config file if it doesn't exist.
		_, err = os.Create(configFile)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	providerName := viper.GetString(provider)
	regionCode := viper.GetString(region)

	switch providerName {
	case ember.Provider:
		if regionCode == "" {
			regionCode, err = getCountryCode()
			if err != nil {
				return err
			}

			viper.Set(region, regionCode)
		}

		err = getEmberGridIntensityForCountry(regionCode)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("provider %q not recognized", providerName)
	}

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
