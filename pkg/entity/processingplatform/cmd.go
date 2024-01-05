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
	_ = common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
	return cmd
}
