package table

import (
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

var apiContext context.Context

var ppclient ProcessingPlatformsServiceClient
var catclient DataCatalogsServiceClient

func SetupClient(ppclient_ ProcessingPlatformsServiceClient, catclient_ DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	ppclient = ppclient_
	catclient = catclient_
}

func list(cmd *cobra.Command) error {
	flags := cmd.Flags()
	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &processing_platformsv1alpha.ListTablesRequest{
			PlatformId:     platformId,
			PageParameters: common.PageParameters(cmd),
		}
		response, err := ppclient.ListTables(apiContext, req)
		if err != nil {
			return err
		}
		return common.Print(printer, err, response)
	} else {
		catalogId, databaseId, schemaId, err := common.GetCatalogCoordinates(flags)
		if err != nil {
			return err
		}
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
