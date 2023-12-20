package completion

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	"pace/pace/pkg/entity/catalog"
	"pace/pace/pkg/entity/processingplatform"
	"pace/pace/pkg/util"
)

func AddCatalogFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", common.CatalogFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.CatalogFlag, catalog.IdsCompletion)
}

func AddProcessingPlatformFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", common.ProcessingPlatformFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, processingplatform.IdsCompletion)
}

func AddDatabaseFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", common.DatabaseFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.DatabaseFlag, databaseCompletion)
}

func AddSchemaFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.SchemaFlag, common.SchemaFlagShort, "", common.SchemaFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.SchemaFlag, schemaCompletion)
}

/*
	schemaCompletion

Provides completions based on whether or not we're connecting to a processing platform or a catalog.
*/
func schemaCompletion(cmd *cobra.Command, args []string, s string) ([]string, cobra.ShellCompDirective) {
	flags := cmd.Flags()
	ppId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if ppId != "" {
		return processingplatform.SchemaIdsCompletion(cmd, args, s)

	} else {
		return catalog.SchemaIdsCompletion(cmd, args, s)
	}
}

/*
	databaseCompletion

Provides completions based on whether or not we're connecting to a processing platform or a catalog.
*/
func databaseCompletion(cmd *cobra.Command, args []string, s string) ([]string, cobra.ShellCompDirective) {
	flags := cmd.Flags()
	ppId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if ppId != "" {
		return processingplatform.DatabaseIdsCompletion(cmd, args, s)

	} else {
		return catalog.DatabaseIdsCompletion(cmd, args, s)
	}
}
