package table

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/completion"
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
	completion.AddProcessingPlatformFlag(cmd, flags)
	completion.AddCatalogFlag(cmd, flags)
	completion.AddDatabaseFlag(cmd, flags)
	completion.AddSchemaFlag(cmd, flags)

	_ = cmd.MarkFlagRequired(common.DatabaseFlag)
	_ = cmd.MarkFlagRequired(common.SchemaFlag)
	cmd.MarkFlagsOneRequired(common.ProcessingPlatformFlag, common.CatalogFlag)
	return cmd
}
