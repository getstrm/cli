package catalog

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "catalogs",
		Short:             "List Catalogs",
		Long:              listCatalogsDocs,
		Example:           listCatalogsExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	_ = common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	return cmd
}
