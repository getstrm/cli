package plugin

import (
	plugins "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/plugins/v1alpha/pluginsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/plugins/v1alpha"
	data_policy_generatorsv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/plugins/data_policy_generators/v1alpha"
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/anypb"
	"os"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
)

var apiContext context.Context

var client plugins.PluginsServiceClient

func SetupClient(clientConnection plugins.PluginsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func IdsCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	req := &ListPluginsRequest{}
	response, err := client.ListPlugins(apiContext, req)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := lo.Map(response.Plugins, func(p *Plugin, _ int) string {
		return p.Id
	})

	return names, cobra.ShellCompDirectiveNoFileComp
}

func getPluginById(pluginId *string) *Plugin {
	req := &ListPluginsRequest{}
	response, err := client.ListPlugins(apiContext, req)
	CliExit(err)
	plugin, found := lo.Find(response.Plugins, func(p *Plugin) bool {
		return p.Id == *pluginId
	})
	if !found {
		CliExit(errors.New(fmt.Sprintf("plugin with id %s not found", *pluginId)))
	}
	return plugin
}

func invokePlugin(cmd *cobra.Command, pluginId *string) {
	plugin := getPluginById(pluginId)
	payload := unmarshalPayload(cmd)
	switch plugin.PluginType {
	case PluginType_DATA_POLICY_GENERATOR:
		invokeDataPolicyGenerator(plugin, payload)
	default:
		CliExit(errors.New(fmt.Sprintf("plugin type %s not supported", plugin.PluginType)))
	}
}

func invokeDataPolicyGenerator(plugin *Plugin, payload *anypb.Any) {
	req := &InvokeDataPolicyGeneratorRequest{
		PluginId: &plugin.Id,
		Payload:  payload,
	}
	response, err := client.InvokeDataPolicyGenerator(apiContext, req)
	CliExit(err)
	printer.Print(response.DataPolicy)
}

func unmarshalPayload(cmd *cobra.Command) *anypb.Any {
	fileName := GetStringAndErr(cmd.Flags(), common.PluginPayloadFlag)
	_, err := os.ReadFile(fileName)
	CliExit(err)
	// Todo: actually unmarshal from yaml/json based on descriptor
	req := &data_policy_generatorsv1alpha.OpenAIDataPolicyGeneratorPayload{}
	payload, err := anypb.New(req)
	CliExit(err)
	return payload
}
