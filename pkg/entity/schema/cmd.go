package schema

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/util"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "schemas",
		Short:             "List Schemas",
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = util.ConfigurePrinter(cmd, availablePrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(flags)
	return cmd
}
