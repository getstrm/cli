package schema

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/completion"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "schemas",
		Short:             "List Schemas",
		Long:              listSchemasLongDocs,
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
	completion.AddCatalogFlag(cmd, flags)
	completion.AddProcessingPlatformFlag(cmd, flags)
	completion.AddDatabaseFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.DatabaseFlag)
	cmd.MarkFlagsOneRequired(common.ProcessingPlatformFlag, common.CatalogFlag)
	return cmd
}
