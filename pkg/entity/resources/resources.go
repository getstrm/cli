package resources

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/resources/v1alpha/resourcesv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/resources/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"strings"
)

var apiContext context.Context

var resourcesClient resourcesv1alphagrpc.ResourcesServiceClient

func SetupClient(resources_client_ resourcesv1alphagrpc.ResourcesServiceClient, ctx context.Context) {
	apiContext = ctx
	resourcesClient = resources_client_
}

func list(cmd *cobra.Command, s []string) error {
	req := &ListResourcesRequest{
		PageParameters: common.PageParameters(cmd),
	}

	if len(s) > 0 {
		resourcePathElements := strings.Split(s[0], "/")
		req.IntegrationId = resourcePathElements[0]
		req.ResourcePath = lo.Drop(resourcePathElements, 1)
	}

	response, err := resourcesClient.ListResources(apiContext, req)
	return common.Print(printer, err, response)

}
