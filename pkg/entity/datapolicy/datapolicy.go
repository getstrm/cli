package datapolicy

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

func get(cmd *cobra.Command, tableId *string) {
	flags := cmd.Flags()
	if util.GetBoolAndErr(flags, bareFlag) {
		getBare(cmd, tableId)
	} else {
		req := &data_policiesv1alpha.GetDataPolicyRequest{
			DataPolicyId: *tableId,
		}
		response, err := client.GetDataPolicy(apiContext, req)
		common.CliExit(err)
		printer.Print(response)
	}
}

func getBare(cmd *cobra.Command, tableId *string) {
	flags := cmd.Flags()
	platformId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if platformId != "" {
		getBarePolicyFromProcessingPlatform(platformId, tableId)
	} else {
		getBarePolicyFromCatalog(flags, tableId)
	}
}

func getBarePolicyFromCatalog(flags *pflag.FlagSet, tableId *string) {
	catalogId, databaseId, schemaId := common.CheckCatalogCoords(flags)
	req := &data_policiesv1alpha.GetCatalogBarePolicyRequest{
		CatalogId:  catalogId,
		DatabaseId: databaseId,
		SchemaId:   schemaId,
		TableId:    *tableId,
	}
	response, err := client.GetCatalogBarePolicy(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func getBarePolicyFromProcessingPlatform(platformId string, tableId *string) {
	req := &data_policiesv1alpha.GetProcessingPlatformBarePolicyRequest{
		PlatformId: platformId,
		Table:      *tableId,
	}
	response, err := client.GetProcessingPlatformBarePolicy(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
func list(cmd *cobra.Command) {
	req := &data_policiesv1alpha.ListDataPoliciesRequest{}
	response, err := client.ListDataPolicies(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
