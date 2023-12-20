package datapolicy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	"pace/pace/pkg/completion"
	. "pace/pace/pkg/util"
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
		ValidArgsFunction: idsCompletion,
	}

	flags := cmd.Flags()
	completion.AddProcessingPlatformFlag(cmd, flags)
	cmd.MarkFlagRequired(common.ProcessingPlatformFlag)

	return cmd
}

func EvaluateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (policy-id)",
		Short:             "Evaluate an existing data policy by applying it to sample data provided in a csv file",
		Long:              evaluateLongDocs,
		Example:           evaluateExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, evaluatePrinters())
		},
		Run: func(cmd *cobra.Command, args []string) {
			evaluate(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy id
		ValidArgsFunction: idsCompletion,
	}

	flags := cmd.Flags()
	completion.AddProcessingPlatformFlag(cmd, flags)
	cmd.MarkFlagRequired(common.ProcessingPlatformFlag)
	addSampleDataFlag(cmd, flags)
	cmd.MarkFlagRequired(common.SampleDataFlag)

	common.ConfigureExtraPrinters(cmd, cmd.Flags(), evaluatePrinters())

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
	completion.AddProcessingPlatformFlag(cmd, flags)
	completion.AddCatalogFlag(cmd, flags)
	completion.AddDatabaseFlag(cmd, flags)
	completion.AddSchemaFlag(cmd, flags)
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

func addSampleDataFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.String(common.SampleDataFlag, "", common.SampleDataUsage)
	CliExit(
		cmd.RegisterFlagCompletionFunc(common.SampleDataFlag,
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return []string{"csv"}, cobra.ShellCompDirectiveFilterFileExt
			}),
	)
}
