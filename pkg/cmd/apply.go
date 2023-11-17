package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
)

var ApplyCmd = &cobra.Command{
	Use:               common.ApplyCommandName,
	DisableAutoGenTag: true,
	Short:             "Apply a specification",
	Long:              "Applies an existing configuration, such as a data policy",
}

func init() {
	ApplyCmd.AddCommand(datapolicy.ApplyCmd())
}
