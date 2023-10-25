package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/processingplatform"
)

var ListCmd = &cobra.Command{
	Use:               common.ListCommandName,
	DisableAutoGenTag: true,
	Short:             "List entities",
}

func init() {
	ListCmd.AddCommand(processingplatform.ListCmd())

	ListCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
