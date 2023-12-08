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

func invokePlugin(cmd *cobra.Command, args []string) {
	plugin := getPluginById(&args[0])
	actionType := Action_Type(Action_Type_value[args[1]])

	request := &InvokePluginRequest{
		PluginId: plugin.Id,
		Action: &Action{
			Type: actionType,
		},
		Parameters: nil,
	}
	addPluginRequestParameters(cmd, plugin, actionType, request)

	response, err := client.InvokePlugin(apiContext, request)
	CliExit(err)
	printResult(plugin, response)
}

// addPluginRequestParameters adds the correct parameters to the request based on the plugin type
// unfortunately, this is necessary, as the interface that the request Parameters implement, is not exported
func addPluginRequestParameters(cmd *cobra.Command, plugin *Plugin, actionType Action_Type, request *InvokePluginRequest) {
	switch actionType {
	case Action_GENERATE_SAMPLE_DATA:
		request.Parameters = &InvokePluginRequest_SampleDataGeneratorParameters{
			SampleDataGeneratorParameters: &SampleDataGenerator_Parameters{
				Payload: *readPayload(cmd, &plugin.Id),
			},
		}

	case Action_GENERATE_DATA_POLICY:
		request.Parameters = &InvokePluginRequest_DataPolicyGeneratorParameters{
			DataPolicyGeneratorParameters: &DataPolicyGenerator_Parameters{
				Payload: *readPayload(cmd, &plugin.Id),
			},
		}
	default:
		CliExit(errors.New(fmt.Sprintf("plugin type %s not supported", actionType)))
	}
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

// printResult ensures that the correct element of the result is extracted and then printed
func printResult(plugin *Plugin, response *InvokePluginResponse) {
	switch response.Result.(type) {
	case *InvokePluginResponse_DataPolicyGeneratorResult:
		dataPolicy := response.GetResult().(*InvokePluginResponse_DataPolicyGeneratorResult).DataPolicyGeneratorResult.DataPolicy
		printer.Print(dataPolicy)
	case *InvokePluginResponse_SampleDataGeneratorResult:
		csv := response.GetResult().(*InvokePluginResponse_SampleDataGeneratorResult).SampleDataGeneratorResult.Data
		fmt.Println(csv)
	}
}
