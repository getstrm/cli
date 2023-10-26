package schema

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

// strings used in the cli

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_policiesv1alphagrpc.NewDataPolicyServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	catalogId := util.GetStringAndErr(flags, common.CatalogFlag)
	databaseId := util.GetStringAndErr(flags, common.DatabaseFlag)
	req := &data_policiesv1alpha.ListSchemasRequest{
		CatalogId:  catalogId,
		DatabaseId: databaseId,
	}
	response, err := client.ListSchemas(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
