package catalog

import (
	catalogs "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
)

var apiContext context.Context
var client catalogs.DataCatalogsServiceClient

func SetupClient(clientConnection catalogs.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list() error {
	response, err := client.ListCatalogs(apiContext, &ListCatalogsRequest{})
	return common.Print(printer, err, response)
}

func IdsCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	response, err := client.ListCatalogs(apiContext, &ListCatalogsRequest{})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	return lo.Map(response.Catalogs, func(catalog *DataCatalog, _ int) string {
		return catalog.Id
	}), cobra.ShellCompDirectiveNoFileComp
}

func DatabaseIdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags := cmd.Flags()
	catalogId, _ := flags.GetString(common.CatalogFlag)
	if catalogId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := client.ListDatabases(apiContext, &ListDatabasesRequest{
		CatalogId: catalogId,
	})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	return lo.Map(response.Databases, func(catalog *DataCatalog_Database, _ int) string {
		return catalog.Id
	}), cobra.ShellCompDirectiveNoFileComp
}

func SchemaIdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags := cmd.Flags()
	catalogId, _ := flags.GetString(common.CatalogFlag)
	if catalogId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	databaseId, _ := flags.GetString(common.DatabaseFlag)
	if databaseId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := client.ListSchemas(apiContext, &ListSchemasRequest{
		CatalogId:  catalogId,
		DatabaseId: &databaseId,
	})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	return lo.Map(response.Schemas, func(catalog *DataCatalog_Schema, _ int) string {
		return catalog.Id
	}), cobra.ShellCompDirectiveNoFileComp
}

func AddCatalogFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", common.CatalogFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.CatalogFlag, IdsCompletion)
}

func AddDatabaseFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", common.DatabaseFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.DatabaseFlag, DatabaseIdsCompletion)
}

func AddSchemaFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.SchemaFlag, common.SchemaFlagShort, "", common.SchemaFlagUsage)
	_ = cmd.RegisterFlagCompletionFunc(common.SchemaFlag, SchemaIdsCompletion)
}
