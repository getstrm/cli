package database

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var apiContext context.Context
var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_policiesv1alphagrpc.NewDataPolicyServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	catalogId := util.GetStringAndErr(flags, common.CatalogFlag)
	req := &data_policiesv1alpha.ListDatabasesRequest{
		CatalogId: catalogId,
	}
	response, err := client.ListDatabases(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
