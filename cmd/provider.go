package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(providerCmd)
}

var (
	providerCmd = &cobra.Command{
		Use:   "provider",
		Short: "Details about providers of carbon intensity data",
		Long:  "Details about providers of carbon intensity data",
	}
)
