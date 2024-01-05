package table

import (
	catalogApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	platformApi "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"strings"
)

var printer common.Printer

type listTablePrinter struct{}
type listPlainPrinter struct{}

// listPrinters
// printers that can handle the output of the list command
func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

func (p listTablePrinter) Print(data interface{}) {
	rows := lo.Map(toTables(data), func(catalogTable *entities.Table, _ int) table.Row {
		return table.Row{
			catalogTable.Id,
			catalogTable.Name,
			strings.Join(catalogTable.Tags, ","),
		}
	})
	headers := table.Row{
		"ID", "Name", "Tags",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	for _, t := range toTables(data) {
		fmt.Println(t.Id, t.Name, strings.Join(t.Tags, ","))
	}
}

func toTables(data interface{}) []*entities.Table {
	var tables []*entities.Table
	if listResponse, ok := (data).(*catalogApi.ListTablesResponse); ok {
		tables = listResponse.Tables
	} else if listResponse, ok := (data).(*platformApi.ListTablesResponse); ok {
		tables = listResponse.Tables
	} else {
		tables = nil
	}
	return tables
}
