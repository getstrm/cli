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
	flags := cmd.Flags()
	_ = common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	_ = processingplatform.AddProcessingPlatformFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.ProcessingPlatformFlag)
	return cmd
}
