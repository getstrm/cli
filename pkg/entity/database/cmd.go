package database

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/completion"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "databases",
		Short:             "List Databases",
		Long:              listDatabasesLongDocs,
		Example:           listDatabasesExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
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
	_ = common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	completion.AddProcessingPlatformFlag(cmd, flags)
	completion.AddCatalogFlag(cmd, flags)
	cmd.MarkFlagsOneRequired(common.ProcessingPlatformFlag, common.CatalogFlag)
	return cmd
}
