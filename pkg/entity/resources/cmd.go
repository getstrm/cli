package resources

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "resources (resource-path)",
		Short:             "list resources",
		Long:              listDocs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd, args)
		},
		// fqn of table
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: resourcesCompletion,
	}
	flags := cmd.Flags()
	_ = common.ConfigureExtraPrinters(cmd, flags, listPrinters())
	return cmd
}
