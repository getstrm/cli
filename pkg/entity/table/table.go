package table

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

var apiContext context.Context

var ppclient processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient
var catclient data_catalogsv1alphagrpc.DataCatalogsServiceClient

func SetupClient(ppclient_ processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient, catclient_ data_catalogsv1alphagrpc.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	ppclient = ppclient_
	catclient = catclient_
}

func list(cmd *cobra.Command) error {
	flags := cmd.Flags()
	catalogId, databaseId, schemaId, _ := common.GetCatalogCoordinates(flags)
	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &processing_platformsv1alpha.ListTablesRequest{
			PlatformId:     platformId,
			DatabaseId:     &databaseId,
			SchemaId:       &schemaId,
			PageParameters: common.PageParameters(cmd),
		}
		response, err := ppclient.ListTables(apiContext, req)
		return common.Print(printer, err, response)
	} else {
		req := &data_catalogsv1alpha.ListTablesRequest{
			CatalogId:      catalogId,
			DatabaseId:     &databaseId,
			SchemaId:       &schemaId,
			PageParameters: common.PageParameters(cmd),
		}
		response, err := catclient.ListTables(apiContext, req)
		return common.Print(printer, err, response)
	}
}
