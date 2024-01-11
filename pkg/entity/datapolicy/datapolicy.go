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
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"pace/pace/pkg/common"
	"sigs.k8s.io/yaml"
	"strings"
)

const FqnFlag = "fqn"

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

func upsert(cmd *cobra.Command, filename *string) error {
	policy, err := readPolicy(*filename)
	if err != nil {
		return err
	}
	applyFlag, _ := cmd.Flags().GetBool(common.ApplyFlag)
	req := &UpsertDataPolicyRequest{
		DataPolicy: policy,
		Apply:      applyFlag,
	}
	response, err := polClient.UpsertDataPolicy(apiContext, req)
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

func get(cmd *cobra.Command, dataPolicyOrTableId *string) error {
	flags := cmd.Flags()
	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	blueprint, _ := flags.GetBool(common.BlueprintFlag)
	if blueprint {
		// a blueprint policy only exists on processing platforms or catalogs
		if platformId != "" {
			return getBlueprintPolicyFromProcessingPlatform(flags, platformId, dataPolicyOrTableId)
		} else {
			return getBlueprintPolicyFromCatalog(flags, dataPolicyOrTableId)
		}
	} else {
		// return a data policy from the PACE database.
		response, err := getDataPolicy(dataPolicyOrTableId, platformId)
		return common.Print(printer, err, response.DataPolicy)
	}
}

func apply(cmd *cobra.Command, dataPolicyId *string) error {
	processingPlatform, _ := cmd.Flags().GetString(common.ProcessingPlatformFlag)
	req := &ApplyDataPolicyRequest{
		DataPolicyId: *dataPolicyId,
		PlatformId:   processingPlatform,
	}
	response, err := polClient.ApplyDataPolicy(apiContext, req)
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

func evaluate(cmd *cobra.Command, dataPolicyId *string) error {
	processingPlatform, _ := cmd.Flags().GetString(common.ProcessingPlatformFlag)
	sampleDataFileName, _ := cmd.Flags().GetString(common.SampleDataFlag)
	sampleDataFile, err := os.ReadFile(sampleDataFileName)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

func getDataPolicy(dataPolicyOrTableId *string, platformId string) (*GetDataPolicyResponse, error) {
	// return a data policy from the PACE database.
	req := &GetDataPolicyRequest{
		DataPolicyId: *dataPolicyOrTableId,
		PlatformId:   platformId,
	}
	return polClient.GetDataPolicy(apiContext, req)
}

func getBlueprintPolicyFromCatalog(flags *pflag.FlagSet, tableId *string) error {
	catalogId, databaseId, schemaId, _ := common.GetCatalogCoordinates(flags)
	req := &catalogentities.GetBlueprintPolicyRequest{
		CatalogId:  catalogId,
		DatabaseId: &databaseId,
		SchemaId:   &schemaId,
		TableId:    *tableId,
	}
	response, err := catClient.GetBlueprintPolicy(apiContext, req)
	if err != nil {
		return err
	}
	return common.Print(printer, err, response.DataPolicy)
}

func getBlueprintPolicyFromProcessingPlatform(flags *pflag.FlagSet, platformId string, tableId *string) error {
	_, databaseId, schemaId, _ := common.GetCatalogCoordinates(flags)
	fqn, _ := flags.GetBool(FqnFlag)
	fqnValue := ""
	if fqn {
		fqnValue = *tableId
	}
	req := &ppentities.GetBlueprintPolicyRequest{
		PlatformId: platformId,
		Table: &Table{
			Id: *tableId,
			Schema: &Schema{
				Id: schemaId,
				Database: &Database{
					Id: databaseId,
				},
			},
		},
		Fqn: fqnValue,
	}
	response, err := pClient.GetBlueprintPolicy(apiContext, req)
	if err != nil {
		return err
	}
	if response.Violation != nil && response.Violation.Description != "" {
		message := fmt.Sprintf("Bare policy violation: %s\n", response.Violation.Description)
		return errors.New(message)
	}
	return common.Print(printer, err, response.DataPolicy)
}

func list(cmd *cobra.Command) error {
	req := &ListDataPoliciesRequest{
		PageParameters: common.PageParameters(cmd),
	}
	response, err := polClient.ListDataPolicies(apiContext, req)
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

// readPolicy
// read a json or yaml encoded policy from the filesystem.
func readPolicy(filename string) (*DataPolicy, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// check if the file is yaml and convert it to json
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		if err != nil {
			return nil, err
		}
	}
	dataPolicy := &DataPolicy{}
	err = protojson.Unmarshal(file, dataPolicy)
	return dataPolicy, err
}

func TableOrDataPolicyIdsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags := cmd.Flags()
	blueprint, _ := flags.GetBool(common.BlueprintFlag)

	// talking to the PACE database
	if !blueprint {
		return idsCompletion(cmd, args, toComplete)
	}

	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	catalogId, databaseId, schemaId, _ := common.GetCatalogCoordinates(flags)
	if databaseId == "" || schemaId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if platformId != "" {
		response, err := pClient.ListTables(apiContext, &ppentities.ListTablesRequest{
			PlatformId: platformId,
			DatabaseId: &databaseId,
			SchemaId:   &schemaId,
		})

		if err != nil {
			return common.CobraCompletionError(err)
		}
		return lo.Map(response.Tables, func(table *Table, _ int) string {
			return table.Id
		}), cobra.ShellCompDirectiveNoFileComp
	}

	response, err := catClient.ListTables(apiContext, &catalogentities.ListTablesRequest{
		CatalogId:  catalogId,
		DatabaseId: &databaseId,
		SchemaId:   &schemaId,
	})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	return lo.Map(response.Tables, func(table *Table, _ int) string {
		return table.Id
	}), cobra.ShellCompDirectiveNoFileComp
}

/*
idsCompletion.

returns data policy ids optionally limited to a certain platform.
*/
func idsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := polClient.ListDataPolicies(apiContext, &ListDataPoliciesRequest{})
	if err != nil {
		return common.CobraCompletionError(err)
	}

	var policies = response.DataPolicies

	platformId, _ := cmd.Flags().GetString(common.ProcessingPlatformFlag)

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
