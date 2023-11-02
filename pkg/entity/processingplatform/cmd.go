package processingplatform

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "processing-platforms",
		Short:             "List Processing Platforms",
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
	common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
	return cmd
}
