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
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
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

func upsert(cmd *cobra.Command, filename *string) {
	policy := readPolicy(*filename)
	apply := GetBoolAndErr(cmd.Flags(), common.ApplyFlag)
	req := &UpsertDataPolicyRequest{
		DataPolicy: policy,
		Apply:      apply,
	}
	response, err := polClient.UpsertDataPolicy(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func get(cmd *cobra.Command, dataPolicyOrTableId *string) {
	flags := cmd.Flags()

	platformId := GetStringAndErr(flags, common.ProcessingPlatformFlag)
	blueprint := GetBoolAndErr(flags, common.BlueprintFlag)
	if blueprint {
		// a blueprint policy only exists on processing platforms or catalogs
		if platformId != "" {
			getBlueprintPolicyFromProcessingPlatform(platformId, dataPolicyOrTableId)
		} else {
			getBlueprintPolicyFromCatalog(flags, dataPolicyOrTableId)
		}
	} else {
		// return a data policy from the PACE database.
		response := getDataPolicy(dataPolicyOrTableId, platformId)
		printer.Print(response.DataPolicy)
	}
}

func apply(cmd *cobra.Command, dataPolicyId *string) *ApplyDataPolicyResponse {
	processingPlatform := GetStringAndErr(cmd.Flags(), common.ProcessingPlatformFlag)

	req := &ApplyDataPolicyRequest{
		DataPolicyId: *dataPolicyId,
		PlatformId:   processingPlatform,
	}
	response, err := polClient.ApplyDataPolicy(apiContext, req)
	CliExit(err)

	return response
}

func evaluate(cmd *cobra.Command, dataPolicyId *string) {
	processingPlatform := GetStringAndErr(cmd.Flags(), common.ProcessingPlatformFlag)
	sampleDataFileName := GetStringAndErr(cmd.Flags(), common.SampleDataFlag)
	sampleDataFile, err := os.ReadFile(sampleDataFileName)
	CliExit(err)
	sampleData := string(sampleDataFile)

	req := &EvaluateDataPolicyRequest{
		DataPolicyId: *dataPolicyId,
		PlatformId:   processingPlatform,
		Evaluation: &EvaluateDataPolicyRequest_FullEvaluation_{
			FullEvaluation: &EvaluateDataPolicyRequest_FullEvaluation{
				SampleCsv: sampleData,
			},
		},
	}
	response, err := polClient.EvaluateDataPolicy(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func getDataPolicy(dataPolicyOrTableId *string, platformId string) *GetDataPolicyResponse {
	// return a data policy from the PACE database.
	req := &GetDataPolicyRequest{
		DataPolicyId: *dataPolicyOrTableId,
		PlatformId:   platformId,
	}
	response, err := polClient.GetDataPolicy(apiContext, req)
	CliExit(err)
	return response
}

func getBlueprintPolicyFromCatalog(flags *pflag.FlagSet, tableId *string) {
	catalogId, databaseId, schemaId := common.GetCatalogCoordinates(flags)
	req := &catalogentities.GetBlueprintPolicyRequest{
		CatalogId:  catalogId,
		DatabaseId: &databaseId,
		SchemaId:   &schemaId,
		TableId:    *tableId,
	}
	response, err := catClient.GetBlueprintPolicy(apiContext, req)
	CliExit(err)
	printer.Print(response.DataPolicy)
}

func getBlueprintPolicyFromProcessingPlatform(platformId string, tableId *string) {
	req := &ppentities.GetBlueprintPolicyRequest{
		PlatformId: platformId,
		TableId:    *tableId,
	}
	response, err := pClient.GetBlueprintPolicy(apiContext, req)
	CliExit(err)
	printer.Print(response.DataPolicy)
	if response.Violation != nil && response.Violation.Description != "" {
		fmt.Fprintf(os.Stderr, "Bare policy violation: %s\n", response.Violation.Description)
		os.Exit(10)
	}
}

func list(cmd *cobra.Command) {
	req := &ListDataPoliciesRequest{
		PageParameters: common.PageParameters(cmd),
	}
	response, err := polClient.ListDataPolicies(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

// readPolicy
// read a json or yaml encoded policy from the filesystem.
func readPolicy(filename string) *DataPolicy {
	file, err := os.ReadFile(filename)
	CliExit(err)

	// check if the file is yaml and convert it to json
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		CliExit(err)
	}

	dataPolicy := &DataPolicy{}
	err = protojson.Unmarshal(file, dataPolicy)
	CliExit(err)
	return dataPolicy
}

func TableOrDataPolicyIdsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	flags := cmd.Flags()

	blueprint := GetBoolAndErr(flags, common.BlueprintFlag)

	// talking to the PACE database
	if !blueprint {
		return IdsCompletion(cmd, args, toComplete)
	}

	platformId := GetStringAndErr(flags, common.ProcessingPlatformFlag)
	catalogId := GetStringAndErr(flags, common.CatalogFlag)

	// talking to a processing platform
	if platformId != "" {
		response, err := pClient.ListTables(apiContext, &ppentities.ListTablesRequest{
			PlatformId: platformId,
		})
		CliExit(err)
		return response.Tables, cobra.ShellCompDirectiveNoFileComp
	}

	// talking to a catalog!
	if catalogId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	databaseId, _ := flags.GetString(common.DatabaseFlag)
	if databaseId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	schemaId, _ := flags.GetString(common.SchemaFlag)
	if schemaId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := catClient.ListTables(apiContext, &catalogentities.ListTablesRequest{
		CatalogId:  catalogId,
		DatabaseId: &databaseId,
		SchemaId:   &schemaId,
	})
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}
	names := lo.Map(response.Tables, func(table *Table, _ int) string {
		return table.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func IdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := polClient.ListDataPolicies(apiContext, &ListDataPoliciesRequest{})
	CliExit(err)
	var policies = response.DataPolicies

	platformId := GetStringAndErr(cmd.Flags(), common.ProcessingPlatformFlag)

	// If the platform id is provided, make sure we only suggest policies for that platform
	if platformId != "" {
		policies = lo.Filter(policies, func(policy *DataPolicy, _ int) bool {
			return policy.Platform.Id == platformId
		})
	}

	return lo.Map(policies, func(dataPolicy *DataPolicy, _ int) string {
		return dataPolicy.Id
	}), cobra.ShellCompDirectiveNoFileComp
}
