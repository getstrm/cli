package datapolicy

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/processingplatform"
)

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "datapolicy (table-id|policy-id)",
		Short:             "Get a bare DataPolicy",
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", "")
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", "")
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", "")
	flags.StringP(common.SchemaFlag, common.SchemaFlagShort, "", "")
	flags.BoolP("bare", "b", false, "")
	err := cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, completion)
	common.CliExit(err)
	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "datapolicies",
		Short:             "List Datapolicies",
		Example:           "",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	return cmd
}

func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := processingplatform.PlatformIdsCompletion(cmd, args, complete)
	return s, c
}
