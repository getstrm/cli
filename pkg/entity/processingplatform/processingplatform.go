package processingplatform

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/api/data_policies/v1alpha/data_policiesv1alphagrpc"
	datapolicies "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
)

// strings used in the cli

var apiContext context.Context

var client data_policiesv1alphagrpc.DataPolicyServiceClient

func SetupClient(clientConnection data_policiesv1alphagrpc.DataPolicyServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) {
	req := &datapolicies.ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
func IdsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &datapolicies.ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := lo.Map(response.ProcessingPlatforms, func(p *datapolicies.DataPolicy_ProcessingPlatform, _ int) string {
		return p.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func AddProcessingPlatformFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", common.ProcessingPlatformFlagUsage)
	err := cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, IdsCompletion)
	common.CliExit(err)
}
