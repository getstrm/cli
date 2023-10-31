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
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	return cmd
}
