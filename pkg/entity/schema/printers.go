package schema

import (
	datapolicies "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var printer util.Printer

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
	listResponse, _ := (data).(*datapolicies.ListSchemasResponse)
	rows := lo.Map(listResponse.Schemas, func(schema *datapolicies.DataCatalog_Schema, _ int) table.Row {
		return table.Row{
			schema.Id,
			schema.Name,
		}
	})
	headers := table.Row{
		"ID", "Name",
	}
	util.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*datapolicies.ListSchemasResponse)
	for _, schema := range listResponse.Schemas {
		fmt.Println(schema)
	}
}
