package globaltransform

import (
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

func UpsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "global-transform (yaml or json file)",
		Short:             "Upsert a global transform",
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
			return []string{"yaml", "json"}, cobra.ShellCompDirectiveFilterFileExt
		},
	}
	return cmd
}

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "global-transform  (ref) (type)",
		Short:             "Get a global transform",
		Long:              getLongDoc,
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, args[0], args[1])
		},
		Args: cobra.ExactArgs(2), // the policy or table id
	}
	return cmd
}

func DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "global-transform  (ref) (type)",
		Short:             "delete a global transform",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, args[0], args[1])
		},
		Args: cobra.ExactArgs(1), // the policy or table id
	}
	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "data-policies",
		Short:             "List Datapolicies",
		Example:           listExample,
		Long:              listLongDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
	}
	return cmd
}
