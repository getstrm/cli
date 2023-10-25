package table

import (
	"buf.build/gen/go/getstrm/daps/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policiesv1alpha "buf.build/gen/go/getstrm/daps/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

// strings used in the cli

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_policiesv1alphagrpc.NewDataPolicyServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	platformId := util.GetStringAndErr(flags, "processing-platform")
	req := &data_policiesv1alpha.ListProcessingPlatformTablesRequest{
		Platform: &data_policiesv1alpha.DataPolicy_ProcessingPlatform{
			Id: platformId,
		},
	}
	response, err := client.ListProcessingPlatformTables(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
