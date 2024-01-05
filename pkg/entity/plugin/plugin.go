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
	"sigs.k8s.io/yaml"
	"strings"
)

var apiContext context.Context

var client plugins.PluginsServiceClient

func SetupClient(clientConnection plugins.PluginsServiceClient, ctx context.Context) {
	apiContext = ctx
	client = clientConnection
}

func IdsCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	response, err := list()
	if err != nil {
		return common.CobraCompletionError(err)
	}

	if len(args) == 1 && cmd.Name() == pluginCommand {
		plugin := lo.Filter(response.Plugins, func(p *Plugin, _ int) bool {
			return p.Id == args[0]
		})

		if len(plugin) == 1 {
			actions := lo.Map(response.Plugins[0].Actions, func(action *Action, _ int) string {
				return action.Type.String()
			})
			return actions, cobra.ShellCompDirectiveNoFileComp
		} else {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	} else if len(args) == 0 && cmd.Name() == pluginCommand {
		names := lo.Map(response.Plugins, func(p *Plugin, _ int) string {
			return p.Id
		})
		return names, cobra.ShellCompDirectiveNoFileComp
	} else {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func listPlugins() error {
	response, err := list()
	return common.Print(printer, err, response)
}

func list() (*ListPluginsResponse, error) {
	return client.ListPlugins(apiContext, &ListPluginsRequest{})
}

func getPluginById(pluginId *string) (*Plugin, error) {
	req := &ListPluginsRequest{}
	response, err := client.ListPlugins(apiContext, req)
	if err != nil {
		return nil, err
	}
	plugin, found := lo.Find(response.Plugins, func(p *Plugin) bool {
		return p.Id == *pluginId
	})
	if !found {
		return nil, fmt.Errorf("plugin with id %s not found", *pluginId)
	}
	return plugin, nil
}

func invokePlugin(cmd *cobra.Command, args []string) error {
	plugin, err := getPluginById(&args[0])
	if err != nil {
		return err
	}

	actionType := Action_Type(Action_Type_value[args[1]])

	request := &InvokePluginRequest{
		PluginId: plugin.Id,
		Action: &Action{
			Type: actionType,
		},
		Parameters: nil,
	}
	_ = addPluginRequestParameters(cmd, plugin, actionType, request)
	response, err := client.InvokePlugin(apiContext, request)
	if err != nil {
		return err
	}
	return printResult(response)
}

// addPluginRequestParameters adds the correct parameters to the request based on the plugin type
// unfortunately, this is necessary, as the interface that the request Parameters implement, is not exported
func addPluginRequestParameters(cmd *cobra.Command, plugin *Plugin, actionType Action_Type, request *InvokePluginRequest) error {
	switch actionType {
	case Action_GENERATE_SAMPLE_DATA:
		payload, err := readPayload(cmd, &plugin.Id, actionType)
		if err != nil {
			return err
		}
		request.Parameters = &InvokePluginRequest_SampleDataGeneratorParameters{
			SampleDataGeneratorParameters: &SampleDataGenerator_Parameters{
				Payload: *payload,
			},
		}
	case Action_GENERATE_DATA_POLICY:
		payload, err := readPayload(cmd, &plugin.Id, actionType)
		if err != nil {
			return err
		}
		request.Parameters = &InvokePluginRequest_DataPolicyGeneratorParameters{
			DataPolicyGeneratorParameters: &DataPolicyGenerator_Parameters{
				Payload: *payload,
			},
		}
	default:
		return fmt.Errorf("plugin type %s not supported", actionType)
	}
	return nil
}

func readPayload(cmd *cobra.Command, pluginId *string, actionType Action_Type) (*string, error) {
	fileName, _ := cmd.Flags().GetString(common.PluginPayloadFlag)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		file, err = yaml.YAMLToJSON(file)
		if err != nil {
			return nil, err
		}
	}

	err = validatePayload(pluginId, actionType, file)
	if err != nil {
		return nil, err
	}
	byte64EncodedJson := base64.StdEncoding.EncodeToString(file)
	return &byte64EncodedJson, nil
}

func validatePayload(pluginId *string, actionType Action_Type, payload []byte) error {
	req := &GetPayloadJSONSchemaRequest{
		PluginId: *pluginId,
		Action: &Action{
			Type: actionType,
		},
	}
	jsonSchema, err := client.GetPayloadJSONSchema(apiContext, req)
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema.Schema)
	documentLoader := gojsonschema.NewBytesLoader(payload)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		var errMsg = strings.Builder{}
		errMsg.WriteString(fmt.Sprintf("payload validation failed for plugin %s:\n", *pluginId))
		for _, err := range result.Errors() {
			errMsg.WriteString(fmt.Sprintf("- %s\n", err))
		}
		return errors.New(errMsg.String())
	}
	return nil
}

// printResult ensures that the correct element of the result is extracted and then printed
func printResult(response *InvokePluginResponse) error {
	switch response.Result.(type) {
	case *InvokePluginResponse_DataPolicyGeneratorResult:
		dataPolicy := response.GetResult().(*InvokePluginResponse_DataPolicyGeneratorResult).DataPolicyGeneratorResult.DataPolicy
		return common.Print(printer, nil, dataPolicy)
	case *InvokePluginResponse_SampleDataGeneratorResult:
		csv := response.GetResult().(*InvokePluginResponse_SampleDataGeneratorResult).SampleDataGeneratorResult.Data
		fmt.Println(csv)
		return nil
	}
	return nil
}
