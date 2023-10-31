package datapolicy

import (
	catalogs "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	datapolicies "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	catalogentities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_policies/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	ppentities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
	"sigs.k8s.io/yaml"
	"strings"
)

// strings used in the cli

var apiContext context.Context

var polClient datapolicies.DataPoliciesServiceClient
var catClient catalogs.DataCatalogsServiceClient
var pClient processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(policiesServiceClient datapolicies.DataPoliciesServiceClient, catalogsClient catalogs.DataCatalogsServiceClient, ppClient processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	polClient = policiesServiceClient
	catClient = catalogsClient
	pClient = ppClient
}

func upsert(_ *cobra.Command, filename *string) {
	policy := readPolicy(*filename)
	req := &UpsertDataPolicyRequest{
		DataPolicy: policy,
	}
	response, err := polClient.UpsertDataPolicy(apiContext, req)
	util.CliExit(err)
	printer.Print(response)
}

func get(cmd *cobra.Command, tableId *string) {
	flags := cmd.Flags()
	if util.GetBoolAndErr(flags, bareFlag) {
		getBare(cmd, tableId)
	} else {
		req := &GetDataPolicyRequest{
			DataPolicyId: *tableId,
		}
		response, err := polClient.GetDataPolicy(apiContext, req)
		util.CliExit(err)
		printer.Print(response.DataPolicy)
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
	catalogId, databaseId, schemaId := common.GetCatalogCoordinates(flags)
	req := &catalogentities.GetBarePolicyRequest{
		CatalogId:  catalogId,
		DatabaseId: databaseId,
		SchemaId:   schemaId,
		TableId:    *tableId,
	}
	response, err := catClient.GetBarePolicy(apiContext, req)
	util.CliExit(err)
	printer.Print(response.DataPolicy)
}

func getBarePolicyFromProcessingPlatform(platformId string, tableId *string) {
	req := &ppentities.GetBarePolicyRequest{
		PlatformId: platformId,
		Table:      *tableId,
	}
	response, err := pClient.GetBarePolicy(apiContext, req)
	util.CliExit(err)
	printer.Print(response.DataPolicy)
}

func list(_ *cobra.Command) {
	req := &ListDataPoliciesRequest{}
	response, err := polClient.ListDataPolicies(apiContext, req)
	util.CliExit(err)
	printer.Print(response)
}
func readPolicy(filename string) *DataPolicy {
	file, err := os.ReadFile(filename)
	util.CliExit(err)

	if strings.HasSuffix(filename, ".yaml") {
		file, _ = yaml.YAMLToJSON(file)
	}
	dataPolicy := &DataPolicy{}
	err = protojson.Unmarshal(file, dataPolicy)
	util.CliExit(err)
	return dataPolicy
}
