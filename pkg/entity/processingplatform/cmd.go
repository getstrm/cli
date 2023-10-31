package processingplatform

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "processing-platforms",
		Short:             "List Processing Platforms",
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
	return cmd
}
