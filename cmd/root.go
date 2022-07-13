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

	"github.com/thegreenwebfoundation/grid-intensity-go/ember"
	"github.com/thegreenwebfoundation/grid-intensity-go/watttime"
)

const (
	cacheDir       = ".cache/grid-intensity"
	cacheFileName  = "watttime.org.json"
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

This can be used to make your software carbon aware so it runs at times when the
grid is greener or at locations where carbon intensity is lower.

	grid-intensity --region ARG
	grid-intensity -r BOL`,

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
	// Use persistent flags so they are available to all subcommands.
	rootCmd.PersistentFlags().StringP(provider, "p", ember.ProviderName, "Provider of carbon intensity data")
	rootCmd.PersistentFlags().StringP(region, "r", "", "Region code for provider")

	viper.BindPFlag(provider, rootCmd.PersistentFlags().Lookup(provider))
	viper.BindPFlag(region, rootCmd.PersistentFlags().Lookup(region))

	// Also support environment variables.
	viper.SetEnvPrefix("grid_intensity")
	viper.BindEnv(provider)
	viper.BindEnv(region)
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

func getWattTimeGridIntensity(ctx context.Context, region string) error {
	user := os.Getenv(wattTimeUserEnvVar)
	if user == "" {
		return fmt.Errorf("%q env var must be set", wattTimeUserEnvVar)
	}

	password := os.Getenv(wattTimePasswordEnvVar)
	if user == "" {
		return fmt.Errorf("%q env var must be set", wattTimePasswordEnvVar)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	cacheFile := filepath.Join(homeDir, cacheDir, cacheFileName)
	c, err := watttime.New(user, password, watttime.WithCacheFile(cacheFile))
	if err != nil {
		return fmt.Errorf("could not make provider %v", err)
	}

	result, err := c.GetCarbonIntensityData(ctx, region)
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

	fmt.Printf("Provider %s needs an ISO country code as a region parameter.\n", ember.ProviderName)
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

	providerName, regionCode, err := readConfig()
	if err != nil {
		return err
	}

	switch providerName {
	case ember.ProviderName:
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
	case watttime.ProviderName:
		err = getWattTimeGridIntensity(ctx, regionCode)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("provider %q not recognized", providerName)
	}

	err = writeConfig()
	if err != nil {
		return err
	}

	return nil
}
