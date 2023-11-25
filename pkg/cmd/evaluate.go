package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
)

var EvaluateCmd = &cobra.Command{
	Use:               common.EvaluateCommandName,
	DisableAutoGenTag: true,
	Short:             "Evaluate a specification",
	Long:              "Evaluates an existing specification, such as a data policy",
}

func init() {
	EvaluateCmd.AddCommand(datapolicy.EvaluateCmd())
}
