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

func upsert(_ *cobra.Command, filename *string) {
	policy := readPolicy(*filename)
	req := &UpsertDataPolicyRequest{
		DataPolicy: policy,
	}
	response, err := polClient.UpsertDataPolicy(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func get(cmd *cobra.Command, tableId *string) {
	flags := cmd.Flags()
	platformId, _, bare := isBlueprint(flags)
	if bare {
		// a bare policy only exists on processing platforms or catalogs
		if platformId != "" {
			getBlueprintPolicyFromProcessingPlatform(platformId, tableId)
		} else {
			getBlueprintPolicyFromCatalog(flags, tableId)
		}
	} else {
		// return a data policy from the Pace database.
		req := &GetDataPolicyRequest{
			DataPolicyId: *tableId,
			PlatformId:   platformId,
		}
		response, err := polClient.GetDataPolicy(apiContext, req)
		CliExit(err)
		printer.Print(response.DataPolicy)
	}
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

func list(_ *cobra.Command) {
	req := &ListDataPoliciesRequest{}
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

func TableIdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	flags := cmd.Flags()

	platformId, catalogId, bare := isBlueprint(flags)
	// talking to the Pace database
	if !bare {
		response, err := polClient.ListDataPolicies(apiContext, &ListDataPoliciesRequest{})
		CliExit(err)
		return lo.Map(response.DataPolicies, func(table *DataPolicy, _ int) string {
			return table.Id
		}), cobra.ShellCompDirectiveNoFileComp
	}

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
	names := lo.Map(response.Tables, func(table *DataCatalog_Table, _ int) string {
		return table.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

/*
	isBlueprint

Reads catalog and platform flags and determines if this is
a call for a bare data policy.
*/
func isBlueprint(flags *pflag.FlagSet) (string, string, bool) {
	platformId := GetStringAndErr(flags, common.ProcessingPlatformFlag)
	catalogId := GetStringAndErr(flags, common.CatalogFlag)
	return platformId, catalogId, platformId != "" || catalogId != ""
}
