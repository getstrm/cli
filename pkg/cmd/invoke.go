package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/plugin"
)

var InvokeCmd = &cobra.Command{
	Use:               common.InvokeCommandName,
	DisableAutoGenTag: true,
	Short:             "Invoke a functionality",
	Long:              "Invoke a functionality, such as a plugin",
}

func init() {
	InvokeCmd.AddCommand(plugin.InvokeCmd())
}
