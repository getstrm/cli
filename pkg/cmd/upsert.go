package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/datapolicy"
	"pace/pace/pkg/entity/globaltransform"
)

var UpsertCmd = &cobra.Command{
	Use:               common.UpsertCommandName,
	DisableAutoGenTag: true,
	Short:             "Upsert entities",
	Long:              "Insert or Update an entity",
}

func init() {
	UpsertCmd.AddCommand(datapolicy.UpsertCmd())

	UpsertCmd.AddCommand(globaltransform.UpsertCmd())
}
