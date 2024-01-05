package globaltransform

import (
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/global_transforms/v1alpha/global_transformsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/global_transforms/v1alpha"
	"context"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"pace/pace/pkg/common"
	"sigs.k8s.io/yaml"
	"strings"
)

// strings used in the cli

var apiContext context.Context

var client GlobalTransformsServiceClient

func SetupClient(c GlobalTransformsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = c
}

func upsert(_ *cobra.Command, filename *string) error {
	transform, err := readGlobalTransform(*filename)
	if err != nil {
		return err
	}

	req := &UpsertGlobalTransformRequest{
		Transform: transform,
	}
	response, err := client.UpsertGlobalTransform(apiContext, req)

	if err != nil {
		return err
	}

	return common.Print(printer, err, response.Transform)
}

func get(cmd *cobra.Command, ref string) error {
	flags := cmd.Flags()
	typ, _ := flags.GetString(policyTypeFlag)
	req := &GetGlobalTransformRequest{
		Ref:  ref,
		Type: typ,
	}
	response, err := client.GetGlobalTransform(apiContext, req)
	if err != nil {
		return err
	}
	return common.Print(printer, err, response.Transform)
}

func del(cmd *cobra.Command, ref string) error {
	flags := cmd.Flags()
	typ, _ := flags.GetString(policyTypeFlag)
	req := &DeleteGlobalTransformRequest{
		Ref:  ref,
		Type: typ,
	}
	response, err := client.DeleteGlobalTransform(apiContext, req)

	if err != nil {
		return err
	}

	return common.Print(printer, err, response)
}

func list(_ *cobra.Command) error {
	response, err := client.ListGlobalTransforms(apiContext, &ListGlobalTransformsRequest{})
	if err != nil {
		return err
	}
	return common.Print(printer, err, response)
}

// readGlobalTransform
// read a json or yaml encoded policy from the filesystem.
func readGlobalTransform(filename string) (*GlobalTransform, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// check if the file is yaml and convert it to json
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		if err != nil {
			return nil, err
		}
	}
	transform := &GlobalTransform{}
	err = protojson.Unmarshal(file, transform)
	return transform, err
}

func refCompletionFunction(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	response, err := client.ListGlobalTransforms(apiContext, &ListGlobalTransformsRequest{})

	if err != nil {
		return common.CobraCompletionError(err)
	}

	// TODO handle other types of transforms. Currently only one type though
	refs := lo.Map(response.GlobalTransforms, func(t *GlobalTransform, _ int) string {
		return t.GetTagTransform().TagContent
	})
	return refs, cobra.ShellCompDirectiveNoFileComp
}
