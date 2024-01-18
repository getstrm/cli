package resources

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/resources/v1alpha/resourcesv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/resources/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"strings"
)

var apiContext context.Context

var resources_client resourcesv1alphagrpc.ResourcesServiceClient

func SetupClient(resources_client_ resourcesv1alphagrpc.ResourcesServiceClient, ctx context.Context) {
	apiContext = ctx
	resources_client = resources_client_
}

func list(cmd *cobra.Command, s []string) error {
	req := &ListResourcesRequest{
		PageParameters: common.PageParameters(cmd),
	}

	if len(s) > 0 {
		resourcePath := lo.Map(strings.Split(s[0], "/"), func(pp string, _ int) *ResourceNode {
			return &ResourceNode{Name: pp}
		})
		req.Urn = &ResourceUrn{ResourcePath: resourcePath}
	}
	response, err := resources_client.ListResources(apiContext, req)
	return common.Print(printer, err, response)

}
