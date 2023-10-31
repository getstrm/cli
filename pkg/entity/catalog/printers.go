package catalog

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/data_catalogs/v1alpha"
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/entities/v1alpha"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

func availablePrinters() map[string]common.Printer {
	return common.MergePrinterMaps(
		common.DefaultPrinters,
		map[string]common.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

type listTablePrinter struct{}
type listPlainPrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListCatalogsResponse)
	rows := lo.Map(listResponse.Catalogs, func(catalog *DataCatalog, _ int) table.Row {
		return table.Row{
			catalog.Id,
			catalog.Type,
		}
	})

	headers := table.Row{
		"ID",
		"Type",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListCatalogsResponse)
	for _, catalog := range listResponse.Catalogs {
		fmt.Println(catalog.Id, catalog.Type)
	}
}
