package database

import (
	catalogs "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

var apiContext context.Context

var client catalogs.DataCatalogsServiceClient

func SetupClient(clientConnection catalogs.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) error {
	v, _ := cmd.Flags().GetString(common.CatalogFlag)
	catalogId := v
	response, err := client.ListDatabases(apiContext, &ListDatabasesRequest{
		CatalogId:      catalogId,
		PageParameters: common.PageParameters(cmd),
	})

	if err != nil {
		return err
	}

	return common.Print(printer, err, response)
	return nil

}
