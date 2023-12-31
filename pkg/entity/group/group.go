package group

import (
	processingplatforms "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

// strings used in the cli

var apiContext context.Context

var client processingplatforms.ProcessingPlatformsServiceClient

func SetupClient(clientConnection processingplatforms.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func list(cmd *cobra.Command) error {
	v, _ := cmd.Flags().GetString(common.ProcessingPlatformFlag)
	platformId := v
	response, err := client.ListGroups(apiContext, &ListGroupsRequest{
		PlatformId:     platformId,
		PageParameters: common.PageParameters(cmd),
	})
	return common.Print(printer, err, response)
}
