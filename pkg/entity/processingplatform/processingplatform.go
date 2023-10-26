package processingplatform

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"pace/pace/pkg/common"
)

// strings used in the cli

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_policiesv1alphagrpc.NewDataPolicyServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	req := &data_policiesv1alpha.ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
func PlatformIdsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &data_policiesv1alpha.ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.ProcessingPlatforms))
	for _, s := range response.ProcessingPlatforms {
		names = append(names, s.Id)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
