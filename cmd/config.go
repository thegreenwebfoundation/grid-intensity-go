package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	electricityMapAPITokenEnvVar = "ELECTRICITY_MAP_API_TOKEN"
	wattTimeUserEnvVar           = "WATT_TIME_USER"
	wattTimePasswordEnvVar       = "WATT_TIME_PASSWORD"
)

func getConfigFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	return filepath.Join(homeDir, configDir, configFileName), nil
}

func configFileExists(configFile string) (bool, error) {
	_, err := os.Stat(configFile)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func readConfig() (string, string, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return "", "", err
	}

	viper.SetConfigFile(configFile)

	err = viper.ReadInConfig()
	if errors.Is(err, os.ErrNotExist) {
		// Config file may not be available e.g. when running as a container.
	} else if err != nil {
		return "", "", err
	}

	providerName := viper.GetString(providerKey)
	regionCode := viper.GetString(regionKey)

	return providerName, regionCode, nil
}

func writeConfig() error {
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}

	fileExists, err := configFileExists(configFile)
	if err != nil {
		return err
	}

	if !fileExists {
		// Create config dir if it doesn't exist.
		err = os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return nil
		}

		// Create config file if it doesn't exist.
		_, err = os.Create(configFile)
		if err != nil {
			// If we can't create file don't try to write config.
			return nil
		}
	}

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
