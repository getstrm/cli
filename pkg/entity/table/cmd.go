package table

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
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
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", "snowflake-demo")
	err := cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, completion)
	common.CliExit(err)
	return cmd
}
func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := processingplatform.PlatformIdsCompletion(cmd, args, complete)
	return s, c
}
