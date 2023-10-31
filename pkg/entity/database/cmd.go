package database

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "databases",
		Short:             "List Databases",
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, availablePrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	common.ConfigureExtraPrinters(cmd, flags, availablePrinters())
	catalog.AddCatalogFlag(cmd, flags)
	return cmd
}
