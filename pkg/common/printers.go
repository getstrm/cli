package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"pace/pace/pkg/util"
	"sigs.k8s.io/yaml"
)

type Printer interface {
	Print(data interface{})
}

/*
TODO BvD proposes this mechanism instead of the string concat
which would allow better error handling for absent printers.

type PrinterKey struct {
	format  string
	command string
}
var DefaultPrinters = map[PrinterKey]Printer {
	PrinterKey{
		format: common.OutputFormatJson, command: common.ListCommandName,
	}: ProtoMessageJsonPrettyPrinter{},
}

*/

var DefaultPrinters = map[string]Printer{
	OutputFormatJson + ListCommandName:      ProtoMessageJsonPrettyPrinter{},
	OutputFormatJson + GetCommandName:       ProtoMessageJsonPrettyPrinter{},
	OutputFormatJson + DeleteCommandName:    ProtoMessageJsonPrettyPrinter{},
	OutputFormatJson + CreateCommandName:    ProtoMessageJsonPrettyPrinter{},
	OutputFormatJson + UpsertCommandName:    ProtoMessageJsonPrettyPrinter{},
	OutputFormatYaml + ListCommandName:      ProtoMessageYamlPrinter{},
	OutputFormatYaml + GetCommandName:       ProtoMessageYamlPrinter{},
	OutputFormatYaml + DeleteCommandName:    ProtoMessageYamlPrinter{},
	OutputFormatYaml + CreateCommandName:    ProtoMessageYamlPrinter{},
	OutputFormatYaml + UpsertCommandName:    ProtoMessageYamlPrinter{},
	OutputFormatJsonRaw + ListCommandName:   ProtoMessageJsonRawPrinter{},
	OutputFormatJsonRaw + GetCommandName:    ProtoMessageJsonRawPrinter{},
	OutputFormatJsonRaw + DeleteCommandName: ProtoMessageJsonRawPrinter{},
	OutputFormatJsonRaw + CreateCommandName: ProtoMessageJsonRawPrinter{},
	OutputFormatJsonRaw + UpsertCommandName: ProtoMessageJsonRawPrinter{},
}

type ProtoMessageJsonRawPrinter struct{}
type ProtoMessageJsonPrettyPrinter struct{}
type ProtoMessageYamlPrinter struct{}

func ConfigurePrinter(command *cobra.Command, printers map[string]Printer) Printer {
	outputFormat, _ := command.Flags().GetString(OutputFormatFlag)
	p := printers[outputFormat+command.Parent().Name()]
	// TODO this mechanism is not suitable for correctly providing feedback for entities that do not support the
	// default or plain mechanism.
	if p == nil {
		util.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, OutputFormatFlagAllowedValuesText)))
	}
	return p
}

func (p ProtoMessageJsonRawPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	rawJson := protoMessageToRawJson(protoContent)
	fmt.Println(string(rawJson.Bytes()))
}

func (p ProtoMessageJsonPrettyPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	prettyJson := protoMessageToPrettyJson(protoContent)
	fmt.Println(string(prettyJson.Bytes()))
}

func (p ProtoMessageYamlPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	jsonBytes := protoMessageToRawJson(protoContent)
	m, _ := yaml.JSONToYAML(jsonBytes.Bytes())
	fmt.Println(string(m))
}

func protoMessageToPrettyJson(proto proto.Message) bytes.Buffer {
	return PrettifyJson(protoMessageToRawJson(proto))
}

func protoMessageToRawJson(proto proto.Message) bytes.Buffer {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.MarshalOptions{
		UseProtoNames: true,
	}.Marshal(proto)
	return CompactJson(marshal)
}

func CompactJson(rawJson []byte) bytes.Buffer {
	buffer := bytes.Buffer{}
	errCompact := json.Compact(&buffer, rawJson)
	util.CliExit(errCompact)
	return buffer
}

func PrettifyJson(rawJson bytes.Buffer) bytes.Buffer {
	prettyJson := bytes.Buffer{}
	errIndent := json.Indent(&prettyJson, rawJson.Bytes(), "", "    ")
	util.CliExit(errIndent)
	return prettyJson
}

func MergePrinterMaps(maps ...map[string]Printer) (result map[string]Printer) {
	return lo.Assign[string, Printer](maps...)
}

func RenderTable(headers table.Row, rows []table.Row) {
	if len(rows) == 0 {
		fmt.Println("No entities of this resource type exist.")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(headers)
		t.AppendSeparator()
		t.AppendRows(rows)
		t.SetStyle(noBordersStyle)
		t.Render()
	}
}

var noBordersStyle = table.Style{
	Name:    "StyleNoBorders",
	Options: table.OptionsNoBorders,
	Title:   table.TitleOptionsDefault,
	Format:  table.FormatOptionsDefault,
	Box: table.BoxStyle{
		BottomLeft:       " ",
		BottomRight:      " ",
		BottomSeparator:  " ",
		EmptySeparator:   text.RepeatAndTrim(" ", text.RuneWidthWithoutEscSequences(" ")),
		Left:             " ",
		LeftSeparator:    " ",
		MiddleHorizontal: " ",
		MiddleSeparator:  " ",
		MiddleVertical:   " ",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            " ",
		RightSeparator:   " ",
		TopLeft:          " ",
		TopRight:         " ",
		TopSeparator:     " ",
		UnfinishedRow:    "...",
	},
}
