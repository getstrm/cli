package table

import (
	api "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	entities "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
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

func availablePrinters() orderedmap.OrderedMap[string, common.Printer] {
	printers := common.StandardPrinters.Copy()
	printers.Set(common.OutputFormatTable, listTablePrinter{})
	printers.Set(common.OutputFormatPlain, listPlainPrinter{})
	return *printers
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*api.ListTablesResponse)
	rows := lo.Map(listResponse.Tables, func(catalogTable *entities.DataCatalog_Table, _ int) table.Row {
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
	listResponse, _ := (data).(*api.ListTablesResponse)
	for _, t := range listResponse.Tables {
		fmt.Println(t.Id, t.Name, strings.Join(t.Tags, ","))
	}
}
