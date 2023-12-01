package plugin

import (
	plugins "buf.build/gen/go/getstrm/pace/grpc/go/getstrm/pace/api/plugins/v1alpha/pluginsv1alphagrpc"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/plugins/v1alpha"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"pace/pace/pkg/common"
	. "pace/pace/pkg/util"
	"sigs.k8s.io/yaml"
	"strings"
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
	base64EncodedJsonPayload := readPayload(cmd, pluginId)
	switch plugin.PluginType {
	case PluginType_DATA_POLICY_GENERATOR:
		invokeDataPolicyGenerator(plugin, base64EncodedJsonPayload)
	default:
		CliExit(errors.New(fmt.Sprintf("plugin type %s not supported", plugin.PluginType)))
	}
}

func invokeDataPolicyGenerator(plugin *Plugin, payload *string) {
	//req := &InvokeDataPolicyGeneratorRequest{
	//	PluginId: &plugin.Id,
	//	Payload: *payload,
	//}
	//response, err := client.InvokeDataPolicyGenerator(apiContext, req)
	//CliExit(err)
	//printer.Print(response.DataPolicy)
}

func readPayload(cmd *cobra.Command, pluginId *string) *string {
	fileName := GetStringAndErr(cmd.Flags(), common.PluginPayloadFlag)
	file, err := os.ReadFile(fileName)
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		CliExit(err)
	}
	CliExit(err)
	validatePayload(pluginId, file)
	byte64EncodedJson := base64.StdEncoding.EncodeToString(file)
	return &byte64EncodedJson
}

func validatePayload(pluginId *string, payload []byte) {
	req := &GetPayloadJSONSchemaRequest{PluginId: *pluginId}
	jsonSchema, err := client.GetPayloadJSONSchema(apiContext, req)
	CliExit(err)
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema.Schema)
	documentLoader := gojsonschema.NewBytesLoader(payload)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	CliExit(err)
	if !result.Valid() {
		var errMsg = strings.Builder{}
		errMsg.WriteString(fmt.Sprintf("payload validation failed for plugin %s:\n", *pluginId))
		for _, err := range result.Errors() {
			errMsg.WriteString(fmt.Sprintf("- %s\n", err))
		}
		CliExit(errors.New(errMsg.String()))
	}
}
