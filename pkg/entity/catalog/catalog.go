package catalog

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	datapolicies "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
)

var apiContext context.Context
var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection data_policiesv1alphagrpc.DataPolicyServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list() {
	response, err := client.ListCatalogs(apiContext, &datapolicies.ListCatalogsRequest{})
	common.CliExit(err)
	printer.Print(response)
}

func IdsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &datapolicies.ListCatalogsRequest{}
	response, err := client.ListCatalogs(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}
	names := lo.Map(response.Catalogs, func(catalog *datapolicies.DataCatalog, _ int) string {
		return catalog.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func AddCatalogFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.CatalogFlag, common.CatalogFlagShort, "", common.CatalogFlagUsage)
	err := cmd.RegisterFlagCompletionFunc(common.CatalogFlag, IdsCompletion)
	common.CliExit(err)
}
func AddDatabaseFlag(flags *pflag.FlagSet) {
	flags.StringP(common.DatabaseFlag, common.DatabaseFlagShort, "", common.DatabaseFlagUsage)
}
func AddSchemaFlag(flags *pflag.FlagSet) {
	flags.StringP(common.SchemaFlag, common.SchemaFlagShort, "", common.SchemaFlagUsage)
}
