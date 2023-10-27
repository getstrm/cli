package table

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	datapolicies "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection data_policiesv1alphagrpc.DataPolicyServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	platformId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &datapolicies.ListProcessingPlatformTablesRequest{
			PlatformId: platformId,
		}
		response, err := client.ListProcessingPlatformTables(apiContext, req)
		common.CliExit(err)
		printer.Print(response)
	} else {
		catalogId, databaseId, schemaId := common.GetCatalogCoordinates(flags)
		req := &datapolicies.ListTablesRequest{
			CatalogId:  catalogId,
			DatabaseId: databaseId,
			SchemaId:   schemaId,
		}
		response, err := client.ListTables(apiContext, req)
		common.CliExit(err)
		printer.Print(response)

	}
}
