package schema

import (
	catalogs "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

// strings used in the cli
var apiContext context.Context

var client catalogs.DataCatalogsServiceClient

func SetupClient(clientConnection catalogs.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) error {
	flags := cmd.Flags()
	v, _ := flags.GetString(common.CatalogFlag)
	catalogId := v
	v2, _ := flags.GetString(common.DatabaseFlag)
	databaseId := v2
	req := &ListSchemasRequest{
		CatalogId:      catalogId,
		DatabaseId:     &databaseId,
		PageParameters: common.PageParameters(cmd),
	}
	response, err := client.ListSchemas(apiContext, req)

	if err != nil {
		return err
	}

	return common.Print(printer, err, response)
}
