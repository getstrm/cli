package plugin

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
)

const pluginCommand = "plugin"

func InvokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               fmt.Sprintf("%s (plugin-id) (action) --payload (payload-file)", pluginCommand),
		Short:             "Invoke an action for a plugin with the provided payload (JSON or YAML)",
		Long:              invokeLongDocs,
		Example:           invokeExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			invokePlugin(cmd, args)
		},
		Args:              cobra.RangeArgs(1, 2), // the plugin id and optional action
		ValidArgsFunction: IdsCompletion,
	}

	flags := cmd.Flags()
	addPayloadFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.PluginPayloadFlag)

	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               fmt.Sprintf("%ss", pluginCommand),
		Short:             "List plugins",
		Long:              listLongDocs,
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, listPrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			listPlugins()
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
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
