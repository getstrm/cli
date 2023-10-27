package catalog

import (
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
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
	listResponse, _ := (data).(*data_policiesv1alpha.ListCatalogsResponse)
	rows := lo.Map(listResponse.Catalogs, func(catalog *data_policiesv1alpha.DataCatalog, _ int) table.Row {
		return table.Row{
			catalog.Id,
			catalog.Type,
		}
	})

	headers := table.Row{
		"ID",
		"Type",
	}
	util.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListCatalogsResponse)
	for _, catalog := range listResponse.Catalogs {
		fmt.Println(catalog.Id, catalog.Type)
	}
}
