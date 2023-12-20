package database

import (
	catalogs "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/data_catalogs/v1alpha/data_catalogsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var apiContext context.Context

var client catalogs.DataCatalogsServiceClient

func SetupClient(clientConnection catalogs.DataCatalogsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) {
	catalogId := util.GetStringAndErr(cmd.Flags(), common.CatalogFlag)
	response, err := client.ListDatabases(apiContext, &ListDatabasesRequest{
		CatalogId:      catalogId,
		PageParameters: common.PageParameters(cmd),
	})
	util.CliExit(err)
	printer.Print(response)
}
