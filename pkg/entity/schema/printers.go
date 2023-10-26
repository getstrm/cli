package schema

import (
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListSchemasResponse)
	printTable(listResponse.Schemas)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListSchemasResponse)
	printPlain(listResponse.Schemas)
}

func printTable(schemas []*data_policiesv1alpha.DataCatalog_Schema) {
	rows := make([]table.Row, 0, len(schemas))
	for _, schema := range schemas {
		row := table.Row{
			schema.Id,
			schema.Name,
		}
		rows = append(rows, row)
	}

	headers := table.Row{
		"ID", "Name",
	}
	util.RenderTable(headers, rows)
}
func printPlain(schemas []*data_policiesv1alpha.DataCatalog_Schema) {
	for _, schema := range schemas {
		fmt.Println(schema)
	}
}
