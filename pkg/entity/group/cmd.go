package group

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/processingplatform"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "groups",
		Short:             "List Groups",
		Long:              listLongDocs,
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
	processingplatform.AddProcessingPlatformFlag(cmd, flags)

	cmd.MarkFlagRequired(common.ProcessingPlatformFlag)
	return cmd
}
