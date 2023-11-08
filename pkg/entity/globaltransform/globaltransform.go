package globaltransform

import (
	. "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/global_transforms/v1alpha/global_transformsv1alphagrpc"
	v1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/global_transforms/v1alpha"
	"context"
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
	printer.Print(response)
}

func get(cmd *cobra.Command, ref string, typ string) {
	// return a global policy from the Pace database.
	req := &GetGlobalTransformRequest{
		RefAndType: &v1alpha.GlobalTransform_RefAndType{
			Ref:  ref,
			Type: typ,
		},
	}
	response, err := client.GetGlobalTransform(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func delete(cmd *cobra.Command, ref string, typ string) {
	// return a data policy from the Pace database.
	refAndType := []*v1alpha.GlobalTransform_RefAndType{
		&v1alpha.GlobalTransform_RefAndType{
			Ref:  ref,
			Type: typ,
		},
	}
	req := &DeleteGlobalTransformRequest{
		RefAndTypes: refAndType,
	}
	response, err := client.DeleteGlobalTransform(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

func list(_ *cobra.Command) {
	req := &ListGlobalTransformsRequest{}
	response, err := client.ListGlobalTransforms(apiContext, req)
	CliExit(err)
	printer.Print(response)
}

// readGlobalTransform
// read a json or yaml encoded policy from the filesystem.
func readGlobalTransform(filename string) *v1alpha.GlobalTransform {
	file, err := os.ReadFile(filename)
	CliExit(err)

	// check if the file is yaml and convert it to json
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		CliExit(err)
	}
	transform := &v1alpha.GlobalTransform{}
	err = protojson.Unmarshal(file, transform)
	CliExit(err)
	return transform
}
