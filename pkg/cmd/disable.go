package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/welcome"
)

var DisableCmd = &cobra.Command{
	Use:               common.DisableCommandName,
	DisableAutoGenTag: true,
	Short:             "Disable operation",
	Long:              "Disables an existing configuration, such as welcoming message",
}

func init() {
	DisableCmd.AddCommand(welcome.DisableCmd())
}
