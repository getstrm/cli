package globaltransform

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	"strings"
)

const (
	policyTypeFlag      = "type"
	policyTypeFlagShort = "t"
)

var transformTypes = []string{"TAG_TRANSFORM"}

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
		Use:               "global-transform (ref)",
		Short:             "Get a global transform",
		Long:              getLongDoc,
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, args[0])
		},
		Args:              cobra.ExactArgs(1), // ref
		ValidArgsFunction: refCompletionFunction,
	}
	_ = setupFlags(cmd)
	return cmd
}

func DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "global-transform (ref)",
		Short:             "delete a global transform",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = common.ConfigurePrinter(cmd, common.StandardPrinters)
		},
		Run: func(cmd *cobra.Command, args []string) {
			delete(cmd, args[0])
		},
		Args:              cobra.ExactArgs(1), // ref and type
		ValidArgsFunction: refCompletionFunction,
	}
	_ = setupFlags(cmd)
	return cmd
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "global-transforms",
		Short:             "List Global Transforms",
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

func typeCompletionFunction(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return transformTypes, cobra.ShellCompDirectiveNoFileComp
}

func setupFlags(cmd *cobra.Command) *pflag.FlagSet {
	flags := cmd.Flags()
	_ = flags.StringP(policyTypeFlag, policyTypeFlagShort, transformTypes[0], "type of global transform: "+strings.Join(transformTypes, ","))
	_ = cmd.RegisterFlagCompletionFunc(policyTypeFlag, typeCompletionFunction)
	return flags
}
