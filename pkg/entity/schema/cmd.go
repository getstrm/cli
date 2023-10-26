package schema

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "schemas",
		Short:             "List Schemas",
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
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", "")
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", "")
	err := cmd.RegisterFlagCompletionFunc(common.CatalogFlag, completion)
	common.CliExit(err)
	return cmd
}

func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := catalog.CatalogIdsCompletion(cmd, args, complete)
	return s, c
}
