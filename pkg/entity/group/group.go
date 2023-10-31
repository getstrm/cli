package group

import (
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

// strings used in the cli

var apiContext context.Context

var client processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(clientConnection processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) {
	flags := cmd.Flags()
	platformId := util.GetStringAndErr(flags, common.ProcessingPlatformFlag)
	req := &ListGroupsRequest{
		PlatformId: platformId,
	}
	response, err := client.ListGroups(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}
