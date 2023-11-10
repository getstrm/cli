package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
	"pace/pace/pkg/entity/globaltransform"
)

var GetCmd = &cobra.Command{
	Use:               common.GetCommandName,
	DisableAutoGenTag: true,
	Short:             "Get a single entity",
}

func init() {
	GetCmd.AddCommand(datapolicy.GetCmd())
	GetCmd.AddCommand(globaltransform.GetCmd())
}
