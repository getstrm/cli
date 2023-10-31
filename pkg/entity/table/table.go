package table

import (
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var apiContext context.Context

var ppclient ProcessingPlatformsServiceClient
var catclient DataCatalogsServiceClient

func SetupClient(ppclient_ ProcessingPlatformsServiceClient, catclient_ DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	ppclient = ppclient_
	catclient = catclient_
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	platformId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &processing_platformsv1alpha.ListTablesRequest{
			PlatformId: platformId,
		}
		response, err := ppclient.ListTables(apiContext, req)
		util.CliExit(err)
		printer.Print(response)
	} else {
		catalogId, databaseId, schemaId := common.GetCatalogCoordinates(flags)
		req := &data_catalogsv1alpha.ListTablesRequest{
			CatalogId:  catalogId,
			DatabaseId: databaseId,
			SchemaId:   schemaId,
		}
		response, err := catclient.ListTables(apiContext, req)
		util.CliExit(err)
		printer.Print(response)

	}
}
