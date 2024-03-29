package cmd

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/database"
	"pace/pace/pkg/entity/datapolicy"
	"pace/pace/pkg/entity/globaltransform"
	"pace/pace/pkg/entity/group"
	"pace/pace/pkg/entity/plugin"
	"pace/pace/pkg/entity/processingplatform"
	"pace/pace/pkg/entity/resources"
	"pace/pace/pkg/entity/schema"
	"pace/pace/pkg/entity/table"
)

var ListCmd = &cobra.Command{
	Use:               common.ListCommandName,
	DisableAutoGenTag: true,
	Short:             "List entities",
	Long:              "return 0 or more entities of a certain type",
}

func init() {
	ListCmd.AddCommand(
		processingplatform.ListCmd(),
		catalog.ListCmd(),
		table.ListCmd(),
		group.ListCmd(),
		schema.ListCmd(),
		database.ListCmd(),
		datapolicy.ListCmd(),
		datapolicy.ScanLineage(),
		globaltransform.ListCmd(),
		plugin.ListCmd(),
		resources.ListCmd(),
	)

	flags := ListCmd.PersistentFlags()
	flags.Uint32P(common.PageSizeFlag, "P", 10, "the maximum number of records per page")
	flags.Uint32P(common.PageSkipFlag, "S", 0, "the number of records that need to be skipped")
	flags.StringP(common.PageTokenFlag, "T", "", "next page token. Used by BigQuery")
}
