package table

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
)

var printer common.Printer

type listTablePrinter struct{}
type listPlainPrinter struct{}

func availablePrinters() map[string]common.Printer {
	return common.MergePrinterMaps(
		common.DefaultPrinters,
		map[string]common.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListTablesResponse)
	rows := lo.Map(listResponse.Tables, func(group string, _ int) table.Row {
		return table.Row{
			group,
		}
	})
	headers := table.Row{
		"ID",
	}
	common.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListTablesResponse)
	for _, t := range listResponse.Tables {
		fmt.Println(t)
	}
}
