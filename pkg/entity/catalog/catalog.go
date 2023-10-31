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
	"pace/pace/pkg/util"
)

var apiContext context.Context
var client catalogs.DataCatalogsServiceClient

func SetupClient(clientConnection catalogs.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list() {
	req := &ListCatalogsRequest{}
	response, err := client.ListCatalogs(apiContext, req)
	util.CliExit(err)
	printer.Print(response)
}

func IdsCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := client.ListCatalogs(apiContext, &ListCatalogsRequest{})
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}
	names := lo.Map(response.Catalogs, func(catalog *DataCatalog, _ int) string {
		return catalog.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func AddCatalogFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", common.CatalogFlagUsage)
	err := cmd.RegisterFlagCompletionFunc(common.CatalogFlag, IdsCompletion)
	util.CliExit(err)
}

func AddDatabaseFlag(flags *pflag.FlagSet) {
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", common.DatabaseFlagUsage)
}

func AddSchemaFlag(flags *pflag.FlagSet) {
	flags.StringP(common.SchemaFlag, common.SchemaFlagShort, "", common.SchemaFlagUsage)
}
