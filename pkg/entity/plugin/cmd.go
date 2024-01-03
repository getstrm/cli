package plugin

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
)

const pluginCommand = "plugin"

func InvokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               fmt.Sprintf("%s (plugin-id) (action) --payload (payload-file)", pluginCommand),
		Short:             "Invoke an action for a plugin with the provided payload (JSON or YAML)",
		Long:              invokeLongDocs,
		Example:           invokeExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, common.StandardPrinters)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return invokePlugin(cmd, args)
		},
		Args:              cobra.RangeArgs(1, 2), // the plugin id and optional action
		ValidArgsFunction: IdsCompletion,
	}

	flags := cmd.Flags()
	_ = addPayloadFlag(cmd, flags)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPlugins()
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
	return cmd
}

func addPayloadFlag(cmd *cobra.Command, flags *pflag.FlagSet) error {
	flags.String(common.PluginPayloadFlag, "", common.PluginPayloadFlagUsage)
	return cmd.RegisterFlagCompletionFunc(common.PluginPayloadFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"yaml", "json"}, cobra.ShellCompDirectiveFilterFileExt
		})
}
