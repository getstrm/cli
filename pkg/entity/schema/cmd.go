package schema

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/processingplatform"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "schemas",
		Short:             "List Schemas",
		Long:              listSchemasLongDocs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, listPrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	catalog.AddCatalogFlag(cmd, flags)
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	catalog.AddDatabaseFlag(cmd, flags)

	cmd.MarkFlagRequired(common.DatabaseFlag)
	cmd.MarkFlagsOneRequired(common.ProcessingPlatformFlag, common.CatalogFlag)

	return cmd
}
