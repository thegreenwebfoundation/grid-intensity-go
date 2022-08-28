package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xuanwo/go-locale"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

const (
	cacheDir              = ".cache/grid-intensity"
	configDir             = ".config/grid-intensity"
	configFileName        = "config.yaml"
	locationKey           = "location"
	providerKey           = "provider"
	wattTimeCacheFileName = "watttime.org.json"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grid-intensity",
	Short: "Get carbon intensity data for electricity grids",
	Long: `A tool for getting the carbon intensity data for electricity grids.

This can be used to make your software carbon aware so it runs at times when the
grid is greener or at locations where carbon intensity is lower.

	grid-intensity --provider Ember --location ARG
	grid-intensity -p Ember -l BOL`,

	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(locationKey, cmd.Flags().Lookup(locationKey))
		viper.BindPFlag(providerKey, cmd.Flags().Lookup(providerKey))
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runRoot()
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
	rootCmd.Flags().StringP(locationKey, "l", "", "Location code for provider")
	rootCmd.Flags().StringP(providerKey, "p", provider.Ember, "Provider of carbon intensity data")

	// Also support environment variables.
	viper.SetEnvPrefix("grid_intensity")
	viper.BindEnv(locationKey)
	viper.BindEnv(providerKey)
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

	fmt.Printf("Provider %s needs an ISO country code as a location parameter.\n", provider.Ember)
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

func runRoot() error {
	ctx := context.Background()

	providerName, err := readConfig(providerKey)
	if err != nil {
		return fmt.Errorf("could not read config for %#q, %w", providerKey, err)
	}
	locationCode, err := readConfig(locationKey)
	if err != nil {
		return fmt.Errorf("could not read config for %#q, %w", locationKey, err)
	}

	var cacheFile string

	switch providerName {
	case provider.CarbonIntensityOrgUK:
		if locationCode == "" {
			locationCode = "UK"
		}
		if locationCode != "UK" {
			return fmt.Errorf("only location UK is supported")
		}
		viper.Set(locationKey, locationCode)
	case provider.Ember:
		if locationCode == "" {
			locationCode, err = getCountryCode()
			if err != nil {
				return err
			}
			viper.Set(locationKey, locationCode)
		}
	case provider.WattTime:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		// Use file cache to avoid API rate limiting.
		cacheFile = filepath.Join(homeDir, cacheDir, wattTimeCacheFileName)
	}

	client, err := getClient(providerName, cacheFile)
	if err != nil {
		return fmt.Errorf("could not get client, %w", err)
	}

	res, err := client.GetCarbonIntensity(ctx, locationCode)
	if err != nil {
		return fmt.Errorf("could not get carbon intensity, %w", err)
	}

	bytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return fmt.Errorf("could not marshal json, %w", err)
	}

	fmt.Println(string(bytes))

	err = writeConfig()
	if err != nil {
		return fmt.Errorf("could not write config, %w", err)
	}

	return nil
}
