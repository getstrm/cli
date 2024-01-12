package lineage

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/processing_platforms/v1alpha/processing_platformsv1alphagrpc"
	processing_platformsv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"context"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

var apiContext context.Context

var ppclient processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient

func SetupClient(ppclient_ processing_platformsv1alphagrpc.ProcessingPlatformsServiceClient, ctx context.Context) {
	apiContext = ctx
	ppclient = ppclient_
}

func get(cmd *cobra.Command, s string) error {
	flags := cmd.Flags()
	platformId, _ := flags.GetString(common.ProcessingPlatformFlag)
	if platformId != "" {
		req := &processing_platformsv1alpha.GetLineageRequest{
			PlatformId: platformId,
			Fqn:        s,
		}
		response, err := ppclient.GetLineage(apiContext, req)

		return common.Print(printer, err, response)
	}
	return nil
}
