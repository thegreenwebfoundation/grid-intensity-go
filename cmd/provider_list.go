package cmd

import (
	"log"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/thegreenwebfoundation/grid-intensity-go/pkg/provider"
)

func init() {
	providerCmd.AddCommand(providerListCmd)
}

var (
	providerListCmd = &cobra.Command{
		Use:   "list",
		Short: "List supported providers of carbon intensity data",
		Long: `List all supported providers of carbon intensity data
for electricity grids.

	grid-intensity provider list`,
		Run: func(cmd *cobra.Command, args []string) {
			err := runProviderList()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func runProviderList() error {
	providers := provider.GetProviderDetails()

	tbl := table.New("NAME", "URL")

	for _, p := range providers {
		tbl.AddRow(p.Name, p.URL)
	}

	tbl.Print()

	return nil
}
