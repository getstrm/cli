package group

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
	listResponse, _ := (data).(*data_policiesv1alpha.ListProcessingPlatformGroupsResponse)
	printTable(listResponse.Groups)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_policiesv1alpha.ListProcessingPlatformGroupsResponse)
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
	util.RenderTable(headers, rows)
}
func printPlain(groups []string) {
	for _, group := range groups {
		fmt.Println(group)
	}
}
