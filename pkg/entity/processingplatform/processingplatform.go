package processingplatform

import (
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

// strings used in the cli

var apiContext context.Context

var client processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(clientConnection processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(_ *cobra.Command) error {
	response, err := client.ListProcessingPlatforms(apiContext, &ListProcessingPlatformsRequest{})
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

func IdsCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	req := &ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	if err != nil {
		return common.CobraCompletionError(err)
	}
	return lo.Map(response.ProcessingPlatforms, func(p *ProcessingPlatform, _ int) string {
		return p.Id
	}), cobra.ShellCompDirectiveNoFileComp
}

func DatabaseIdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags := cmd.Flags()
	ppId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if ppId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := client.ListDatabases(apiContext, &ListDatabasesRequest{
		PlatformId:     ppId,
		PageParameters: common.PageParameters(cmd),
	})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	names := lo.Map(response.Databases, func(db *Database, _ int) string {
		return db.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func SchemaIdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	flags := cmd.Flags()
	ppId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if ppId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	databaseId, _ := flags.GetString(common.DatabaseFlag)
	if databaseId == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	response, err := client.ListSchemas(apiContext, &ListSchemasRequest{
		PlatformId:     ppId,
		DatabaseId:     &databaseId,
		PageParameters: common.PageParameters(cmd),
	})
	if err != nil {
		return common.CobraCompletionError(err)
	}
	names := lo.Map(response.Schemas, func(catalog *Schema, _ int) string {
		return catalog.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}
