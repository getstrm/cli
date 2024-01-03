package schema

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "schemas",
		Short:             "List Schemas",
		Long:              listSchemasLongDocs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(cmd, flags)

	cmd.MarkFlagRequired(common.CatalogFlag)
	cmd.MarkFlagRequired(common.DatabaseFlag)

	return cmd
}
