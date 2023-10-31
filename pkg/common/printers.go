package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"pace/pace/pkg/util"
	"sigs.k8s.io/yaml"
	"strings"
)

var DefaultPrinters orderedmap.OrderedMap[string, Printer]

func init() {
	// the order is important. The first one is the default output for every command
	DefaultPrinters = *orderedmap.NewOrderedMap[string, Printer]()
	DefaultPrinters.Set(OutputFormatYaml, ProtoMessageYamlPrinter{})
	DefaultPrinters.Set(OutputFormatJson, ProtoMessageJsonPrettyPrinter{})
	DefaultPrinters.Set(OutputFormatJsonRaw, ProtoMessageJsonRawPrinter{})
}

type Printer interface {
	Print(data interface{})
}

type ProtoMessageJsonRawPrinter struct{}
type ProtoMessageJsonPrettyPrinter struct{}
type ProtoMessageYamlPrinter struct{}

// ConfigurePrinter
// this function is called just before the command execution
// the output format has already been set.
func ConfigurePrinter(command *cobra.Command, printers orderedmap.OrderedMap[string, Printer]) Printer {
	outputFormat, _ := command.Flags().GetString(OutputFormatFlag)
	printer, ok := printers.Get(outputFormat)
	if !ok {
		util.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, strings.Join(printers.Keys(), ", "))))
	}
	return printer
}

// ConfigureExtraPrinters
// this function is called after construction of the Cobra command in case a method wants to provide more than the DefaultPrinters
// entities that only need the default universal printers don't need to call this.
func ConfigureExtraPrinters(cmd *cobra.Command, flags *pflag.FlagSet, printers orderedmap.OrderedMap[string, Printer]) {
	formats := printers.Keys()
	flags.StringP(OutputFormatFlag, OutputFormatFlagShort, formats[0],
		fmt.Sprintf("output formats [%v]", strings.Join(formats, ", ")))
	err := cmd.RegisterFlagCompletionFunc(OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return formats, cobra.ShellCompDirectiveNoFileComp
	})
	util.CliExit(err)
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

func protoMessageToRawJson(proto proto.Message) bytes.Buffer {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.MarshalOptions{
		UseProtoNames: true,
	}.Marshal(proto)
	buffer := bytes.Buffer{}
	errCompact := json.Compact(&buffer, marshal)
	util.CliExit(errCompact)
	return buffer
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

func protoMessageToPrettyJson(proto proto.Message) bytes.Buffer {
	prettyJson := bytes.Buffer{}
	errIndent := json.Indent(&prettyJson, protoMessageToRawJson(proto).Bytes(), "", "    ")
	util.CliExit(errIndent)
	return prettyJson
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
