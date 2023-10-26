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
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(flags)
	catalog.AddSchemaFlag(flags)
	return cmd
}
