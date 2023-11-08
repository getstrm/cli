package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/globaltransform"
)

var DeleteCmd = &cobra.Command{
	Use:               common.DeleteCommandName,
	DisableAutoGenTag: true,
	Short:             "Delete entities",
}

func init() {
	DeleteCmd.AddCommand(globaltransform.DeleteCmd())
}
