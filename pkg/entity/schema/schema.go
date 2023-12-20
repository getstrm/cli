package schema

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	processing_platformsv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
)

// strings used in the cli
var apiContext context.Context

var ppclient processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient
var catclient data_catalogsv1alphagrpc.DataCatalogsServiceClient

func SetupClient(ppclient_ processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient, catclient_ data_catalogsv1alphagrpc.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	ppclient = ppclient_
	catclient = catclient_
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	catalogId := GetStringAndErr(flags, common.CatalogFlag)
	databaseId := GetStringAndErr(flags, common.DatabaseFlag)
	platformId := GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &processing_platformsv1alpha.ListSchemasRequest{
			PlatformId:     platformId,
			DatabaseId:     &databaseId,
			PageParameters: common.PageParameters(cmd),
		}
		response, err := ppclient.ListSchemas(apiContext, req)
		CliExit(err)
		printer.Print(response)
	} else {
		req := &data_catalogsv1alpha.ListSchemasRequest{
			CatalogId:      catalogId,
			DatabaseId:     &databaseId,
			PageParameters: common.PageParameters(cmd),
		}
		response, err := catclient.ListSchemas(apiContext, req)
		CliExit(err)
		printer.Print(response)
	}
}
