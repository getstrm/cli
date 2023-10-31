package catalog

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "catalogs",
		Short:             "List Catalogs",
		Long:              docs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, availablePrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
}
