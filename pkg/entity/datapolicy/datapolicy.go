package datapolicy

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policies "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
	"strings"
)

// strings used in the cli

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_policiesv1alphagrpc.NewDataPolicyServiceClient(clientConnection)
}

func upsert(cmd *cobra.Command, filename *string) {
	policy := readPolicy(*filename)
	req := &data_policies.UpsertDataPolicyRequest{
		DataPolicy: policy,
	}
	response, err := client.UpsertDataPolicy(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(cmd *cobra.Command, tableId *string) {
	flags := cmd.Flags()
	if util.GetBoolAndErr(flags, bareFlag) {
		getBare(cmd, tableId)
	} else {
		req := &data_policies.GetDataPolicyRequest{
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
	req := &data_policies.GetCatalogBarePolicyRequest{
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
	req := &data_policies.GetProcessingPlatformBarePolicyRequest{
		PlatformId: platformId,
		Table:      *tableId,
	}
	response, err := client.GetProcessingPlatformBarePolicy(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func list(cmd *cobra.Command) {
	req := &data_policies.ListDataPoliciesRequest{}
	response, err := client.ListDataPolicies(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
func readPolicy(filename string) *data_policies.DataPolicy {
	file, _ := os.ReadFile(filename)
	dataPolicy := &data_policies.DataPolicy{}

	var err error
	if strings.HasSuffix(filename, ".json") {
		err = protojson.Unmarshal(file, dataPolicy)
	} else {
		err = yaml.Unmarshal(file, dataPolicy)
	}
	common.CliExit(err)
	return dataPolicy
}
