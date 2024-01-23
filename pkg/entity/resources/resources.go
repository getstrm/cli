package resources

import (
	"buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/resources/v1alpha/resourcesv1alphagrpc"
	entitiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/resources/v1alpha"
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func resourcesCompletion(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// pace list resources bigquery-dev/stream-machine-development/batch_job_demo

	req := &ListResourcesRequest{
		PageParameters: common.PageParameters(cmd),
	}

	var response *ListResourcesResponse
	var err error

	resourcePathElements := strings.Split(toComplete, "/")
	req.IntegrationId = resourcePathElements[0]
	resourcePath := lo.Drop(resourcePathElements, 1)
	last, _ := lo.Last(resourcePath)

	var resourceNames []string

	if last == "" {
		req.ResourcePath = lo.DropRight(resourcePath, 1)
	} else {
		// First try to get the resources for this path and treat the last element
		// as an existing resource that might have children
		req.ResourcePath = resourcePath
		response, err = resourcesClient.ListResources(apiContext, req)

		if err != nil {
			st, _ := status.FromError(err)

			if st.Code() == codes.NotFound {
				// If the resource doesn't exist, try to get the resources for the
				// parent path and treat the last element as a prefix
				req.ResourcePath = lo.DropRight(resourcePath, 1)
			} else {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		}
	}

	response, err = resourcesClient.ListResources(apiContext, req)

	noLeafResources := lo.Filter(response.Resources, func(resource *entitiesv1alpha.ResourceUrn, _ int) bool {
		if len(resource.ResourcePath) > 0 {
			lastResource, _ := lo.Last(resource.ResourcePath)
			return !lastResource.IsLeaf
		}

		return true
	})

	resourceNames = lo.Map(noLeafResources, func(resource *entitiesv1alpha.ResourceUrn, _ int) string {
		var integrationId string

		if resource.GetPlatform() != nil {
			integrationId = resource.GetPlatform().Id
		} else {
			integrationId = resource.GetCatalog().Id
		}

		resourcePathString := strings.Join(lo.Map(resource.ResourcePath, func(resourceNode *entitiesv1alpha.ResourceNode, _ int) string {
			return resourceNode.Name
		}), "/")

		return fmt.Sprintf("%s/%s", integrationId, resourcePathString)
	})

	if last != "" {
		resourceNames = lo.Filter(resourceNames, func(resourceName string, _ int) bool {
			return strings.Contains(resourceName, last)
		})
	}

	return resourceNames, cobra.ShellCompDirectiveNoSpace
}
