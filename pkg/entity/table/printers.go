package table

import (
	data_policiesv1alpha "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/api/data_policies/v1alpha"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"pace/pace/pkg/common"
	"pace/pace/pkg/util"
)

var printer util.Printer

type listTablePrinter struct{}
type listPlainPrinter struct{}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListProcessingPlatformTablesResponse)
	rows := lo.Map(listResponse.Tables, func(group string, _ int) table.Row {
		return table.Row{
			group,
		}
	})
	headers := table.Row{
		"ID",
	}
	util.RenderTable(headers, rows)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListProcessingPlatformTablesResponse)
	for _, t := range listResponse.Tables {
		fmt.Println(t)
	}
}
