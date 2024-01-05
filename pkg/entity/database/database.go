package database

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	data_catalogsv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	processing_platformsv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
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
	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if platformId != "" {
		response, err := ppclient.ListDatabases(apiContext, &processing_platformsv1alpha.ListDatabasesRequest{
			PlatformId:     platformId,
			PageParameters: common.PageParameters(cmd),
		})
		return common.Print(printer, err, response)
	} else {
		catalogId, _ := flags.GetString(common.CatalogFlag)
		response, err := catclient.ListDatabases(apiContext, &data_catalogsv1alpha.ListDatabasesRequest{
			CatalogId:      catalogId,
			PageParameters: common.PageParameters(cmd),
		})
		return common.Print(printer, err, response)
	}
}
