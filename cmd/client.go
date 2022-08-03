package cmd

import (
	"fmt"
	"os"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func getClient(providerName string, cacheFile string) (provider.Interface, error) {
	var client provider.Interface
	var err error

	switch providerName {
	case provider.CarbonIntensityOrgUK:
		c := provider.CarbonIntensityUKConfig{}
		client, err = provider.NewCarbonIntensityUK(c)
		if err != nil {
			return nil, fmt.Errorf("could not make carbon intensity uk provider, %w", err)
		}
	case provider.ElectricityMap:
		token := os.Getenv(electricityMapAPITokenEnvVar)
		if token == "" {
			return nil, fmt.Errorf("%q env var must be set", electricityMapAPITokenEnvVar)
		}

		c := provider.ElectricityMapConfig{
			Token: token,
		}
		client, err = provider.NewElectricityMap(c)
		if err != nil {
			return nil, fmt.Errorf("could not make electricity map provider, %w", err)
		}
	case provider.Ember:
		client, err = provider.NewEmber()
		if err != nil {
			return nil, fmt.Errorf("could not make ember provider, %w", err)
		}
	case provider.WattTime:
		user := os.Getenv(wattTimeUserEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimeUserEnvVar)
		}

		password := os.Getenv(wattTimePasswordEnvVar)
		if user == "" {
			return nil, fmt.Errorf("%q env var must be set", wattTimePasswordEnvVar)
		}

		c := provider.WattTimeConfig{
			APIUser:     user,
			APIPassword: password,
			CacheFile:   cacheFile,
		}
		client, err = provider.NewWattTime(c)
		if err != nil {
			return nil, fmt.Errorf("could not make watt time provider, %w", err)
		}
	default:
		return nil, fmt.Errorf("provider %q not supported", providerName)
	}

	return client, nil
}
