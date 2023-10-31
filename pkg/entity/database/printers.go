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

func availablePrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.DefaultPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListDatabasesResponse)
	rows := lo.Map(listResponse.Databases, func(catalog *DataCatalog_Database, _ int) table.Row {
		return table.Row{
			catalog.Id,
			catalog.DisplayName,
		}
	})
	headers := table.Row{
		"ID", "Name",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListDatabasesResponse)
	for _, database := range listResponse.Databases {
		fmt.Println(database)
	}
}
