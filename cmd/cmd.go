package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	electricityMapAPITokenEnvVar = "ELECTRICITY_MAP_API_TOKEN"
	wattTimeUserEnvVar           = "WATT_TIME_USER"
	wattTimePasswordEnvVar       = "WATT_TIME_PASSWORD"
)

func getConfig() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", nil
	}

	configFile := filepath.Join(homeDir, configDir, configFileName)
	viper.SetConfigFile(configFile)

	err = viper.ReadInConfig()
	if _, ok := err.(*fs.PathError); ok {
		// Create config dir if it doesn't exist.
		os.Mkdir(filepath.Join(homeDir, configDir), os.ModePerm)

		// Create config file if it doesn't exist.
		_, err = os.Create(configFile)
		if err != nil {
			return "", "", err
		}
	} else if err != nil {
		return "", "", err
	}

	providerName := viper.GetString(provider)
	regionCode := viper.GetString(region)

	return providerName, regionCode, nil
}
