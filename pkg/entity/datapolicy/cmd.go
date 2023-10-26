package datapolicy

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/processingplatform"
)

const bareFlag = "bare"
const bareFlagShort = "b"

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (table-id|policy-id)",
		Short:             "Get a data policy",
		Long:              getHelp,
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy or table id
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	common.SetOutputFormats(flags, common.OutputFormatYaml, common.OutputFormatJson, common.OutputFormatJsonRaw)
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(flags)
	catalog.AddSchemaFlag(flags)
	flags.BoolP(bareFlag, bareFlagShort, false, "when true ask platform or catalog, otherwise ask Pace")
	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policies",
		Short:             "List Datapolicies",
		Example:           "",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	flags := cmd.Flags()
	common.SetOutputFormats(flags, common.OutputFormatYaml, common.OutputFormatJson, common.OutputFormatJsonRaw)
	return cmd
}
