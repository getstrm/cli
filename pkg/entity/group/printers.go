package group

import (
	. "buf.build/gen/go/getstrm/pace/protocolbuffers/go/getstrm/pace/api/processing_platforms/v1alpha"
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
	listResponse, _ := (data).(*ListGroupsResponse)
	printTable(listResponse.Groups)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*ListGroupsResponse)
	for _, group := range listResponse.Groups {
		fmt.Println(group)
	}
}

func printTable(groups []string) {
	rows := lo.Map(groups, func(group string, _ int) table.Row {
		return table.Row{
			group,
		}
	})
	headers := table.Row{
		"Name",
	}
	common.RenderTable(headers, rows)
}
func printPlain(groups []string) {
	for _, group := range groups {
		fmt.Println(group)
	}
}
