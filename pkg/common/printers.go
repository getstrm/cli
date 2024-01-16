package common

import (
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/proto"
	"os"
	"pace/pace/pkg/util"
	"strings"
)

// StandardPrinters
// these are the printers that can always be used because they
// understand every type of protobuf message
var StandardPrinters orderedmap.OrderedMap[string, Printer]

const (
	OutputFormatYaml       = "yaml"
	OutputFormatJson       = "json"
	OutputFormatJsonRaw    = "json-raw"
	OutputFormatTable      = "table"
	OutputFormatPlain      = "plain"
	OutputFormatSimpleYaml = "yaml-simple"

	OutputFormatFlag      = "output"
	OutputFormatFlagShort = "o"
)

func init() {
	// the order is important. The first one is the default output for every command
	StandardPrinters = *orderedmap.NewOrderedMap[string, Printer]()
	StandardPrinters.Set(OutputFormatYaml, ProtoMessageYamlPrinter{})
	StandardPrinters.Set(OutputFormatJson, ProtoMessageJsonPrettyPrinter{})
	StandardPrinters.Set(OutputFormatJsonRaw, ProtoMessageJsonRawPrinter{})
}

type Printer interface {
	// TODO this should be a private func
	Print(data interface{})
}

type ProtoMessageJsonRawPrinter struct{}
type ProtoMessageJsonPrettyPrinter struct{}
type ProtoMessageYamlPrinter struct{}

/*
ConfigurePrinter

this function is called just before the command execution
the output format has already been set.
*/
func ConfigurePrinter(command *cobra.Command, printers orderedmap.OrderedMap[string, Printer]) (Printer, error) {
	outputFormat, _ := command.Flags().GetString(OutputFormatFlag)
	printer, ok := printers.Get(outputFormat)
	if !ok {
		return nil, AbortError("Output format '%v' is not supported. Allowed values: %v", outputFormat, strings.Join(printers.Keys(), ", "))
	}
	return printer, nil
}

// ConfigureExtraPrinters
// this function is called after construction of the Cobra command in case a method wants to provide more than the StandardPrinters
// entities that only need the standard universal printers don't need to call this.
func ConfigureExtraPrinters(cmd *cobra.Command, flags *pflag.FlagSet, printers orderedmap.OrderedMap[string, Printer]) error {
	formats := printers.Keys()
	flags.StringP(OutputFormatFlag, OutputFormatFlagShort, formats[0],
		fmt.Sprintf("output formats [%v]", strings.Join(formats, ", ")))
	err := cmd.RegisterFlagCompletionFunc(OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return formats, cobra.ShellCompDirectiveNoFileComp
	})
	return err
}

func Print(printer Printer, err error, data interface{}) error {
	if err != nil {
		return err
	}
	printer.Print(data)
	return nil
}

func (p ProtoMessageJsonRawPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	rawJson := util.ProtoMessageToRawJson(protoContent)
	fmt.Println(rawJson.String())
	printIfProtoMessageIsEmpty(protoContent)
}

func (p ProtoMessageJsonPrettyPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	prettyJson := util.ProtoMessageToPrettyJson(protoContent)
	fmt.Println(prettyJson.String())
	printIfProtoMessageIsEmpty(protoContent)
}

func (p ProtoMessageYamlPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	yaml := util.ProtoMessageToYaml(protoContent)
	fmt.Println(yaml.String())
	printIfProtoMessageIsEmpty(protoContent)
}

func printIfProtoMessageIsEmpty(protoContent proto.Message) {
	// An alternative to this approach is to use proto.Size(protoContent) == 0, though this way we ensure that we do a
	// deep comparison.
	cloned := proto.Clone(protoContent)
	proto.Reset(cloned)

	if proto.Equal(cloned, protoContent) {
		// To ensure that stdout always contains valid json / yaml, we print to stderr if the message is empty
		fmt.Fprintln(os.Stderr, "No entities of this resource type exist.")
	}
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

func SliceToRow(slice []string) table.Row {
	row := make(table.Row, len(slice))
	for i, value := range slice {
		row[i] = value
	}
	return row
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
