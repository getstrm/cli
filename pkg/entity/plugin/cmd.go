package plugin

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
)

func InvokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "plugin (plugin-id) --payload (payload-file)",
		Short:             "Invoke a plugin with the provided payload (JSON or YAML)",
		Long:              "tbd", // Todo - add long docs
		Example:           "tbd", // Todo - add example
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			invokePlugin(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy id
		ValidArgsFunction: IdsCompletion,
	}

	flags := cmd.Flags()
	addPayloadFlag(cmd, flags)
	cmd.MarkFlagRequired(common.PluginPayloadFlag)

	return cmd
}

func addPayloadFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.String(common.PluginPayloadFlag, "", common.PluginPayloadFlagUsage)
	CliExit(
		cmd.RegisterFlagCompletionFunc(common.PluginPayloadFlag,
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return []string{"yaml", "json"}, cobra.ShellCompDirectiveFilterFileExt
			}),
	)
}
