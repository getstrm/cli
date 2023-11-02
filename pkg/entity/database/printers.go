package database

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

// listPrinters
// printers that can handle the output of the list command
func listPrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListDatabasesResponse)
	common.RenderTable(table.Row{
		"ID", "Name", "Type",
	}, lo.Map(listResponse.Databases, func(database *DataCatalog_Database, _ int) table.Row {
		return table.Row{
			database.Id,
			database.DisplayName,
			database.Type,
		}
	}))
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListDatabasesResponse)
	for _, database := range listResponse.Databases {
		fmt.Println(database.Id, database.DisplayName, database.Type)
	}
}
