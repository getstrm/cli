package datapolicy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	"pace/pace/pkg/completion"
)

func UpsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (yaml or json file)",
		Short:             "Upsert a data policy",
		Long:              upsertLongDocs,
		Example:           upsertExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, common.StandardPrinters)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return upsert(cmd, &args[0])
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
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, common.StandardPrinters)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy id
		ValidArgsFunction: idsCompletion,
	}

	flags := cmd.Flags()
	completion.AddProcessingPlatformFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.ProcessingPlatformFlag)
	return cmd
}

func EvaluateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (policy-id)",
		Short:             "Evaluate an existing data policy by applying it to sample data provided in a csv file",
		Long:              evaluateLongDocs,
		Example:           evaluateExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, evaluatePrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return evaluate(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the policy id
		ValidArgsFunction: idsCompletion,
	}

	flags := cmd.Flags()
	completion.AddProcessingPlatformFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.ProcessingPlatformFlag)
	_ = addSampleDataFlag(cmd, flags)
	_ = cmd.MarkFlagRequired(common.SampleDataFlag)
	_ = common.ConfigureExtraPrinters(cmd, cmd.Flags(), evaluatePrinters())
	return cmd
}

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policy (table-name|policy-id)",
		Short:             "Get a data policy",
		Long:              getLongDoc,
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, common.StandardPrinters)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd, &args[0])
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
	flags.Bool(FqnFlag, false, "use argument as fqn")
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
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			printer, err = common.ConfigurePrinter(cmd, listPrinters())
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
	_ = common.ConfigureExtraPrinters(cmd, cmd.Flags(), listPrinters())
	return cmd
}

func ScanLineage() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "lineage",
		Short:             "List lineage for all stored data-policies",
		Example:           "",
		Long:              "",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scanLineage(cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
		Args:              cobra.NoArgs,
	}
	_ = common.ConfigureExtraPrinters(cmd, cmd.Flags(), common.StandardPrinters)
	return cmd
}

func addSampleDataFlag(cmd *cobra.Command, flags *pflag.FlagSet) error {
	flags.String(common.SampleDataFlag, "", common.SampleDataUsage)
	return cmd.RegisterFlagCompletionFunc(common.SampleDataFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"csv"}, cobra.ShellCompDirectiveFilterFileExt
		})
}
