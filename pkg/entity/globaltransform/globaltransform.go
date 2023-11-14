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
	. "pace/pace/pkg/util"
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

func upsert(_ *cobra.Command, filename *string) {
	transform := readGlobalTransform(*filename)
	req := &UpsertGlobalTransformRequest{
		Transform: transform,
	}
	response, err := client.UpsertGlobalTransform(apiContext, req)
	CliExit(err)
	printer.Print(response.Transform)
}

func get(cmd *cobra.Command, ref string) {
	flags := cmd.Flags()
	typ := GetStringAndErr(flags, policyTypeFlag)
	req := &GetGlobalTransformRequest{
		Ref: ref, Type: typ,
	}
	response, err := client.GetGlobalTransform(apiContext, req)
	CliExit(err)
	printer.Print(response.Transform)
}

func del(cmd *cobra.Command, ref string) {
	flags := cmd.Flags()
	typ := GetStringAndErr(flags, policyTypeFlag)
	req := &DeleteGlobalTransformRequest{
		RefAndTypes: []*GlobalTransform_RefAndType{
			refAndType(typ, ref),
		},
	}
	response, err := client.DeleteGlobalTransform(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func list(_ *cobra.Command) {
	response, err := client.ListGlobalTransforms(apiContext, &ListGlobalTransformsRequest{})
	CliExit(err)
	printer.Print(response)
}

// readGlobalTransform
// read a json or yaml encoded policy from the filesystem.
func readGlobalTransform(filename string) *GlobalTransform {
	file, err := os.ReadFile(filename)
	CliExit(err)

	// check if the file is yaml and convert it to json
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		CliExit(err)
	}
	transform := &GlobalTransform{}
	err = protojson.Unmarshal(file, transform)
	CliExit(err)
	return transform
}

func refAndType(typ string, ref string) *GlobalTransform_RefAndType {
	return &GlobalTransform_RefAndType{
		Ref:  ref,
		Type: typ,
	}
}
func refCompletionFunction(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	response, err := client.ListGlobalTransforms(apiContext, &ListGlobalTransformsRequest{})
	CliExit(err)
	refs := lo.Map(response.GlobalTransforms, func(t *GlobalTransform, _ int) string {
		return t.Ref
	})
	return refs, cobra.ShellCompDirectiveNoFileComp
}
