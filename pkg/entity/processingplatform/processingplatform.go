package processingplatform

import (
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
)

// strings used in the cli

var apiContext context.Context

var client processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(clientConnection processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(_ *cobra.Command) {
	response, err := client.ListProcessingPlatforms(apiContext, &ListProcessingPlatformsRequest{})
	CliExit(err)
	printer.Print(response)
}

func IdsCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	req := &ListProcessingPlatformsRequest{}
	response, err := client.ListProcessingPlatforms(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := lo.Map(response.ProcessingPlatforms, func(p *DataPolicy_ProcessingPlatform, _ int) string {
		return p.Id
	})
	log.Debugln(fmt.Sprintf("IdsCompletion names: %v", names))

	return names, cobra.ShellCompDirectiveNoFileComp
}

func AddProcessingPlatformFlag(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.StringP(common.ProcessingPlatformFlag, common.ProcessingPlatformFlagShort, "", common.ProcessingPlatformFlagUsage)
	CliExit(cmd.RegisterFlagCompletionFunc(common.ProcessingPlatformFlag, IdsCompletion))
}
