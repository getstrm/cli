package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
)

var TranspileCmd = &cobra.Command{
	Use:               common.TranspileCommandName,
	DisableAutoGenTag: true,
	Short:             "Transpile a specification",
	Long:              "Transpiles an existing specification, such as a data policy",
}

func init() {
	TranspileCmd.AddCommand(datapolicy.TranspileCmd())
}
