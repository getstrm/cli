package processingplatform

import (
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
)

// strings used in the cli

var apiContext context.Context

var client processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(clientConnection processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(_ *cobra.Command) {
	req := &ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func IdsCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := lo.Map(response.ProcessingPlatforms, func(p *DataPolicy_ProcessingPlatform, _ int) string {
		return p.Id
	})
	return names, cobra.ShellCompDirectiveNoFileComp
}

func AddProcessingPlatformFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", common.ProcessingPlatformFlagUsage)
	err := cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, IdsCompletion)
	common.CliExit(err)
}
