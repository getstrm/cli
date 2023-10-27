package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
)

var GetCmd = &cobra.Command{
	Use:               common.GetCommandName,
	DisableAutoGenTag: true,
	Short:             "Get entities",
}

func init() {
	GetCmd.AddCommand(datapolicy.GetCmd())
}
