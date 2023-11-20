package datapolicy

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/processingplatform"
)

func UpsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (yaml or json file)",
		Short:             "Upsert a data policy",
		Long:              upsertLongDocs,
		Example:           upsertExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			upsert(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1), // the policy file (yaml or json),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return common.DefaultFileTypesCompletion, cobra.ShellCompDirectiveFilterFileExt
		},
	}

	flags := cmd.Flags()
	flags.BoolP(common.ApplyFlag, common.ApplyFlagShort, false, common.ApplyFlagUsage)
	return cmd
}

func ApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (policy-id)",
		Short:             "Apply an existing data policy",
		Long:              applyLongDocs,
		Example:           applyExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			apply(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy id
		ValidArgsFunction: IdsCompletion,
	}

	flags := cmd.Flags()
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	cmd.MarkFlagRequired(common.ProcessingPlatformFlag)

	return cmd
}

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (table-id|policy-id)",
		Short:             "Get a data policy",
		Long:              getLongDoc,
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy or table id
		ValidArgsFunction: TableOrDataPolicyIdsCompletion,
	}
	flags := cmd.Flags()
	processingplatform.AddProcessingPlatformFlag(cmd, flags)
	catalog.AddCatalogFlag(cmd, flags)
	catalog.AddDatabaseFlag(cmd, flags)
	catalog.AddSchemaFlag(cmd, flags)
	flags.BoolP(common.BlueprintFlag, common.BlueprintFlagShort, false, common.BlueprintFlagUsage)
	cmd.MarkFlagsMutuallyExclusive(common.CatalogFlag, common.ProcessingPlatformFlag)
	cmd.MarkFlagsOneRequired(common.CatalogFlag, common.ProcessingPlatformFlag)
	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policies",
		Short:             "List Data Policies",
		Example:           listExample,
		Long:              listLongDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, listPrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
	return cmd
}
