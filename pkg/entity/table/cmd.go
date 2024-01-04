package table

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/processingplatform"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "tables",
		Short:             "List Tables",
		Long:              listTablesLongDocs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	_ = common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	_ = processingplatform.AddProcessingPlatformFlag(cmd, flags)
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(cmd, flags)
	catalog.AddSchemaFlag(cmd, flags)

	cmd.MarkFlagsRequiredTogether(common.CatalogFlag, common.DatabaseFlag, common.SchemaFlag)
	cmd.MarkFlagsOneRequired(common.ProcessingPlatformFlag, common.CatalogFlag)
	return cmd
}
