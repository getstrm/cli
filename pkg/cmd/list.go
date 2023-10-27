package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/database"
	"pace/pace/pkg/entity/datapolicy"
	"pace/pace/pkg/entity/group"
	"pace/pace/pkg/entity/processingplatform"
	"pace/pace/pkg/entity/schema"
	"pace/pace/pkg/entity/table"
)

var ListCmd = &cobra.Command{
	Use:               common.ListCommandName,
	DisableAutoGenTag: true,
	Short:             "List entities",
}

func init() {
	ListCmd.AddCommand(processingplatform.ListCmd())
	ListCmd.AddCommand(catalog.ListCmd())
	ListCmd.AddCommand(table.ListCmd())
	ListCmd.AddCommand(group.ListCmd())
	ListCmd.AddCommand(schema.ListCmd())
	ListCmd.AddCommand(database.ListCmd())
	ListCmd.AddCommand(datapolicy.ListCmd())
}
